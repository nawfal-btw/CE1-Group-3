-- Development authentication schema only. Not a numbered application migration.
-- Apply explicitly to an approved, isolated application test database, never to
-- Keycloak's database. There is deliberately no IF NOT EXISTS or automatic DDL.
BEGIN;

CREATE SCHEMA auth;

CREATE TABLE auth.identities (
    id uuid PRIMARY KEY CHECK (id <> '00000000-0000-0000-0000-000000000000'),
    issuer text NOT NULL CHECK (issuer <> ''),
    subject text NOT NULL CHECK (subject <> ''),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (issuer, subject)
);

CREATE TABLE auth.login_transactions (
    state_hash bytea PRIMARY KEY CHECK (octet_length(state_hash) = 32),
    browser_hash bytea NOT NULL CHECK (octet_length(browser_hash) = 32),
    nonce_hash bytea NOT NULL CHECK (octet_length(nonce_hash) = 32),
    pkce_verifier text NOT NULL CHECK (length(pkce_verifier) BETWEEN 43 AND 128),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL CHECK (expires_at > created_at)
);

CREATE INDEX login_transactions_expiry ON auth.login_transactions (expires_at);

CREATE TABLE auth.sessions (
    token_hash bytea PRIMARY KEY CHECK (octet_length(token_hash) = 32),
    identity_id uuid NOT NULL REFERENCES auth.identities (id),
    roles jsonb NOT NULL CHECK (jsonb_typeof(roles) = 'array'),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL CHECK (expires_at > created_at)
);

CREATE INDEX sessions_identity ON auth.sessions (identity_id);
CREATE INDEX sessions_expiry ON auth.sessions (expires_at);

COMMIT;
