[CmdletBinding()]
param(
    [string]$ServerUrl = 'http://localhost:8081',
    [string]$EnvironmentFile,
    [switch]$ImportedRealm
)

. "$PSScriptRoot/phase2a-common.ps1"
if (!$EnvironmentFile) {
    $EnvironmentFile = Join-Path $PSScriptRoot '.env'
}
$script:KeycloakUrl = $ServerUrl.TrimEnd('/')
Connect-Keycloak $EnvironmentFile
$definition = Get-AuthDefinition
$policy = $definition.Policy
$clientTemplate = $definition.Client
$objects = Get-AuthObjects
$passkeys = @()
if (!$objects.Client -or !$objects.Flow) {
    throw 'BFF client or passwordless flow is missing.'
}
if ($ImportedRealm) {
    if (!$script:AuthSettings['AUTH_BFF_CLIENT_SECRET']) {
        throw 'Imported-realm verification requires AUTH_BFF_CLIENT_SECRET in the environment file.'
    }
    $localState = [pscustomobject]@{ clientSecret = $script:AuthSettings['AUTH_BFF_CLIENT_SECRET'] }
} else {
    $localState = Get-Content -LiteralPath "$PSScriptRoot/local/phase2a-secrets.json" -Raw | ConvertFrom-Json
    Assert-ObjectOwnership $objects $localState
}

foreach ($role in $definition.Realm.roles.realm) {
    if (@($objects.Roles | Where-Object name -CEQ $role.name).Count -ne 1) {
        throw "Missing application role: $($role.name)."
    }
}
if ($objects.Realm.browserFlow -cne 'browser') {
    throw 'The realm default browser flow was changed.'
}
$browserExecutions = @(Invoke-Keycloak GET "$script:RealmApi/authentication/flows/browser/executions")
if (@($browserExecutions | Where-Object { $_.PSObject.Properties['providerId'] -and $_.providerId -ceq 'auth-cookie' }).Count -ne 1 -or
    @($browserExecutions | Where-Object { $_.PSObject.Properties['providerId'] -and $_.providerId -ceq 'auth-username-password-form' }).Count -ne 1) {
    throw 'The default browser flow is missing its normal Cookie or Username Password step.'
}
$policyDifferences = @(Get-PolicyDifferences $objects.Realm $policy)
if ($policyDifferences.Count -gt 0) {
    throw "Incorrect passwordless policy: $($policyDifferences -join ', ')."
}
if ($objects.Realm.sslRequired -cne $definition.Realm.sslRequired) {
    throw 'Realm TLS policy does not match environment settings.'
}
Assert-PasswordlessFlow $objects.Executions
Assert-BffClient $objects.Client $objects.Flow $clientTemplate
$registrationAction = @($objects.RequiredActions | Where-Object alias -CEQ 'webauthn-register-passwordless')
if ($registrationAction.Count -ne 1 -or !$registrationAction[0].enabled -or $registrationAction[0].defaultAction) {
    throw 'Passwordless registration must be enabled but not a realm-wide default.'
}
$roleScope = @(Invoke-Keycloak GET "$script:RealmApi/clients/$($objects.Client.id)/scope-mappings/realm/composite")
$roleNames = @($roleScope | ForEach-Object name | Sort-Object)
$expectedRoles = @($definition.RoleScope | Sort-Object)
if (($roleNames -join ',') -cne ($expectedRoles -join ',')) {
    throw 'BFF role scope must contain only the three human application roles.'
}
$clientSecret = Invoke-Keycloak GET "$script:RealmApi/clients/$($objects.Client.id)/client-secret"
if (!$localState.clientSecret -or $clientSecret.value -cne $localState.clientSecret) {
    throw 'BFF client secret differs from the supplied secret.'
}

if (!$ImportedRealm) {
    # Only the original local development utility owns this disposable account.
    $user = $objects.User
    if (!$user.enabled) {
        throw 'Development dispatcher is disabled.'
    }
    $mapping = Invoke-Keycloak GET "$script:RealmApi/users/$($user.id)/role-mappings"
    $effectiveRoles = @(Invoke-Keycloak GET "$script:RealmApi/users/$($user.id)/role-mappings/realm/composite")
    $applicationRoles = @($effectiveRoles | Where-Object { $_.name -cin @('administrator', 'dispatcher', 'operator', 'responder') })
    if ($applicationRoles.Count -ne 1 -or $applicationRoles[0].name -cne 'dispatcher') {
        throw 'The development dispatcher must have only the dispatcher application role.'
    }
    $clientMappings = $mapping.PSObject.Properties['clientMappings']
    if ($clientMappings -and $clientMappings.Value.PSObject.Properties['realm-management']) {
        throw 'Dispatcher has realm-management assignments.'
    }
    $directRoles = @(Invoke-Keycloak GET "$script:RealmApi/users/$($user.id)/role-mappings/realm")
    if (@($directRoles | Where-Object { $_.name -cnotin @('dispatcher', 'default-roles-ce1-dev') }).Count -ne 0) {
        throw 'Unexpected direct realm role assignment.'
    }
    $credentials = @(Invoke-Keycloak GET "$script:RealmApi/users/$($user.id)/credentials")
    $passkeys = @($credentials | Where-Object type -CEQ 'webauthn-passwordless')
    if ($passkeys.Count -eq 0 -and @($user.requiredActions) -cnotcontains 'webauthn-register-passwordless') {
        throw 'Unenrolled dispatcher is missing its registration action.'
    }
}
Write-Output 'PASS: roles, client/secret, claim mappings, role scope, passwordless policy/flow, registration action and default browser flow.'

# One client for the entire run. Do not follow redirects or carry login cookies.
Add-Type -AssemblyName System.Net.Http
$handler = New-Object System.Net.Http.HttpClientHandler
$handler.AllowAutoRedirect = $false
$handler.UseCookies = $false
$httpClient = New-Object System.Net.Http.HttpClient($handler)
$httpClient.Timeout = [TimeSpan]::FromSeconds(20)

function Invoke-OidcProbe([string]$Path, $Values, [switch]$Post) {
    $form = @($Values.Keys | ForEach-Object {
        [uri]::EscapeDataString([string]$_) + '=' + [uri]::EscapeDataString([string]$Values[$_])
    }) -join '&'
    $request = New-Object System.Net.Http.HttpRequestMessage
    $response = $null
    try {
        if ($Post) {
            $request.Method = [System.Net.Http.HttpMethod]::Post
            $request.RequestUri = "$script:KeycloakUrl$Path"
            $request.Content = New-Object System.Net.Http.StringContent($form, [Text.Encoding]::UTF8, 'application/x-www-form-urlencoded')
        } else {
            $request.Method = [System.Net.Http.HttpMethod]::Get
            $request.RequestUri = "$script:KeycloakUrl$Path" + '?' + $form
        }
        $response = $httpClient.SendAsync($request).GetAwaiter().GetResult()
        return @{
            status = [int]$response.StatusCode
            body = $response.Content.ReadAsStringAsync().GetAwaiter().GetResult()
            location = [string]$response.Headers.Location
        }
    } finally {
        if ($response) {
            $response.Dispose()
        }
        $request.Dispose()
    }
}

function New-AuthorizationParameters {
    return [ordered]@{
        client_id = 'ce1-bff-dev'
        redirect_uri = $clientTemplate.redirectUris[0]
        response_type = 'code'
        scope = 'openid'
        state = [guid]::NewGuid().ToString('N')
        nonce = [guid]::NewGuid().ToString('N')
        code_challenge = ('a' * 43)
        code_challenge_method = 'S256'
    }
}

function Assert-AuthorizationRejected($Parameters, [string]$Label) {
    $response = Invoke-OidcProbe '/realms/ce1-dev/protocol/openid-connect/auth' $Parameters
    $rejected = ($response.status -eq 400) -or
        ($response.status -eq 302 -and $response.location -match '[?#&]error=')
    if (!$rejected) {
        throw "Protocol test failed: $Label was not rejected (HTTP $($response.status))."
    }
    if ($response.location -and $response.location -notmatch ('^' + [regex]::Escape($clientTemplate.redirectUris[0]) + '[?#]')) {
        throw 'Error redirected outside the approved callback.'
    }
    Write-Output "PASS: $Label rejected."
}

try {
    $parameters = New-AuthorizationParameters
    $response = Invoke-OidcProbe '/realms/ce1-dev/protocol/openid-connect/auth' $parameters
    if ($response.status -ne 200 -or
        $response.body -notmatch 'name="username"' -or
        $response.body -match 'type="password"') {
        throw 'Valid code/S256 request did not reach the username-only form.'
    }
    Write-Output 'PASS: valid code/S256 request reaches username-only authentication.'

    $parameters = New-AuthorizationParameters
    $parameters.Remove('code_challenge')
    $parameters.Remove('code_challenge_method')
    Assert-AuthorizationRejected $parameters 'missing PKCE'

    $parameters = New-AuthorizationParameters
    $parameters.code_challenge_method = 'plain'
    Assert-AuthorizationRejected $parameters 'plain PKCE'

    foreach ($redirectUri in @(
        ($clientTemplate.redirectUris[0] + '/'),
        'http://localhost:8082/auth/other',
        'http://127.0.0.1:8082/auth/callback',
        'http://example.invalid/auth/callback'
    )) {
        $parameters = New-AuthorizationParameters
        $parameters.redirect_uri = $redirectUri
        Assert-AuthorizationRejected $parameters 'nonmatching redirect'
    }
    $parameters = New-AuthorizationParameters
    $parameters.response_type = 'token'
    Assert-AuthorizationRejected $parameters 'implicit flow'

    $tokenPath = '/realms/ce1-dev/protocol/openid-connect/token'
    foreach ($grantType in @('password', 'client_credentials', 'urn:ietf:params:oauth:grant-type:token-exchange')) {
        $grantParameters = @{
            client_id = 'ce1-bff-dev'
            client_secret = $localState.clientSecret
            grant_type = $grantType
        }
        if ($grantType -eq 'urn:ietf:params:oauth:grant-type:token-exchange') {
            $grantParameters.subject_token = 'deliberately-invalid-test-token'
            $grantParameters.subject_token_type = 'urn:ietf:params:oauth:token-type:access_token'
        }
        $response = Invoke-OidcProbe $tokenPath $grantParameters -Post
        if ($response.status -notin @(400, 401)) {
            throw 'Disabled grant unexpectedly succeeded.'
        }
        $oauthError = $response.body | ConvertFrom-Json
        $capabilityDenied = $oauthError.error -in @('unauthorized_client', 'unsupported_grant_type')
        if ($grantType -eq 'urn:ietf:params:oauth:grant-type:token-exchange') {
            # This version reports a disabled exchange as invalid_request.
            $capabilityDenied = $oauthError.error -ceq 'invalid_request' -and
                $oauthError.error_description -ceq 'Standard token exchange is not enabled for the requested client'
        }
        if (!$capabilityDenied) {
            throw "Grant $grantType was not rejected for the expected capability restriction ($($oauthError.error))."
        }
        Write-Output "PASS: $grantType disabled."
    }

    foreach ($grantPath in @(
        '/realms/ce1-dev/protocol/openid-connect/auth/device',
        '/realms/ce1-dev/protocol/openid-connect/ext/ciba/auth'
    )) {
        $response = Invoke-OidcProbe $grantPath @{
            client_id = 'ce1-bff-dev'
            client_secret = $localState.clientSecret
            scope = 'openid'
            login_hint = 'dev-dispatcher'
        } -Post
        if ($response.status -notin @(400, 401)) {
            throw 'Device/CIBA grant unexpectedly accepted.'
        }
        $oauthError = $response.body | ConvertFrom-Json
        $capabilityDenied = $oauthError.error -ceq 'unauthorized_client'
        if ($grantPath.EndsWith('/ext/ciba/auth')) {
            $capabilityDenied = $oauthError.error -ceq 'invalid_grant' -and
                $oauthError.error_description -ceq 'Client not allowed OIDC CIBA Grant'
        }
        if (!$capabilityDenied) {
            throw 'Device/CIBA grant was not rejected by the expected client restriction.'
        }
    }
    Write-Output 'PASS: device authorization and CIBA disabled.'

    $response = Invoke-OidcProbe $tokenPath @{
        client_id = 'ce1-bff-dev'
        grant_type = 'authorization_code'
        code = 'deliberately-invalid-test-code'
        redirect_uri = $clientTemplate.redirectUris[0]
        code_verifier = ('a' * 43)
    } -Post
    $oauthError = $response.body | ConvertFrom-Json
    if ($response.status -ne 401 -or $oauthError.error -cne 'unauthorized_client' -or
        $oauthError.error_description -cne 'Invalid client or Invalid client credentials') {
        throw 'Token endpoint did not reject missing confidential-client credentials as expected.'
    }
    Write-Output 'PASS: unauthenticated token request rejected.'
    if ($ImportedRealm) {
        Write-Output 'MANUAL: user provisioning and passkey enrollment are intentionally excluded from the import seed. No browser login was tested.'
    } else {
        Write-Output "MANUAL: $($passkeys.Count) passwordless credential(s) registered. Browser enrollment and login are not verified by this script."
    }
    if ($script:AdminWrites -ne 0) {
        throw 'Verification unexpectedly modified configuration.'
    }
} finally {
    $httpClient.Dispose()
}
