# Authentication persistence

PostgreSQL authentication persistence and its opt-in integration tests have
passed against the isolated `ce1_auth_test` fixture. Local Phase 2C browser testing
also verified session creation and logout deletion before expiry. This directory
contains development SQL, **not a numbered application migration**. Never apply
it to Keycloak's database.

## Files and integration contract

- `schema.sql` defines `auth.identities`, `auth.login_transactions` and
  `auth.sessions` in one explicit transaction. It fails if the schema already
  exists; there is no automatic DDL or silent drift handling.
- `apps/internal/repository/postgres/auth.go` provides `AuthStore` using
  parameterized SQL through `database/sql`. At application integration, reuse
  the existing GORM pool through `gormDB.DB()`; callers own its lifecycle.
- `auth_test.go` covers input validation and database failure handling.
  `auth_integration_test.go` covers real PostgreSQL behavior and concurrency.
  `apps/internal/auth/login_test.go` adds opt-in HTTP/store integration coverage.

Independent development uses container `ce1-auth-persistence-test`, volume
`ce1-auth-persistence-test-data`, loopback port `55432` and database/user
`ce1_auth_test`. This is separate from the primary Keycloak PostgreSQL instance.
Private connection settings remain under ignored `infra/auth/local/`.

Production integration must use the team's existing application PostgreSQL
server. Coordinate a version in `apps/db/migrations` and paired up/down migrations
with the database teammate; no version has been allocated here. Do not apply this
SQL and then run an equivalent numbered migration without coordinating history.
No AutoMigrate or per-replica schema creation is intended.

Keep DDL ownership separate from runtime access. Runtime needs schema USAGE,
identity SELECT/INSERT/UPDATE, login-transaction SELECT/INSERT/DELETE and session
SELECT/INSERT/DELETE. Tests also need fixture UPDATE and identity DELETE.
Certificate-verified database TLS remains a deployment requirement; the isolated
local fixture currently uses an explicitly approved unencrypted connection.

## Identity and one-time transactions

`GetOrCreateIdentity` maps the opaque `(issuer, subject)` pair to an application
UUID. A unique constraint and `INSERT ... ON CONFLICT ... DO UPDATE ... RETURNING`
preserve that UUID under concurrent first use across connections. The conflict
update takes a row lock. Changing issuer does not automatically link identities.

The login handler generates high-entropy state, nonce, browser binding and PKCE
values. Persistence receives 32-byte hashes for state, nonce and browser binding,
and temporarily stores the recoverable PKCE verifier. Restrict table/backup access.
Login transactions expire after five minutes under the handler's policy.

`ConsumeLoginTransaction` atomically deletes and returns the matching unexpired
browser-bound row using `DELETE ... RETURNING`. Only one concurrent consumer
succeeds. Wrong bindings leave the row intact; missing, expired and consumed rows
all return `ErrAuthNotFound`. Consumption occurs before code exchange, without
holding a database transaction across a provider call. A later failure requires
a fresh login; the consumed transaction is not restored.

## Sessions and atomic replacement

Sessions contain only a hash of the opaque bearer identifier, identity UUID,
role snapshot, creation time and expiry. No ID/access/refresh tokens are stored.
Nil roles are normalized to `[]`. Creation and lookup enforce expiry using the
database clock. All Go replicas query the shared database; no memory cache,
sliding expiry or last-seen writes are used.

`CreateSession` delegates to `ReplaceSession` without a previous hash.
`ReplaceSession` uses one READ COMMITTED transaction:

1. Validate the new session and optional previous hash.
2. Insert the new, unexpired session.
3. Delete the previous session if a hash was supplied.
4. Commit; only then may the handler issue the new browser cookie.

Statement failures roll back the insertion and preserve the previous session.
Missing, already-revoked or expired previous sessions are acceptable; an expired
row is deleted if present. Account switching is supported. Invalid hash lengths,
identical previous/new hashes, duplicate new hashes and invalid session inputs
are rejected.

Concurrent independently authenticated replacements of the same old session
may all commit with distinct new hashes. This is not a one-active-session-per-user
policy. Contenders using the same new hash cannot both insert it.

A COMMIT error can mean either rollback or a lost acknowledgement after a real
commit. Never automatically retry or assume the old session survived. The handler
returns generic 503 and issues no new cookie on error; a committed-but-unissued
session can remain until expiry/cleanup. No compensating deletion is attempted.

`GetSession` returns only unexpired rows. `DeleteSession` is idempotent; after
successful deletion, subsequent lookups from any replica reject the session.
Already-authenticated requests may finish. Logout clears the browser cookie only
after acknowledged deletion and preserves unrelated pending login transactions.
A deletion acknowledgement failure returns 503 without clearing the cookie.

Store errors can contain driver details. HTTP boundaries sanitize them:
missing/invalid sessions produce 401, genuine lookup failures produce 503.
Do not expose or log raw SQL errors or authentication parameters.

## Tests

Non-database checks from `apps`:

~~~powershell
go test -mod=readonly ./... -count=1
go test -mod=readonly -tags authdev ./... -count=1
go vet -mod=readonly -tags authdev ./...
~~~

Database-writing tests are explicitly opt-in. Only use a separately approved,
pre-provisioned disposable database with this schema. Privately supply
`CE1_AUTH_TEST_DATABASE_URL`; no connection string is committed.

~~~powershell
go test -mod=readonly -tags integration ./internal/auth ./internal/repository/postgres -count=1 -args -auth-postgres
~~~

The persistence suite requires both its build tag and the opt-in flag. Login
database tests require the flag; without it they skip. Both suites reject a
database name not ending in `_auth_test`. That suffix guard is not authorization
to use an unrelated database. Tests create no databases or schemas, perform no
migrations, and remove only their unique fixture rows; interruption can leave
fixtures requiring separately scoped cleanup.

Prior real PostgreSQL tests passed for identity concurrency, replay prevention,
expiry, cross-pool lookup/revocation, first login, replacement, account switching,
missing/expired old sessions, insertion/deletion rollback, concurrent replacement,
transaction visibility, cookie failure behavior and controlled COMMIT errors.
The COMMIT-error cases simulate rollback and lost acknowledgement through test
driver wrappers; they are not proof of real network-partition recovery.
Non-database OIDC tests use a local HTTPS issuer, not the live Keycloak service.

Local browser testing confirmed dispatcher session persistence and physical logout
deletion while the original session expiry remained in the future. This does not
establish production failover, process-restart behavior, real replica deployment
or database TLS. Database-writing tests need not run during ordinary unit CI.

## Remaining work

Shared application wiring and coordinated migrations are pending. Phase 2D
session-management work follows; automatic expired-row cleanup, account-wide
revocation, refresh/renewal, Keycloak SSO logout and device binding are not
implemented. Expired rows are unusable but retained until cleanup. Role snapshots
do not immediately reflect Keycloak account/role changes. Authentication remains
separate from incident authorization and device/resource policy.
