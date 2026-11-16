[CmdletBinding()]
param([switch]$Apply)

. "$PSScriptRoot/phase2a-common.ps1"
Connect-Keycloak
$definition = Get-AuthDefinition
$policy = $definition.Policy
$clientTemplate = $definition.Client
$before = Get-AuthObjects
$registrationAction = @($before.RequiredActions | Where-Object alias -CEQ 'webauthn-register-passwordless')
if ($registrationAction.Count -ne 1 -or !$registrationAction[0].enabled -or $registrationAction[0].defaultAction) {
    throw 'Passwordless registration must be enabled and must not be a realm-wide default action.'
}
foreach ($roleName in @('administrator', 'dispatcher', 'operator', 'responder')) {
    if (@($before.Roles | Where-Object name -CEQ $roleName).Count -ne 1) {
        throw "Required application role is missing: $roleName."
    }
}

$localDirectory = Join-Path $PSScriptRoot 'local'
$secretsPath = Join-Path $localDirectory 'phase2a-secrets.json'
$localState = $null
if (Test-Path -LiteralPath $secretsPath) {
    $localState = Get-Content -LiteralPath $secretsPath -Raw | ConvertFrom-Json
    Assert-ObjectOwnership $before $localState
    Assert-PasswordlessFlow $before.Executions
    Assert-BffClient $before.Client $before.Flow $clientTemplate
    $clientSecret = Invoke-Keycloak GET "$script:RealmApi/clients/$($before.Client.id)/client-secret"
    if (!$localState.clientSecret -or $clientSecret.value -cne $localState.clientSecret) {
        throw 'Client secret does not match local state; no overwrite or rotation attempted.'
    }
} elseif ($before.Client -or $before.User -or $before.Flow) {
    throw 'Phase 2A objects already exist without local ownership state; refusing to adopt or replace them.'
}

$humanRoles = $definition.RoleScope
$scopeRoles = @()
$assignedRoles = @()
$existingCredentials = $null
if ($before.Client) {
    $scopeRoles = @(Invoke-Keycloak GET "$script:RealmApi/clients/$($before.Client.id)/scope-mappings/realm")
    if (@($scopeRoles | Where-Object { $_.name -cnotin $humanRoles }).Count -ne 0) {
        throw 'Unexpected BFF role scope; no existing mappings will be removed.'
    }
}
if ($before.User) {
    $assignedRoles = @(Invoke-Keycloak GET "$script:RealmApi/users/$($before.User.id)/role-mappings/realm")
    if (@($assignedRoles | Where-Object { $_.name -cnotin @('dispatcher', 'default-roles-ce1-dev') }).Count -ne 0) {
        throw 'Unexpected dispatcher roles; no existing assignments will be removed.'
    }
    $existingCredentials = ConvertTo-Json -InputObject @(Invoke-Keycloak GET "$script:RealmApi/users/$($before.User.id)/credentials") -Depth 30 -Compress
}

$policyDifferences = @(Get-PolicyDifferences $before.Realm $policy)
Write-Output 'Plan: configure the localhost policy, dedicated passwordless flow, confidential BFF client and development dispatcher.'
if (!$Apply) {
    Write-Output 'Inspection only. Run with -Apply to make approved changes.'
    return
}

$protectedBefore = Get-ProtectedConfiguration $before $policy
if ($policyDifferences.Count -gt 0) {
    # Send only the owned policy fields, never an entire realm export.
    $null = Invoke-Keycloak PUT $script:RealmApi $policy
}
if (!$localState) {
    if (!(Test-Path -LiteralPath $localDirectory)) {
        $null = New-Item -ItemType Directory -Path $localDirectory
    }
    $localState = [pscustomobject]@{
        owner = [guid]::NewGuid().ToString()
        clientSecret = $null
        bootstrapPassword = New-RandomSecret
        userId = $null
    }
}

$flow = $before.Flow
if (!$flow) {
    $null = Invoke-Keycloak POST "$script:RealmApi/authentication/flows" @{
        alias = $script:FlowName
        description = "Managed by infra/auth Phase 2A: $($localState.owner)"
        providerId = 'basic-flow'
        topLevel = $true
        builtIn = $false
    }
    foreach ($provider in @('auth-username-form', 'webauthn-authenticator-passwordless')) {
        $null = Invoke-Keycloak POST "$script:RealmApi/authentication/flows/$script:FlowName/executions/execution" @{ provider = $provider }
    }
    $executions = @(Invoke-Keycloak GET "$script:RealmApi/authentication/flows/$script:FlowName/executions")
    foreach ($execution in $executions) {
        $null = Invoke-Keycloak PUT "$script:RealmApi/authentication/flows/$script:FlowName/executions" @{
            id = $execution.id
            requirement = 'REQUIRED'
        }
    }
    $flows = Invoke-Keycloak GET "$script:RealmApi/authentication/flows"
    $flow = $flows | Where-Object alias -CEQ $script:FlowName
}

$client = $before.Client
if (!$client) {
    $clientTemplate.attributes | Add-Member NoteProperty $script:OwnerAttribute $localState.owner
    $clientTemplate.authenticationFlowBindingOverrides = @{ browser = $flow.id }
    $null = Invoke-Keycloak POST "$script:RealmApi/clients" $clientTemplate
    $clients = Invoke-Keycloak GET "$script:RealmApi/clients?clientId=ce1-bff-dev"
    $client = @($clients)[0]
    $generatedSecret = Invoke-Keycloak GET "$script:RealmApi/clients/$($client.id)/client-secret"
    $localState.clientSecret = $generatedSecret.value
}
$missingScopeRoles = @($before.Roles | Where-Object {
    $_.name -cin $humanRoles -and $_.name -cnotin @($scopeRoles | ForEach-Object name)
})
if ($missingScopeRoles.Count -gt 0) {
    $null = Invoke-Keycloak POST "$script:RealmApi/clients/$($client.id)/scope-mappings/realm" $missingScopeRoles
}

$user = $before.User
if (!$user) {
    # Enrollment-only password; do not remove it before manual passkey proof.
    # temporary=false avoids forcing a second password before enrollment.
    $null = Invoke-Keycloak POST "$script:RealmApi/users" @{
        username = 'dev-dispatcher'
        enabled = $true
        firstName = 'Development'
        lastName = 'Dispatcher'
        emailVerified = $false
        requiredActions = @('webauthn-register-passwordless')
        credentials = @(@{ type = 'password'; temporary = $false; value = $localState.bootstrapPassword })
    }
    $users = Invoke-Keycloak GET "$script:RealmApi/users?username=dev-dispatcher&exact=true"
    $user = @($users)[0]
    $localState.userId = $user.id
}
if (@($assignedRoles | ForEach-Object name) -cnotcontains 'dispatcher') {
    $dispatcherRole = $before.Roles | Where-Object name -CEQ 'dispatcher'
    $null = Invoke-Keycloak POST "$script:RealmApi/users/$($user.id)/role-mappings/realm" @($dispatcherRole)
}

$after = Get-AuthObjects
Assert-ObjectOwnership $after $localState
Assert-PasswordlessFlow $after.Executions
Assert-BffClient $after.Client $after.Flow $clientTemplate
$protectedAfter = Get-ProtectedConfiguration $after $policy
foreach ($section in $protectedBefore.PSObject.Properties.Name) {
    if ($protectedBefore.$section -cne $protectedAfter.$section) {
        throw "Setup changed protected configuration: $section. Inspect the realm; no automatic rollback attempted."
    }
}
if ($before.Client -and $before.Client.id -cne $after.Client.id) {
    throw 'Existing BFF client ID changed.'
}
if ($before.Flow -and $before.Flow.id -cne $after.Flow.id) {
    throw 'Existing passwordless flow ID changed.'
}
if ($before.User) {
    $credentialsAfter = ConvertTo-Json -InputObject @(Invoke-Keycloak GET "$script:RealmApi/users/$($after.User.id)/credentials") -Depth 30 -Compress
    if ($existingCredentials -cne $credentialsAfter) {
        throw 'Existing dispatcher credentials changed.'
    }
}
$secretAfter = Invoke-Keycloak GET "$script:RealmApi/clients/$($after.Client.id)/client-secret"
if ($secretAfter.value -cne $localState.clientSecret) {
    throw 'BFF client secret changed or does not match local state.'
}
if (!(Test-Path -LiteralPath $secretsPath)) {
    # First-time setup writes this file once. Reruns never rewrite it.
    Write-NewLocalJson $secretsPath $localState
}
Write-Output "PASS: setup completed; $script:AdminWrites Admin API write request(s). Existing local files were not rewritten."
Write-Output 'MANUAL: passkey enrollment and login remain separate. No existing password or passkey was changed.'
