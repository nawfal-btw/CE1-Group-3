# Prepares a manual browser authorization request. Does not open a browser,
# handle a callback, exchange a code, create a user or authenticate anyone.
[CmdletBinding()]
param(
    [ValidatePattern('^[a-zA-Z0-9_-]+\.json$')]
    [string]$RequestFile = 'phase2a-login-test.json',
    [string]$EnvironmentFile
)
. "$PSScriptRoot/phase2a-common.ps1"
if (!$EnvironmentFile) {
    $EnvironmentFile = Join-Path $PSScriptRoot 'local/ce1-auth-local.env'
}
$script:AuthSettings = Read-AuthSettings $EnvironmentFile
$definition = Get-AuthDefinition
if ($script:AuthSettings['AUTH_KEYCLOAK_URL']) {
    $script:KeycloakUrl = $script:AuthSettings['AUTH_KEYCLOAK_URL'].TrimEnd('/')
}
$verifier = (New-RandomSecret).TrimEnd('=').Replace('+', '-').Replace('/', '_')
$sha256 = [Security.Cryptography.SHA256]::Create()
try {
    $challengeBytes = $sha256.ComputeHash([Text.Encoding]::ASCII.GetBytes($verifier))
    $challenge = [Convert]::ToBase64String($challengeBytes).TrimEnd('=').Replace('+', '-').Replace('/', '_')
} finally {
    $sha256.Dispose()
}
$state = [guid]::NewGuid().ToString('N')
$nonce = [guid]::NewGuid().ToString('N')
$parameters = [ordered]@{
    client_id = 'ce1-bff-dev'
    redirect_uri = $definition.Client.redirectUris[0]
    response_type = 'code'
    scope = 'openid'
    code_challenge = $challenge
    code_challenge_method = 'S256'
    state = $state
    nonce = $nonce
}
$query = @($parameters.Keys | ForEach-Object { [uri]::EscapeDataString($_)+'='+[uri]::EscapeDataString($parameters[$_]) }) -join '&'
$localDirectory = Join-Path $PSScriptRoot 'local'
if (!(Test-Path -LiteralPath $localDirectory)) {
    throw 'Run configure-phase2a.ps1 -Apply first.'
}
$requestPath = Join-Path $localDirectory $RequestFile
if (Test-Path -LiteralPath $requestPath) {
    throw 'Request file already exists. Choose another -RequestFile name; existing local files are never overwritten.'
}
$request = @{
    createdUtc = [DateTime]::UtcNow.ToString('o')
    state = $state
    nonce = $nonce
    codeVerifier = $verifier
    loginUrl = "$script:KeycloakUrl/realms/ce1-dev/protocol/openid-connect/auth?$query"
}
# Ampersands are valid JSON characters. Keep them readable when copying the URL
# directly from the file in Windows PowerShell, which otherwise escapes them.
$json = (ConvertTo-Json -InputObject $request).Replace('\u0026', '&')
$stream = [IO.File]::Open($requestPath, [IO.FileMode]::CreateNew, [IO.FileAccess]::Write)
try {
    $bytes = [Text.Encoding]::UTF8.GetBytes($json)
    $stream.Write($bytes, 0, $bytes.Length)
} finally {
    $stream.Dispose()
}
Write-Output "Manual request saved to ignored local/$RequestFile. Open its loginUrl yourself; no browser was launched."
