# Local authentication environment

Phases 2A, independent 2B implementation and local 2C browser testing are complete.
Shared API/router integration is deferred. The temporary `authdev` harness reuses
the authentication components; it is excluded from ordinary builds by its
`authdev` build tag and must not be deployed as the application API.

## Current topology

| Component | Configuration |
| --- | --- |
| Compose project | `ce1-auth-local` (always pass explicitly) |
| Private environment file | `infra/auth/local/ce1-auth-local.env` |
| Keycloak / Admin Console | `https://localhost:8443` / `https://localhost:8443/admin/` |
| Exact issuer | `https://localhost:8443/realms/ce1-dev` |
| Realm / confidential client | `ce1-dev` / `ce1-bff-dev` |
| Harness origin | `https://127.0.0.1:8444` |
| Exact callback | `https://127.0.0.1:8444/auth/callback` |
| WebAuthn RP ID | `localhost` |
| Application session fixture | Database/user `ce1_auth_test` at `127.0.0.1:55432` |
| Keycloak PostgreSQL volume | `ce1-auth-local_postgres_data` |
| Compose network | `ce1-auth-local_auth` |

Keycloak 26.8.1 and PostgreSQL 17.11-bookworm are pinned. Keycloak publishes only
`127.0.0.1:8443`; its main HTTP listener is disabled. Management HTTP port 9000
and Keycloak PostgreSQL have no host publication. Application sessions use the
separate `ce1-auth-persistence-test` fixture and
`ce1-auth-persistence-test-data` volume, never Keycloak's database.

Keep Keycloak on `localhost` to preserve the existing WebAuthn RP ID. Browse the
harness using `127.0.0.1`: different hostnames isolate host-only cookies, whereas
different ports alone do not. Direct Go TLS is used; no reverse proxy or trusted
forwarded headers are configured.

## Certificates and private settings

The verified local setup uses mkcert with a CA trusted in the development user's
Windows certificate store. Certificate preparation and trust installation are
separate prerequisites on another machine; never copy private CA keys.

- `local/tls/keycloak.pem` and `keycloak-key.pem`: DNS SAN `localhost`.
- `local/tls/authdev.pem` and `authdev-key.pem`: IP SAN `127.0.0.1`.
- `local/tls/rootCA.pem`: public CA certificate, explicitly loaded by Go.

Restrict private-key directory/file permissions. Keycloak mounts only its leaf
certificate and key read-only, never the CA private key. Missing mount files fail
rather than creating directories. Verify browser trust without certificate
bypasses; Firefox trust can depend on its configuration. Windows curl/Schannel
may report unknown revocation status for this local CA; do not use `-k`.

The existing private environment file must not be overwritten. For a fresh
machine, use `.env.example` as a template in ignored `local/` and supply independent,
nonempty `AUTH_DB_PASSWORD`, `AUTH_ADMIN_PASSWORD` and `AUTH_BFF_CLIENT_SECRET`.
The base file/example retain historical HTTP defaults; always combine the HTTPS
override for the current topology. It supplies the HTTPS callback, RP ID and
realm SSL policy for a fresh import. Changing environment values does not rotate
persisted credentials or update an existing realm.

Never print private environment files, resolved Compose configuration, raw
authentication errors, callback query strings, cookies or tokens. The entire
`infra/auth/local/` directory and private auth environment files are Git-ignored.
Local plaintext settings remain accessible to the development user and Docker
administrators; they are not a production secrets-management solution.

## Compose operation and preservation

Run these commands from `infra/auth`. Always use both files and the explicit
project; a bare Compose command selects the historical default project.

~~~powershell
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local -f compose.yaml -f compose.https.yaml config --quiet
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local -f compose.yaml -f compose.https.yaml ps
~~~

The override requires Compose 2.24.4+ for `!override`. It changes only Keycloak,
replacing `start-dev` with `start --import-realm` and the HTTP publication with
HTTPS. Private management health checks retain HTTP through
`KC_HTTP_MANAGEMENT_SCHEME=http`. PostgreSQL service, network and volume remain
unchanged. Local `start` mode alone does not make this deployment production-ready.

For an explicitly approved Keycloak recreation, use the same configuration:

~~~powershell
docker compose --env-file local/ce1-auth-local.env -p ce1-auth-local -f compose.yaml -f compose.https.yaml up -d --no-deps --pull never --wait keycloak
~~~

Do not run base-only `up` during normal HTTPS operation. Before a future migration,
record the existing settings and container/volume identity, take an encrypted,
access-restricted PostgreSQL logical backup and verify an isolated restoration.
That backup/restore verification passed before the local HTTPS migration.
Neither a realm seed nor a running-volume filesystem copy replaces a database
backup.

Rollback needs a separately reviewed procedure: revert only targeted persisted
realm/client fields, restore the realm's prior SSL policy before depending on
HTTP administration, then recreate only Keycloak with the original base-only
configuration and the same project/database. Rollback was not exercised during
the successful migration. Never delete volumes, reset/import over the realm or
automatically restore over the primary database.

Retained `ce1-auth-dev` and `ce1-auth-import-test` projects and the isolated fixture
are intentionally kept. Do not clean them up or reuse their credentials/volumes
without approval.

## Realm seed and verified live settings

`keycloak/ce1-dev-realm.json` initializes fresh realms natively. It contains no
application users, credentials or signing-key exports. The client secret remains
an environment placeholder. Existing realms are skipped on startup: editing the
seed does not update live settings.

The client uses confidential Authorization Code Flow, mandatory PKCE S256,
restricted role scopes, an exact callback and empty Web Origins. Its override
requires Username Form followed by WebAuthn Passwordless; user verification and
discoverable credentials are required. Existing users/passkeys, signing keys and
client secret were preserved during HTTPS migration. Realm `ce1-dev` requires
SSL `all`; its WebAuthn RP ID remains `localhost`.

In addition to the role and username mappers, the seed includes the verified
live corrections:

| Mapper | Configuration |
| --- | --- |
| `ce1-api audience` / `oidc-audience-mapper` | Custom audience `ce1-api`; access/introspection claims true; ID/UserInfo claims false |
| `application subject` / `oidc-sub-mapper` | Access/introspection claims true; no explicit ID-token or UserInfo setting |

The Subject mapper fixes the initial callback rejection caused by an access
token without `sub`. Token validation was kept strict. Do not duplicate mappers
or replace the live client to synchronize a seed.

For private administrative readback on Windows PowerShell 5.1, assign REST JSON
arrays directly (`$mappers = Invoke-RestMethod ...`). Wrapping that command in
`@(...)` can create a nested array and produce misleading count/filter failures.

## Go harness and automated checks

Supply these variables privately to the harness process:

| Variable | Purpose |
| --- | --- |
| `CE1_AUTH_TEST_DATABASE_URL` | Existing isolated database at the exact host/port/database/user above |
| `AUTH_BFF_CLIENT_SECRET` | Existing persisted client secret |
| `AUTHDEV_CA_FILE` | Public CA PEM |
| `AUTHDEV_CERT_FILE` | Harness certificate PEM |
| `AUTHDEV_KEY_FILE` | Harness private-key PEM |

From `apps`, certificate paths begin `../infra/auth/local/tls/`. The harness checks
the resolved pgx target and fallbacks before connecting, then verifies
`current_database()` and `current_user`. It creates no schema or migrations.
The isolated fixture's local `sslmode=disable` is a development exception;
certificate-verified PostgreSQL TLS remains a deployment requirement.

Run non-database checks from `apps`:

~~~powershell
go test -mod=readonly ./... -count=1
go test -mod=readonly -tags authdev ./... -count=1
go vet -mod=readonly -tags authdev ./...
~~~

After private settings and the approved fixture/schema are ready, run manually:

~~~powershell
go run -mod=readonly -tags authdev ./cmd/authdev
~~~

Open the harness origin, log in with an enrolled passkey, check `/auth/me`, then
log out. Routes are mounted under `/auth/` without stripping the prefix.
The forms require exact Origin checks; only the root form page uses
`Referrer-Policy: same-origin`. Auth endpoints retain `no-referrer` and
`Cache-Control: no-store`. Session cookies are Secure, HttpOnly, SameSite=Lax
and host-only. Logout ends the local application session, not Keycloak SSO.
Stop only the harness process started for testing; retain database fixtures.

## Verification results and limits

Local Phase 2C testing completed on 2026-11-19:

- Trusted HTTPS discovery/JWKS, existing-passkey login and confidential PKCE code
  exchange succeeded through the unchanged strict OIDC validator.
- `/auth/me` returned an application UUID and the dispatcher role; a read-only
  check confirmed its active, unexpired PostgreSQL session and identity mapping.
- Manual logout was followed by `/auth/me` returning 401. The exact active session
  row was physically absent while its original expiry was still 59 seconds away,
  proving revocation by deletion rather than expiry.
- The successful callback consumed its browser-bound one-time transaction.
  Unrelated pending flows were preserved. Primary container/storage identities
  and realm security settings remained intact; the harness was stopped afterward.

See [persistence/README.md](persistence/README.md) for atomic replacement,
concurrency semantics and opt-in database tests. Automated TLS issuer tests cover
token failures and key rotation; browser cancellation, pre-enrollment denial,
password-authenticated SSO bypass resistance, cross-platform passkeys and
enrollment-credential cleanup remain unverified.

Phase 2D session-management work is next. Shared API wiring, coordinated versioned
database migrations, incident authorization, dispatcher device/session binding,
responder device authentication and Kubernetes/Tailscale deployment are pending.
Roles are snapshots; Keycloak role/account changes do not immediately revoke a
local session. No refresh tokens, automatic renewal or expired-row cleanup exist.
Keep this environment loopback-only and free of real incident data.

Historical enrollment context: [PHASE2A.md](PHASE2A.md).
Official references: [Keycloak imports](https://www.keycloak.org/server/importExport),
[HTTPS](https://www.keycloak.org/server/enabletls),
[management interface](https://www.keycloak.org/server/management-interface),
[Compose merging](https://docs.docker.com/reference/compose-file/merge/),
[mkcert](https://github.com/FiloSottile/mkcert).
