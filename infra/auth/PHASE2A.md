# Phase 2A: OIDC client and passwordless configuration

Keycloak 26.8.1 natively imports the complete configuration from
`keycloak/ce1-dev-realm.json` on a fresh database. Existing realms are skipped.
PowerShell is optional for host-side development and verification, never a server
or Kubernetes startup requirement. This phase does not create a callback handler, session, device
credential, mTLS configuration, or backend authorization policy. All commands
below run in PowerShell from `infra/auth`.

## Completion status

The primary environment is `ce1-auth-local` at **http://localhost:8083**, with
issuer `http://localhost:8083/realms/ce1-dev` and private settings in
`local/ce1-auth-local.env`. The old `ce1-auth-dev` and disposable
`ce1-auth-import-test` environments are stopped and preserved. See the
[README](README.md) for explicit Compose project commands; do not use bare
Compose defaults to select the primary environment.

Previously completed automated checks passed for fresh native import, client
configuration/secret matching, role scopes and claim mappings, WebAuthn policy
and flow binding, registration actions, preserved default browser flow, PKCE,
redirect restrictions, disabled grants and persistence across stop/start.
The subsequent read-only review also confirmed healthy primary containers,
the exact discovery issuer and a reachable JWKS endpoint.

The developer confirmed successful manual enrollment and private-window
passwordless login as `dev-dispatcher` with the `dispatcher` role. No Keycloak
password was entered during client login. Keycloak returned an authorization-code
redirect with matching state and no OAuth error to
`http://localhost:8082/auth/callback`. Connection refusal was expected because
there is no BFF callback yet. This is a manual browser result, not an automated
test of code exchange, token validation or sessions.

Not yet reported as tested: browser cancellation, login before passkey enrollment,
and resistance to bypass using an existing Account Console password-authenticated
SSO session. Cross-platform passkeys, production TLS and Kubernetes remain
untested. Enrollment-password removal and termination of enrollment sessions
remain unconfirmed. Phase 2B implementation has not begun.

## Application identities

The workflow is **Dispatcher -> Agency (Operator) -> Agent (Responder)**.
The four existing realm roles are preserved without renaming or changing their
descriptions or permissions:

| Application role | Intended responsibility |
| --- | --- |
| `administrator` | Manage application users, role assignments, agency associations and devices; no automatic Keycloak administration privilege |
| `dispatcher` | Receive incidents, coordinate responses and request assistance from agencies |
| `operator` | Individual human coordinating an agency's resources and agents |
| `responder` | Agent in the field; initial authentication will use device mTLS in a later phase |

Agency is an organization. Agent is a domain entity representing a responder.
Resources and vehicles are operational entities. None of these is a substitute
for an individual identity. The eventual backend must enforce both roles and
agency/resource-specific permissions.

## Native initialization and optional local utilities

Follow the [README](README.md) for native startup, placeholders and verification.
The seed includes the client, role scopes, custom flow, its
client-specific binding, policies and required actions. It includes no users or
credentials. Keep the native default browser flow; do not bind the custom flow
realm-wide. Registration actions include the standard 26.8.1 actions because an
explicit import list replaces defaults.

For optional verification, first confirm both services are healthy. The selected
environment file's administrator settings must match the existing administrator;
changing bootstrap environment variables does not change persisted credentials.

```powershell
# Primary native-imported realm; no local ownership file required.
powershell -NoProfile -ExecutionPolicy RemoteSigned -File .\verify-phase2a.ps1 -ImportedRealm -ServerUrl http://localhost:8083 -EnvironmentFile .\local\ce1-auth-local.env
```

`-ImportedRealm` checks configuration and protocol restrictions but skips the
historical fixture's user/credential assertions. Its manual-test reminder reports
what the script itself did not test; it does not invalidate the manual result
recorded above. Do not omit the explicit primary server/environment settings.
The unchanged Phase 1 `verify.ps1` targets historical port 8081.

`RemoteSigned` applies only to each child process. No persistent execution policy
is changed. The scripts authenticate the existing master administrator to obtain
an in-memory Admin API token; that administrative token is unrelated to dispatcher
login. No credential values or token response bodies are printed. Do not enable
PowerShell tracing/transcription, share terminal dumps, or print local secret files.

`configure-phase2a.ps1 -Apply` remains an optional Admin API utility for the
original, owned local development setup. It reads client and policy settings from
the canonical realm JSON. Do not run it as a hook after native import: imported
objects intentionally lack its local ownership state, so it refuses to adopt
them. It never imports over a realm, deletes a database, rotates secrets/keys or
resets an existing user's password. Applying live changes still requires approval.

The original API-based setup keeps these **ignored historical local files**:

- `local/phase2a-secrets.json`: client secret, enrollment bootstrap password,
  setup ownership marker and recorded development-user ID.
- `local/phase2a-baseline.json`: a historical artifact from the initial setup.
  Refactored scripts leave it untouched and no longer read or create it.
- `local/phase2a-login-test.json`: generated only by the manual-test preparation
  command; contains a PKCE verifier and request details. Treat it as sensitive.

The primary environment uses `local/ce1-auth-local.env` and separately named
`local/ce1-auth-local-login-test*.json` requests. Do not reuse historical fixture
credentials for the primary environment. All of `local/` is ignored by Git.

Client and flow ownership markers and the recorded user ID prevent accidentally
adopting unrelated objects with matching names. Missing local state, an unexpected
owned configuration, or a changed secret causes a failure. Existing secret files
are read only; they are never rewritten to accept a mismatch. Keep local files
protected by your Windows account permissions. They are plaintext development
storage, not a production secret manager.

Setup compares unrelated realm settings, roles, default browser executions,
built-in clients, required actions and public signing keys immediately before and
after its work. These comparisons stay in memory and report which section changed.
Legitimate administrative changes between separate runs do not require updating a
permanent drift file. Routine verification checks the owned client, policy, flow,
roles and credential state.

The optional API utility creates its local secret file once after successful setup.
Setup is not transactional across Admin API calls: interruption can leave partial
objects that require deliberate administrator recovery. Reruns refuse to adopt
objects without matching local state, rather than guessing, resetting credentials
or deleting objects. Successful setup is repeatable with zero writes. Existing
passwords/passkeys are never replaced, and completed registration actions are not
re-added. Neither the historical baseline nor the in-memory checks are a backup.

The complete realm JSON is the only authentication configuration definition;
the old standalone client/policy fragments have been removed. Keycloak resolves
native placeholders from container environment variables. Optional host scripts
read the same file and resolve its documented settings for readback comparisons.
Editing it does not update a persisted realm. Changing an import-time secret does
not rotate the stored secret, and changing the RP ID does not migrate passkeys.

## What the client does

`ce1-bff-dev` is confidential: the future server keeps its secret and presents it
to the token endpoint (`client_secret_basic` is the intended integration method).
The browser never receives that secret. Authorization Code Flow is enabled and
PKCE **S256 is mandatory**. A future login transaction generates a verifier and
sends its SHA-256 challenge; code exchange requires the original verifier as well
as confidential-client authentication.

The local default callback is **http://localhost:8082/auth/callback**, supplied by
`AUTH_BFF_REDIRECT_URI`. It is reserved for Phase 2B and does not exist yet. Deployed
environments must supply their exact HTTPS callback. No wildcard redirects, CORS origins, offline
scope, service accounts, implicit/password grants, device authorization, CIBA, or
standard token exchange are enabled for this client. Refreshing a normal online
session remains possible; offline refresh permissions are not assigned.

No inherited default or optional client scopes are attached. Minimal explicit
mappers expose `preferred_username` in the ID token/UserInfo and scoped realm
roles in access tokens. Role scope contains only `administrator`, `dispatcher`,
and `operator`; `responder` and realm-management privileges are excluded. Scope
mapping limits claims; it does **not** prevent another realm user from logging in.
The future BFF must deny users lacking an allowed application role. API audiences
and complete token validation belong to later integration. ID tokens are not API
access tokens.

## Passwordless policy and flow

The client's browser override points to `ce1-bff-passkey`:

1. **Username Form: REQUIRED** -- selects the individual account without a password.
2. **WebAuthn Passwordless Authenticator: REQUIRED** -- verifies an existing passkey.

There is no Cookie, Password Form, conditional credential skip, or alternative
authentication execution. The realm's default `browser` flow is unchanged, so the
Account Console remains available for controlled initial enrollment. A missing
passkey must not become username-only login. Enrollment uses a separate,
authenticated bootstrap procedure below.

The passwordless registration action is `webauthn-register-passwordless`, enabled
but not a realm-wide default. The seed assigns it to no users; an administrator
assigns it during enrollment. The local test account is `dev-dispatcher`.
After registration Keycloak normally removes the completed action.

Policy: RP ID supplied by `AUTH_WEBAUTHN_RP_ID` (**localhost** locally), user verification **required**, discoverable
credentials **required**, algorithms **ES256/RS256**, attachment unrestricted,
attestation **none**, and no extra origins. User verification means a local PIN or
biometric unlocks the authenticator; Keycloak does not receive that PIN/biometric.
This uses the current Discoverable Credential setting, not its deprecated
replacement compatibility flag.

No operating system is required. Compatible Windows Hello, Apple/Android passkey
implementations, Linux FIDO2 authenticators and roaming security keys may work;
actual support depends on browser and authenticator. No cross-platform success is
claimed without testing. Localhost HTTP is a development browser exception, not a
deployment pattern; production requires HTTPS and a stable hostname. Credentials
registered for localhost do not automatically transfer to a production RP ID.

## Manual enrollment and repeat browser testing

No script attempts browser enrollment or removes the bootstrap password.

Native import does not create `dev-dispatcher` or its password. On a fresh deployment, an
administrator provisions each human account deliberately, assigns its application
role and `Webauthn Register Passwordless` required action, and handles initial
enrollment/recovery privately. Do not copy local fixture passwords into deployment
configuration. Individual users must operate their own authenticators.

The primary account is already enrolled and its positive login has passed. Do
not recreate it or reset credentials to repeat that test; start at step 4.

1. For a new account only, have an administrator provision an individual user,
   assign the approved application role and enable `Webauthn Register Passwordless`.
   Arrange enrollment credentials privately; do not use the historical
   `local/phase2a-secrets.json` in this environment.
2. In a separate browser profile/private window from the administrator, open
   **http://localhost:8083/realms/ce1-dev/account/**. Sign in using the enrollment
   password. If a profile prompt requires email, use disposable
   `dev-dispatcher@example.invalid`; no SMTP/email verification is configured.
3. Complete the required **Webauthn Register Passwordless** action using an
   authenticator you independently control. Give it a recognizable label.
4. Log out of the Account Console. Prepare a fresh manual client request:

   ```powershell
   powershell -NoProfile -ExecutionPolicy RemoteSigned -File .\prepare-passkey-test.ps1 -EnvironmentFile .\local\ce1-auth-local.env -RequestFile ce1-auth-local-login-test-next.json
   ```

   The environment parameter defaults to the primary file. Existing request files
   are preserved; choose a fresh request filename for every repeat test. The
   filename must be a JSON filename without directory components; it stays under
   the ignored `local/` directory.

5. Open the generated ignored request file in your editor and copy its `loginUrl`
   value without surrounding quotes into a private browser window. JSON output
   keeps literal `&` separators for direct copying. Enter `dev-dispatcher`, then
   authenticate using your passkey. No password fallback should appear.
6. After successful authentication, the browser should navigate to the exact
   reserved callback with `code`, `state`, and issuer information. Connection
   refusal is expected because the BFF is absent. Compare the returned state
   privately with the saved request. Do not share the URL/code or treat this
   observation as token/session validation.
7. The following negative tests remain pending: cancel passkey authentication,
   and try the client request while an Account Console password-authenticated SSO
   session exists: neither should bypass the required passkey execution. Test
   pre-enrollment rejection only with a separately approved unenrolled account;
   it must not offer unauthenticated registration. Do not remove an existing
   passkey or reset the primary account to perform this test.
8. Rerun the explicit primary verification command above for configuration checks.
   Imported-realm mode does not assert user credentials. Report manual results
   without sharing credentials. **Retain the bootstrap
   password until successful passkey authentication is confirmed.**

Historically, the optional API setup used Keycloak `temporary=false` for its
enrollment-only password to avoid forcing a second password before registration.
Do not assume the manually provisioned primary account has the same setting.
A retained non-expiring password can still work through the
unchanged default Account Console flow. After manual success,
explicitly remove the user's Password credential in the Admin Console, end the
enrollment sessions, and verify a fresh client passkey login. Separately remove
the retained local bootstrap value when reviewing credential cleanup; do not
delete the client secret or ownership state. Setup does not do any of this for you.
Cleanup in the primary environment has not been confirmed; no credential changes
are authorized as part of documentation closeout.

## Shared workstations and limitations

Each human requires a separate identity and independently controlled authenticator.
A shared workstation PIN proves neither individual identity nor passkey ownership.
Multiple users' passkeys stored under one shared OS profile are not assumed to be
independently protected. For shared workstations prefer independently controlled
security keys or personal passkeys, or properly isolated OS user profiles.
Synced passkeys are allowed here; none establishes an authorized device identity.

Application logout, session expiration and user switching belong to the BFF
phase. Device binding/mTLS and agency/resource authorization are later work.
No responder account is created. `dev-dispatcher` has only the `dispatcher`
application role plus ordinary Keycloak defaults; it has no realm-administration
privileges. Other human test accounts need separate approval.

Noninteractive verification tests readback, exact flow/scopes,
mandatory PKCE, redirect restrictions, disabled grants and confidential client
authentication. They do not prove passkey enrollment, real user authentication,
cross-platform compatibility, code exchange, JWT validation or secure sessions.
Phase 2B will implement the Go BFF's login/callback, PKCE/state/nonce validation,
token handling, application-role admission, and secure server-side sessions,
logout and user switching, subject to the next approved implementation scope.

## References

- [26.8.1 Admin REST API](https://www.keycloak.org/docs-api/26.8.1/rest-api/index.html)
- [26.8.1 WebAuthn documentation](https://github.com/keycloak/keycloak/blob/26.8.1/docs/documentation/server_admin/topics/authentication/webauthn.adoc)
- [26.8.1 authentication flows](https://github.com/keycloak/keycloak/blob/26.8.1/docs/documentation/server_admin/topics/authentication/flows.adoc)
- [Passwordless provider forbids user setup during authentication](https://github.com/keycloak/keycloak/blob/26.8.1/services/src/main/java/org/keycloak/authentication/authenticators/browser/WebAuthnPasswordlessAuthenticatorFactory.java)
- [OAuth security best practice](https://www.rfc-editor.org/rfc/rfc9700.html)
