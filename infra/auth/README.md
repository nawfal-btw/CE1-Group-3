# Primary local authentication environment

**Use `ce1-auth-local` for local development.** It initializes Keycloak directly
from the canonical realm JSON; no PowerShell initialization script is required.

Run every command below from `infra/auth`. Always include the explicit project
and environment file: the shared Compose file retains its legacy default project
name, so a bare `docker compose up` does not select the primary environment.

## Connection and storage

| Item | Primary environment |
| --- | --- |
| Compose project | `ce1-auth-local` |
| Private environment file | `local/ce1-auth-local.env` |
| Keycloak Admin Console | <http://localhost:8083/admin/> |
| Canonical Keycloak URL | `http://localhost:8083` |
| Exact OIDC issuer | `http://localhost:8083/realms/ce1-dev` |
| Realm / confidential client | `ce1-dev` / `ce1-bff-dev` |
| Reserved BFF callback | `http://localhost:8082/auth/callback` |
| WebAuthn RP ID | `localhost` |
| Host binding | `127.0.0.1:8083` -> Keycloak container port 8080 |
| PostgreSQL volume | `ce1-auth-local_postgres_data` |
| Compose network | `ce1-auth-local_auth` |

PostgreSQL and Keycloak management port 9000 have no host-published ports.
The pinned images are Keycloak `26.8.1` and PostgreSQL `17.11-bookworm`.
The callback is reserved for Phase 2B; no Go callback handler exists yet.

The future BFF must use the exact issuer above, client `ce1-bff-dev`, and the
`AUTH_BFF_CLIENT_SECRET` from the primary environment file. Its callback must
match the registered URI exactly. No backend files are configured in this phase.

## Start, stop and restart

Docker Desktop must be running with Linux containers. Verify connectivity first:

~~~sh
docker info --format '{{.OperatingSystem}} | {{.OSType}}'
~~~

Validate and start the primary project:

~~~sh
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local config --quiet
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local up -d --wait --wait-timeout 240
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local ps
~~~

Stop both primary containers, preserving their storage:

~~~sh
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local stop
~~~

Start existing stopped containers and wait for health:

~~~sh
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local start --wait --wait-timeout 240
~~~

For a restart with a readiness check, run the stop command followed by the start
command above. This is the sequence used for the persistence test.

Do not use `down --volumes`, volume prune, or database reset commands. No cleanup
of the retained environments is needed to operate this project.

## Credentials and fresh setup on another machine

This machine's primary environment file already exists. **Do not overwrite it.**
It contains an independent database password, administrator password, and BFF
client secret generated with the operating system's cryptographic random-number
generator. These values differ from both previous environments.

Open the file privately in your editor to obtain the administrator credentials.
The administrator username is `ce1-local-admin`. It is a Keycloak administrator
in the master realm, separate from the application `administrator` role.

On another machine, create a new `local/ce1-auth-local.env` using `.env.example`
as a reference and generate three new independent secrets. Do not copy another
developer's credentials. Set the following nonsecret values explicitly, because
the example and Compose fallbacks retain the historical port 8081:

~~~dotenv
AUTH_ADMIN_USERNAME=ce1-local-admin
AUTH_HTTP_PORT=8083
AUTH_KEYCLOAK_URL=http://localhost:8083
AUTH_BFF_REDIRECT_URI=http://localhost:8082/auth/callback
AUTH_WEBAUTHN_RP_ID=localhost
AUTH_SSL_REQUIRED=none
~~~

Also supply nonempty `AUTH_DB_PASSWORD`, `AUTH_ADMIN_PASSWORD` and
`AUTH_BFF_CLIENT_SECRET`. The future BFF needs the same client secret. Keep it
private; changing the environment file does not rotate an already imported secret.

The entire `local/` directory is ignored by Git. Verify that protection:

~~~sh
git check-ignore infra/auth/local/ce1-auth-local.env
~~~

Run that Git command from the repository root. Do not share secret files, resolved
Compose configuration, full Docker inspection, shell tracing/transcripts, or
credential-bearing exports. Local environment files are plaintext development
storage; Docker administrators can inspect container environment variables.

## Native realm initialization

`keycloak/ce1-dev-realm.json` is the single configuration source. It defines:

- All four application roles.
- Confidential `ce1-bff-dev` with Code Flow and mandatory PKCE S256.
- Exact callback restrictions, disabled unnecessary grants and restricted scopes.
- Required Username Form -> WebAuthn Passwordless Authenticator, bound only to
  the BFF client; the realm's default browser flow remains unchanged.
- Required user verification/discoverable credentials, ES256/RS256, compatible
  platform/roaming authenticators, registration actions and minimal claim mappings.

Compose mounts this file read-only into `/opt/keycloak/data/import` and runs
Keycloak with `--import-realm`. Keycloak resolves native environment placeholders
for the BFF secret, callback, RP ID and realm TLS policy. No host script, Python,
Job, init container or additional configuration tool initializes the realm.

A fresh database imports the definition. On restart, the existing realm is
skipped and PostgreSQL preserves configuration, users, credentials and signing
keys. Editing JSON or import-time environment values does not update an existing
realm. Use a separately reviewed administrative change rather than deleting data.

The seed contains no application users, bootstrap user passwords, private keys or
existing signing-key exports. The bootstrap master administrator comes from the
private environment file.

## Verification

The existing optional verifier supports the native-imported primary environment:

~~~powershell
powershell -NoProfile -ExecutionPolicy RemoteSigned -File .\verify-phase2a.ps1 -ImportedRealm -ServerUrl http://localhost:8083 -EnvironmentFile .\local\ce1-auth-local.env
~~~

This performs configuration readback and negative protocol tests; it does not
initialize the realm or enroll a passkey. `RemoteSigned` applies only to the child
process. Older utilities default to the historical `.env` and port 8081; use the
explicit primary command here. Do not run the legacy API setup utility as a
startup hook for an imported realm.

Optional direct discovery check:

~~~powershell
$discovery = Invoke-RestMethod 'http://localhost:8083/realms/ce1-dev/.well-known/openid-configuration'
if ($discovery.issuer -cne 'http://localhost:8083/realms/ce1-dev') {
    throw 'Unexpected issuer.'
}
~~~

Verified for this primary environment:

- Fresh native import, health checks, loopback binding and separate PostgreSQL volume.
- Exact discovery issuer and public JWKS.
- Imported BFF secret matches the private environment file, without printing it.
- Roles, claim mappings, role scopes, WebAuthn policy, required flow/binding,
  registration action and preserved default browser flow.
- Mandatory S256, exact redirects, disabled grants and required client authentication.
- Full stop/start preserves realm/client/flow IDs, settings and signing keys;
  the configuration tests pass again after restart.
- Zero application users were seeded.

These automated checks do not test browser passkey authentication. The successful
manual test is recorded below. Production TLS and Kubernetes deployment remain
untested.

## Retained environments

`ce1-auth-dev` is stopped and retained as a backup. Its containers, original
`infra/auth/.env`, local ownership/credential files, database and signing keys
remain intact. Its PostgreSQL volume is `ce1-auth-dev_postgres_data`.

`ce1-auth-import-test` remains stopped and untouched, with volume
`ce1-auth-import-test_postgres_data`. It is a separate disposable test project,
not the primary development environment.

Do not pass the primary environment file to commands that start or reconfigure
either retained project. Backup recovery must use its original settings and
credentials. Do not delete any retained containers or volumes without approval.

## Manual test results and remaining checks

Phase 2A's local identity-provider setup and positive manual browser test are
complete. The developer reported the following results in `ce1-auth-local`:

- An administrator created `dev-dispatcher` and assigned the `dispatcher` role.
- A WebAuthn passwordless passkey was registered.
- A private-window client login succeeded without entering a Keycloak password.
- Keycloak redirected to `http://localhost:8082/auth/callback` with an
  authorization code, matching state and no OAuth error.
- The callback connection was refused, as expected: the Go BFF does not exist yet.

This establishes successful browser authentication and code issuance. Code
exchange, token validation and application sessions have not been tested or
implemented. The user and passkey were provisioned manually after native import;
they are not in the seed.

To prepare another manual request, run from `infra/auth`:

~~~powershell
powershell -NoProfile -ExecutionPolicy RemoteSigned -File .\prepare-passkey-test.ps1 -EnvironmentFile .\local\ce1-auth-local.env -RequestFile ce1-auth-local-login-test-next.json
~~~

The environment parameter defaults to `local/ce1-auth-local.env`. Use a new JSON
filename on each run; existing files are never overwritten. Open the generated
file under ignored `local/`, copy only its `loginUrl` value without surrounding
quotes, and use a private browser window. The JSON preserves literal ampersands
for copying. Do not share the file, full URL or callback parameters.

Negative browser tests have not been reported as performed: cancelling passkey
authentication, attempting login before enrollment, and ensuring an existing
Account Console password-authenticated SSO session cannot bypass the client's
required passkey step. Automated configuration checks are not substitutes for
these browser tests. Cross-platform authenticator compatibility is also untested.

Removal of any enrollment password, termination of enrollment sessions and a
fresh login after cleanup remain unconfirmed. No cleanup is automated or claimed.
Keep bootstrap credentials private; a shared OS-profile PIN does not establish
individual identity. The primary Account Console is
<http://localhost:8083/realms/ce1-dev/account/>.

See [PHASE2A.md](PHASE2A.md) for repeat testing, enrollment and security limits.
BFF callbacks, application sessions and device authentication belong to later,
separately approved work.

## Troubleshooting and limitations

- If Docker is unavailable, start Docker Desktop with Linux containers.
- If startup fails, inspect `docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local logs --tail 100 keycloak postgres` privately.
- If port 8083 is occupied, identify its owner; do not stop an unrelated service.
- If credentials fail after an environment-file edit, remember that stored
  database/admin/client credentials do not rotate through startup settings.
- If configuration changes appear ignored, the existing realm was skipped.
- Always use `localhost` consistently for browser access and issuer validation.

This local setup uses `start-dev`, HTTP, `sslRequired: none`, and unencrypted
database traffic inside its isolated bridge network. Do not expose it on
LAN/Tailscale or use real incident data. A later Kubernetes deployment needs
production startup, HTTPS, trusted hostname/proxy settings, private access,
managed secrets, appropriate database security and persistence. It can use this
same JSON without PowerShell. No Kubernetes manifests are created here.

## Official references

- [Realm import and placeholders](https://www.keycloak.org/server/importExport)
- [Hostname configuration](https://www.keycloak.org/server/hostname)
- [Production configuration](https://www.keycloak.org/server/configuration-production)
