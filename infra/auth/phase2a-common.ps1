# Keycloak administration and checks shared by the Phase 2A scripts.
Set-StrictMode -Version 2
$ErrorActionPreference = 'Stop'
$script:AuthRoot = $PSScriptRoot
$script:KeycloakUrl = 'http://localhost:8081'
$script:RealmApi = '/admin/realms/ce1-dev'
$script:FlowName = 'ce1-bff-passkey'
$script:OwnerAttribute = 'ce1.phase2a.owner'
$script:AdminWrites = 0

function Read-AuthSettings([string]$EnvironmentFile) {
    $settings = @{}
    foreach ($line in Get-Content -LiteralPath $EnvironmentFile) {
        if ($line -match '^([A-Z_]+)=(.*)$') {
            $settings[$matches[1]] = $matches[2]
        }
    }
    return $settings
}

function Connect-Keycloak([string]$EnvironmentFile = "$script:AuthRoot/.env") {
    $settings = Read-AuthSettings $EnvironmentFile
    $script:AuthSettings = $settings
    if (!$settings.AUTH_ADMIN_USERNAME -or !$settings.AUTH_ADMIN_PASSWORD) {
        throw 'Missing local administrator settings in .env.'
    }
    try {
        $token = Invoke-RestMethod -Method Post -TimeoutSec 20 -Uri "$script:KeycloakUrl/realms/master/protocol/openid-connect/token" -Body @{
            client_id = 'admin-cli'
            grant_type = 'password'
            username = $settings.AUTH_ADMIN_USERNAME
            password = $settings.AUTH_ADMIN_PASSWORD
        }
        $script:AdminHeaders = @{ Authorization = 'Bearer ' + $token.access_token }
    } catch {
        throw 'Local administrator authentication failed; credentials and response body withheld.'
    }
    $server = Invoke-Keycloak GET '/admin/serverinfo'
    if ($server.systemInfo.version -cne '26.8.1') {
        throw 'These scripts are verified for Keycloak 26.8.1 only; review before upgrading.'
    }
}

function Get-AuthDefinition {
    # Host-side utilities read the same definition as Keycloak. Resolve only the
    # documented environment settings; no general-purpose template engine.
    $realm = Get-Content -LiteralPath "$script:AuthRoot/keycloak/ce1-dev-realm.json" -Raw | ConvertFrom-Json
    $realm.sslRequired = 'none'
    $realm.webAuthnPolicyPasswordlessRpId = 'localhost'
    if ($script:AuthSettings['AUTH_SSL_REQUIRED']) {
        $realm.sslRequired = $script:AuthSettings['AUTH_SSL_REQUIRED']
    }
    if ($script:AuthSettings['AUTH_WEBAUTHN_RP_ID']) {
        $realm.webAuthnPolicyPasswordlessRpId = $script:AuthSettings['AUTH_WEBAUTHN_RP_ID']
    }
    $client = $realm.clients | Where-Object clientId -CEQ 'ce1-bff-dev'
    $client.redirectUris = @('http://localhost:8082/auth/callback')
    if ($script:AuthSettings['AUTH_BFF_REDIRECT_URI']) {
        $client.redirectUris = @($script:AuthSettings['AUTH_BFF_REDIRECT_URI'])
    }
    # API-based local setup may let Keycloak generate a secret. Native import
    # requires an explicitly supplied secret shared with the BFF.
    $client.PSObject.Properties.Remove('secret')
    $policy = [ordered]@{}
    foreach ($setting in $realm.PSObject.Properties | Where-Object Name -Like 'webAuthnPolicyPasswordless*') {
        $policy[$setting.Name] = $setting.Value
    }
    return [pscustomobject]@{
        Realm = $realm
        Client = $client
        Policy = [pscustomobject]$policy
        RoleScope = @($realm.scopeMappings | Where-Object client -CEQ 'ce1-bff-dev' | ForEach-Object roles)
    }
}

function Invoke-Keycloak([ValidateSet('GET', 'POST', 'PUT')][string]$Method, [string]$Path, $Body = $null) {
    $requestParameters = @{
        Method = $Method
        Uri = "$script:KeycloakUrl$Path"
        Headers = $script:AdminHeaders
        TimeoutSec = 25
    }
    if ($null -ne $Body) {
        $requestParameters.ContentType = 'application/json'
        $requestParameters.Body = ConvertTo-Json -InputObject $Body -Depth 40 -Compress
    }
    try {
        $response = Invoke-RestMethod @requestParameters
    } catch {
        $status = 'unavailable'
        if ($_.Exception.Response) {
            $status = [int]$_.Exception.Response.StatusCode
        }
        throw "Keycloak $Method $Path failed (HTTP $status); request/response bodies withheld."
    }
    if ($Method -ne 'GET') {
        $script:AdminWrites++
    }
    return $response
}

function Get-AuthObjects {
    $clients = @(Invoke-Keycloak GET "$script:RealmApi/clients")
    $client = $clients | Where-Object clientId -CEQ 'ce1-bff-dev'
    $users = @(Invoke-Keycloak GET "$script:RealmApi/users?username=dev-dispatcher&exact=true")
    $flows = @(Invoke-Keycloak GET "$script:RealmApi/authentication/flows")
    $flow = $flows | Where-Object alias -CEQ $script:FlowName
    if (@($client).Count -gt 1 -or $users.Count -gt 1 -or @($flow).Count -gt 1) {
        throw 'Duplicate Phase 2A objects; inspect the realm before continuing.'
    }
    $user = $null
    if ($users.Count -eq 1) {
        $user = Invoke-Keycloak GET "$script:RealmApi/users/$($users[0].id)"
    }
    $executions = @()
    if ($flow) {
        $executions = @(Invoke-Keycloak GET "$script:RealmApi/authentication/flows/$script:FlowName/executions")
    }
    return [pscustomobject]@{
        Realm = Invoke-Keycloak GET $script:RealmApi
        Roles = @(Invoke-Keycloak GET "$script:RealmApi/roles")
        Clients = $clients
        Client = $client
        User = $user
        Flow = $flow
        Executions = $executions
        RequiredActions = @(Invoke-Keycloak GET "$script:RealmApi/authentication/required-actions")
    }
}

function Assert-ObjectOwnership($Objects, $LocalState) {
    if (!$Objects.Client -or !$Objects.User -or !$Objects.Flow) {
        throw 'A previously configured object is missing; refusing automatic recreation.'
    }
    if ($Objects.Client.attributes.($script:OwnerAttribute) -cne $LocalState.owner) {
        throw 'BFF client ownership does not match local state.'
    }
    if ($Objects.User.id -cne $LocalState.userId) {
        throw 'Dispatcher ID does not match local state.'
    }
    if ($Objects.Flow.description -cne "Managed by infra/auth Phase 2A: $($LocalState.owner)") {
        throw 'Passwordless flow ownership does not match local state.'
    }
}

function Get-PolicyDifferences($Realm, $Policy) {
    foreach ($field in @(
        'webAuthnPolicyPasswordlessRpEntityName', 'webAuthnPolicyPasswordlessRpId',
        'webAuthnPolicyPasswordlessAttestationConveyancePreference',
        'webAuthnPolicyPasswordlessAuthenticatorAttachment', 'webAuthnPolicyPasswordlessResidentKey',
        'webAuthnPolicyPasswordlessUserVerificationRequirement'
    )) {
        if ($Realm.$field -cne $Policy.$field) {
            $field
        }
    }
    if (($Realm.webAuthnPolicyPasswordlessSignatureAlgorithms -join ',') -cne ($Policy.webAuthnPolicyPasswordlessSignatureAlgorithms -join ',')) {
        'webAuthnPolicyPasswordlessSignatureAlgorithms'
    }
    if (($Realm.webAuthnPolicyPasswordlessExtraOrigins -join ',') -cne ($Policy.webAuthnPolicyPasswordlessExtraOrigins -join ',')) {
        'webAuthnPolicyPasswordlessExtraOrigins'
    }
}

function Assert-PasswordlessFlow($Executions) {
    if ($Executions.Count -ne 2 -or
        $Executions[0].providerId -cne 'auth-username-form' -or
        $Executions[1].providerId -cne 'webauthn-authenticator-passwordless') {
        throw 'Expected exactly Username Form followed by WebAuthn Passwordless Authenticator.'
    }
    foreach ($execution in $Executions) {
        if ($execution.requirement -cne 'REQUIRED' -or $execution.level -ne 0) {
            throw 'Both passwordless flow steps must be REQUIRED at the top level.'
        }
    }
}

function Assert-BffClient($Client, $Flow, $Template) {
    foreach ($field in @(
        'clientId', 'name', 'protocol', 'enabled', 'publicClient', 'bearerOnly',
        'clientAuthenticatorType', 'standardFlowEnabled', 'implicitFlowEnabled',
        'directAccessGrantsEnabled', 'serviceAccountsEnabled', 'fullScopeAllowed', 'consentRequired'
    )) {
        if ($Client.$field -cne $Template.$field) {
            throw "BFF client setting differs from the template: $field."
        }
    }
    # Keycloak omits this field when Authorization Services are disabled.
    $authorizationServices = $Client.PSObject.Properties['authorizationServicesEnabled']
    if ($authorizationServices -and $authorizationServices.Value) {
        throw 'BFF Authorization Services must be disabled.'
    }
    if (@($Client.redirectUris).Count -ne 1 -or $Client.redirectUris[0] -cne $Template.redirectUris[0] -or @($Client.webOrigins).Count -ne 0) {
        throw 'BFF requires one exact development callback and no Web Origins.'
    }
    foreach ($setting in $Template.attributes.PSObject.Properties) {
        if ($Client.attributes.($setting.Name) -cne $setting.Value) {
            throw "BFF client attribute differs from the template: $($setting.Name)."
        }
    }
    if ($Client.authenticationFlowBindingOverrides.browser -cne $Flow.id) {
        throw 'BFF must use its dedicated passwordless flow.'
    }
    if (@($Client.protocolMappers).Count -ne @($Template.protocolMappers).Count) {
        throw 'Unexpected BFF claim mappers.'
    }
    foreach ($expectedMapper in $Template.protocolMappers) {
        $matchingMappers = @($Client.protocolMappers | Where-Object name -CEQ $expectedMapper.name)
        if ($matchingMappers.Count -ne 1) {
            throw "Missing or duplicate claim mapper: $($expectedMapper.name)."
        }
        $mapper = $matchingMappers[0]
        if ($mapper.protocol -cne $expectedMapper.protocol -or $mapper.protocolMapper -cne $expectedMapper.protocolMapper -or $mapper.consentRequired -ne $expectedMapper.consentRequired) {
            throw "Incorrect claim mapper: $($expectedMapper.name)."
        }
        if (@($mapper.config.PSObject.Properties).Count -ne @($expectedMapper.config.PSObject.Properties).Count) {
            throw "Unexpected claim mapper settings: $($expectedMapper.name)."
        }
        foreach ($setting in $expectedMapper.config.PSObject.Properties) {
            if ($mapper.config.($setting.Name) -cne $setting.Value) {
                throw "Incorrect claim mapper setting: $($expectedMapper.name)/$($setting.Name)."
            }
        }
    }
    $defaultScopes = @(Invoke-Keycloak GET "$script:RealmApi/clients/$($Client.id)/default-client-scopes")
    $optionalScopes = @(Invoke-Keycloak GET "$script:RealmApi/clients/$($Client.id)/optional-client-scopes")
    if ($defaultScopes.Count -ne 0 -or $optionalScopes.Count -ne 0) {
        throw 'BFF must have no inherited default or optional/offline client scopes.'
    }
}

function Get-ProtectedConfiguration($Objects, $Policy) {
    # Compare only within this setup run; no persistent drift file or hashes.
    $realmSettings = [ordered]@{}
    foreach ($setting in $Objects.Realm.PSObject.Properties | Sort-Object Name) {
        if (@($Policy.PSObject.Properties.Name) -notcontains $setting.Name) {
            $realmSettings[$setting.Name] = $setting.Value
        }
    }
    $keys = Invoke-RestMethod "$script:KeycloakUrl/realms/ce1-dev/protocol/openid-connect/certs" -TimeoutSec 15
    $browser = @(Invoke-Keycloak GET "$script:RealmApi/authentication/flows/browser/executions")
    $otherClients = @($Objects.Clients | Where-Object clientId -CNE 'ce1-bff-dev' | Sort-Object id)
    return [pscustomobject]@{
        RealmSettings = ConvertTo-Json -InputObject $realmSettings -Depth 40 -Compress
        Roles = ConvertTo-Json -InputObject @($Objects.Roles | Sort-Object id) -Depth 20 -Compress
        DefaultBrowser = ConvertTo-Json -InputObject $browser -Depth 20 -Compress
        OtherClients = ConvertTo-Json -InputObject $otherClients -Depth 40 -Compress
        RequiredActions = ConvertTo-Json -InputObject @($Objects.RequiredActions | Sort-Object alias) -Depth 20 -Compress
        PublicKeys = ConvertTo-Json -InputObject @($keys.keys | Sort-Object kid) -Depth 20 -Compress
    }
}

function New-RandomSecret {
    $bytes = New-Object byte[] 32
    $generator = [Security.Cryptography.RandomNumberGenerator]::Create()
    try {
        $generator.GetBytes($bytes)
        return [Convert]::ToBase64String($bytes)
    } finally {
        $generator.Dispose()
    }
}

function Write-NewLocalJson([string]$Path, $Value) {
    # CreateNew refuses to truncate an existing secret, baseline or test request.
    $stream = [IO.File]::Open($Path, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write)
    try {
        $bytes = [Text.Encoding]::UTF8.GetBytes((ConvertTo-Json -InputObject $Value -Depth 30))
        $stream.Write($bytes, 0, $bytes.Length)
    } finally {
        $stream.Dispose()
    }
}
