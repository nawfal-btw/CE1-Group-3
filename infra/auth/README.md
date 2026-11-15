# Local authentication infrastructure

This isolated Compose project runs Keycloak and its own PostgreSQL database. It
does not run or alter the application, legacy stack, or Kubernetes. Run commands
below in PowerShell from `infra/auth` unless stated otherwise.

## Versions and boundaries

- Keycloak: `quay.io/keycloak/keycloak:26.8.1`.
- PostgreSQL: `postgres:17.11-bookworm`. PostgreSQL 17 is supported by Keycloak.
- Compose project: `ce1-auth-dev`, with its own network and named volume
  `ce1-auth-dev_postgres_data`.
- Only Keycloak is published: `127.0.0.1:8081` forwards to container port 8080.
  PostgreSQL and Keycloak's management/health/metrics port 9000 are not published.
- Canonical hostname: `http://localhost:8081`. Use this consistently for browser
  access and OIDC issuer matching; `127.0.0.1` is the bind address, not the issuer.

Image tags pin release versions, not immutable image digests. Review security
updates and test upgrades deliberately. Never point this stack at an application
database or reuse an unrelated database volume.

## First-time configuration

Prerequisites: Docker Desktop running with Linux containers and Docker Compose.
No native PostgreSQL, Keycloak, Java, or Kubernetes installation is needed.

From the repository root:

```powershell
Set-Location infra/auth
docker context show
docker info --format '{{.OperatingSystem}} | {{.OSType}} | {{.ServerVersion}}'
docker compose version
```

Confirm the server is Docker Desktop with `linux` container type. If the engine
cannot be reached, start Docker Desktop and wait for readiness before continuing.
Do not change contexts or start an unrelated Docker engine automatically.

Create `.env` only if it does not already exist:

```powershell
if (Test-Path -LiteralPath .env) { throw '.env already exists; preserve it.' }
Copy-Item -LiteralPath .env.example -Destination .env
```

Set `AUTH_DB_PASSWORD` and `AUTH_ADMIN_PASSWORD` to different random passwords in
`.env`; leave `AUTH_ADMIN_USERNAME=local-admin` unless a different local name is
needed. Empty passwords intentionally fail Compose validation. Generate each
password separately with this PowerShell-compatible snippet, then paste it into
the local file (do not share the output):

```powershell
$bytes = New-Object byte[] 32
$rng = [System.Security.Cryptography.RandomNumberGenerator]::Create()
try { $rng.GetBytes($bytes) } finally { $rng.Dispose() }
[Convert]::ToBase64String($bytes)
```

`.env`, `.env.*` (except `.env.example`), and `local/` are ignored by Git. Use
`local/` for sensitive local exports if needed. The bootstrap admin is a local
Keycloak administrator in the `master` realm, not a dispatcher test identity.
Credentials are not included in the realm JSON. Keep filesystem access restricted
and never reuse these credentials outside this disposable environment.

## Validate and start

```powershell
docker compose config --quiet
if ($LASTEXITCODE -ne 0) { throw 'Compose validation failed.' }
docker compose up -d --wait --wait-timeout 240
if ($LASTEXITCODE -ne 0) { throw 'Startup failed; inspect status and logs.' }
docker compose ps
```

The first start pulls images, initializes PostgreSQL, prepares Keycloak, and
imports `ce1-dev`. PostgreSQL readiness gates Keycloak startup. Keycloak's
healthcheck requests `/health/ready` on its internal management port using Bash
built-ins because the official image does not include curl. Metrics are enabled
so readiness can include database health. No extra image or software is installed.

Do not share `docker compose config` without `--quiet`, `config --environment`,
full `docker inspect`, or credential-bearing diagnostic output: they can expose
resolved passwords. `.env` is a development convenience, not a secret manager;
Docker administrators can inspect container environment variables.

Open <http://localhost:8081/admin/> and log in using the local bootstrap admin.
Select realm **ce1-dev** and check **Realm roles**:

| Role | Meaning in this phase |
| --- | --- |
| `administrator` | Provisional application administrator, not a Keycloak administrator |
| `dispatcher` | Provisional dispatcher classification |
| `operator` | Provisional operator classification |
| `responder` | Provisional responder classification; no device or vehicle identity |

These roles grant no implemented API permissions and are not automatically
assigned to users. Application users, custom OIDC clients, passkeys, device
credentials, token middleware, sessions, and resource authorization are deferred.
Keycloak creates its normal built-in clients and default realm roles.

## Verify

```powershell
docker compose ps
docker compose exec -T postgres pg_isready -h 127.0.0.1 -U keycloak -d keycloak
powershell -NoProfile -ExecutionPolicy RemoteSigned -File .\verify.ps1
```

The script checks discovery's exact issuer, advertised endpoints and PKCE S256;
fetches JWKS; checks RSA signing keys and absence of private/symmetric key
parameters; and confirms an unknown realm returns HTTP 404. It obtains no tokens
and changes no server state. `RemoteSigned` above applies only to this child
PowerShell process so a locally created script can run; it does not change user
or machine policy. If organizational policy prevents scripts, use the read-only
requests below; do not change persistent policy to run this project.

Manual endpoint checks:

```powershell
$issuer = 'http://localhost:8081/realms/ce1-dev'
$discovery = Invoke-RestMethod "$issuer/.well-known/openid-configuration"
$discovery | Select-Object issuer, authorization_endpoint, token_endpoint, jwks_uri
$jwks = Invoke-RestMethod $discovery.jwks_uri
$jwks.keys | Select-Object kid, kty, use, alg
```

Discovery describes the identity provider's OIDC endpoints. JWKS publishes public
verification keys; successful retrieval does not prove application JWT validation
or authorization. Those are later phases.

## Shutdown, restart, and persistence

```powershell
docker compose stop
docker compose start --wait --wait-timeout 240
```

`stop` preserves containers and the database volume. PostgreSQL holds realm
configuration, users, and signing keys; no separate Keycloak data volume is
required. If configuration changes require container recreation, use the startup
`up` command again after reviewing the change.

For container/network teardown, run `docker compose down` only after deciding to
remove this project's containers/network (agents require explicit approval).
Ordinary `down` preserves the named database volume. Do not add `--volumes` or
use volume/system prune commands: they can permanently erase identities and keys.

Startup import **skips an existing realm**. Editing the JSON does not update the
persisted realm. Do not reset the database to apply changes. Plan a deliberate
Admin Console/API update or separately approved disposable reset. Bootstrap
credentials also apply only during initial administrator creation; editing `.env`
does not rotate an existing administrator password. Change that password in
Keycloak, and coordinate the local file separately.

## Troubleshooting

| Symptom | Action |
| --- | --- |
| Docker named pipe missing | Start Docker Desktop, select Linux containers, and repeat `docker info`. |
| Required environment variable error | Confirm `.env` exists in this directory with both nonempty passwords. |
| Port 8081 in use | Inspect `Get-NetTCPConnection -LocalPort 8081 -State Listen`; do not stop its owner automatically. A port change also requires updating the canonical hostname, verification script, and future clients. |
| Pull/DNS/proxy failure | Check Docker Desktop connectivity/proxy settings and retry; do not disable TLS verification. |
| Startup timeout or unhealthy container | Run `docker compose ps` and `docker compose logs --tail 100 postgres keycloak`; review locally before sharing logs. Initial image preparation can take time. |
| PostgreSQL authentication fails after `.env` edit | Initial volume credentials persist; changing the environment does not change the database password. Coordinate an explicit password rotation rather than deleting the volume. |
| Realm discovery returns 404 | Check Keycloak logs for import errors and ensure `ce1-dev` is selected/enabled. Existing realms are not overwritten by startup import. |
| Unexpected issuer or browser redirects | Use `http://localhost:8081`, not container DNS names or a different hostname. |
| Admin login fails after `.env` edit | Existing administrators are not recreated. Use the existing credential or plan recovery; do not erase the database. |

## Development-only limitations

`start-dev`, HTTP, realm `sslRequired: none`, and unencrypted database traffic
inside the Compose network are local-development shortcuts. Loopback publishing
reduces remote exposure but does not protect against other local users or a
Docker administrator. The database initialization account has elevated database
privileges; production needs deliberate database privilege and secret management.

Do not expose this stack on LAN/Tailscale, use real incident data, or treat it as
production authentication. Shared deployment requires HTTPS and trusted hostnames,
production Keycloak startup/configuration, appropriate database transport and
privileges, managed secrets, backups/recovery, and an access-control review.
Future passkey enrollment needs a stable origin/RP ID; choose that before enrolling
credentials. This phase does not enforce passwordless user or device authentication.

## Official references

- [Keycloak downloads](https://www.keycloak.org/downloads)
- [Supported configurations](https://www.keycloak.org/server/supported-configurations)
- [Container configuration](https://www.keycloak.org/server/containers)
- [Health checks](https://www.keycloak.org/observability/health)
- [Realm import behavior](https://www.keycloak.org/server/importExport)
- [PostgreSQL version policy](https://www.postgresql.org/support/versioning/)
