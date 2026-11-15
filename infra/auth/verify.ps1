# Read-only OIDC checks. Does not obtain tokens or print key material.
[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
$issuer = 'http://localhost:8081/realms/ce1-dev'
$discovery = Invoke-RestMethod -Uri "$issuer/.well-known/openid-configuration" -TimeoutSec 15
if ($discovery.issuer -cne $issuer) {
    throw "Unexpected issuer: $($discovery.issuer)"
}
if ($discovery.jwks_uri -cne "$issuer/protocol/openid-connect/certs") {
    throw 'Discovery did not advertise the expected local JWKS endpoint.'
}
if ($discovery.authorization_endpoint -cne "$issuer/protocol/openid-connect/auth" -or
    $discovery.token_endpoint -cne "$issuer/protocol/openid-connect/token") {
    throw 'Discovery did not advertise the expected local authorization/token endpoints.'
}
if (@($discovery.code_challenge_methods_supported) -cnotcontains 'S256') {
    throw 'Discovery did not advertise PKCE S256 support.'
}

$jwks = Invoke-RestMethod -Uri $discovery.jwks_uri -TimeoutSec 15
$signingKeys = @($jwks.keys | Where-Object {
    $_.use -eq 'sig' -and $_.kid -and $_.kty -eq 'RSA' -and $_.n -and $_.e
})
if ($signingKeys.Count -lt 1) {
    throw 'JWKS contains no RSA signing key with a key ID and public parameters.'
}
foreach ($key in @($jwks.keys)) {
    foreach ($privateField in @('d', 'p', 'q', 'dp', 'dq', 'qi', 'oth', 'k')) {
        if ($key.PSObject.Properties.Name -contains $privateField) {
            throw 'JWKS unexpectedly exposes private or symmetric key material.'
        }
    }
}

$missingRealmStatus = $null
try {
    $null = Invoke-WebRequest -UseBasicParsing -Uri 'http://localhost:8081/realms/ce1-phase1-missing/.well-known/openid-configuration' -TimeoutSec 15
} catch {
    if ($null -eq $_.Exception.Response) { throw }
    $missingRealmStatus = [int]$_.Exception.Response.StatusCode
}
if ($missingRealmStatus -ne 404) {
    throw 'A nonexistent realm did not return HTTP 404.'
}

Write-Output "PASS: discovery issuer/endpoints, PKCE S256, public JWKS ($($signingKeys.Count) RSA signing key(s)), and nonexistent-realm HTTP 404."
