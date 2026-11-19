package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var ErrAuthNotFound = errors.New("authentication record not found or expired")
var ErrInvalidAuthInput = errors.New("invalid authentication persistence input")

// AuthStore shares the application's PostgreSQL pool. For a GORM connection,
// obtain that pool with db.DB(). Queries are parameterized and are not logged.
// The caller owns the pool and must not pass a transaction-scoped connection:
// consuming a login transaction must commit before any OIDC network request.
type AuthStore struct {
	db *sql.DB
}

func NewAuthStore(db *sql.DB) *AuthStore {
	return &AuthStore{db: db}
}

// AuthLoginTransaction contains hashes of random state, browser-binding and
// nonce values. Only the short-lived PKCE verifier is stored in recoverable form.
// Generating and verifying these values belongs to the future OIDC layer.
type AuthLoginTransaction struct {
	StateHash    []byte
	BrowserHash  []byte
	NonceHash    []byte
	PKCEVerifier string
	ExpiresAt    time.Time
}

// AuthSession contains a snapshot of roles from validated access-token claims.
// This store does not validate tokens or grant permissions. The caller must
// bound ExpiresAt by the validated tokens' expiry and the agreed session policy.
type AuthSession struct {
	TokenHash  []byte
	IdentityID uuid.UUID
	Roles      []string
	ExpiresAt  time.Time
}

// GetOrCreateIdentity preserves opaque, issuer-scoped subjects and returns a
// stable application UUID. The unique constraint and conflict update make this
// safe across independent connections/replicas, including concurrent first use.
func (s *AuthStore) GetOrCreateIdentity(ctx context.Context, issuer, subject string) (uuid.UUID, error) {
	if issuer == "" || subject == "" {
		return uuid.Nil, ErrInvalidAuthInput
	}
	id, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate authentication identity: %w", err)
	}
	err = s.db.QueryRowContext(ctx, `
		INSERT INTO auth.identities (id, issuer, subject)
		VALUES ($1, $2, $3)
		ON CONFLICT (issuer, subject) DO UPDATE SET subject = EXCLUDED.subject
		RETURNING id`, id, issuer, subject).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("get or create authentication identity: %w", err)
	}
	return id, nil
}

func (s *AuthStore) CreateLoginTransaction(ctx context.Context, transaction AuthLoginTransaction) error {
	if len(transaction.StateHash) != 32 || len(transaction.BrowserHash) != 32 ||
		len(transaction.NonceHash) != 32 || len(transaction.PKCEVerifier) < 43 ||
		len(transaction.PKCEVerifier) > 128 || transaction.ExpiresAt.IsZero() {
		return ErrInvalidAuthInput
	}
	result, err := s.db.ExecContext(ctx, `
		INSERT INTO auth.login_transactions
			(state_hash, browser_hash, nonce_hash, pkce_verifier, expires_at)
		SELECT $1, $2, $3, $4, $5 WHERE $5 > CURRENT_TIMESTAMP`,
		transaction.StateHash, transaction.BrowserHash, transaction.NonceHash,
		transaction.PKCEVerifier, transaction.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create authentication login transaction: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check authentication login transaction creation: %w", err)
	}
	if count != 1 {
		return ErrInvalidAuthInput
	}
	return nil
}

// ConsumeLoginTransaction atomically deletes and returns a matching unexpired
// transaction. Wrong browser bindings do not consume it. No row means missing,
// expired, mismatched or already consumed; no distinction is exposed to callers.
// A successful call is committed (autocommit), even if later code exchange fails.
func (s *AuthStore) ConsumeLoginTransaction(ctx context.Context, stateHash, browserHash []byte) (AuthLoginTransaction, error) {
	if len(stateHash) != 32 || len(browserHash) != 32 {
		return AuthLoginTransaction{}, ErrInvalidAuthInput
	}
	var transaction AuthLoginTransaction
	err := s.db.QueryRowContext(ctx, `
		DELETE FROM auth.login_transactions
		WHERE state_hash = $1 AND browser_hash = $2 AND expires_at > CURRENT_TIMESTAMP
		RETURNING state_hash, browser_hash, nonce_hash, pkce_verifier, expires_at`,
		stateHash, browserHash).Scan(&transaction.StateHash, &transaction.BrowserHash,
		&transaction.NonceHash, &transaction.PKCEVerifier, &transaction.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return AuthLoginTransaction{}, ErrAuthNotFound
	}
	if err != nil {
		return AuthLoginTransaction{}, fmt.Errorf("consume authentication login transaction: %w", err)
	}
	return transaction, nil
}

func (s *AuthStore) CreateSession(ctx context.Context, session AuthSession) error {
	return s.ReplaceSession(ctx, session, nil)
}

// ReplaceSession commits creation and previous-session revocation together.
// A missing/expired previous session is acceptable after fresh authentication;
// separate authenticated logins may each replace the same old session.
func (s *AuthStore) ReplaceSession(ctx context.Context, session AuthSession, previousTokenHash []byte) error {
	if len(session.TokenHash) != 32 || session.IdentityID == uuid.Nil || session.ExpiresAt.IsZero() {
		return ErrInvalidAuthInput
	}
	if len(previousTokenHash) != 0 && (len(previousTokenHash) != 32 || bytes.Equal(previousTokenHash, session.TokenHash)) {
		return ErrInvalidAuthInput
	}
	// Store an empty array rather than JSON null for identities with no roles.
	roles := session.Roles
	if roles == nil {
		roles = []string{}
	}
	rolesJSON, err := json.Marshal(roles)
	if err != nil {
		return fmt.Errorf("encode authentication session roles: %w", err)
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin authentication session replacement: %w", err)
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx, `
		INSERT INTO auth.sessions (token_hash, identity_id, roles, expires_at)
		SELECT $1, $2, $3::jsonb, $4 WHERE $4 > CURRENT_TIMESTAMP`,
		session.TokenHash, session.IdentityID, string(rolesJSON), session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("create authentication session: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check authentication session creation: %w", err)
	}
	if count != 1 {
		return ErrInvalidAuthInput
	}
	if len(previousTokenHash) != 0 {
		if _, err := tx.ExecContext(ctx, `DELETE FROM auth.sessions WHERE token_hash = $1`, previousTokenHash); err != nil {
			return fmt.Errorf("revoke previous authentication session: %w", err)
		}
	}
	// A commit error may mean its acknowledgement was lost after PostgreSQL
	// committed. Never retry or assume that a deferred rollback can undo it.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit authentication session replacement: %w", err)
	}
	return nil
}

func (s *AuthStore) GetSession(ctx context.Context, tokenHash []byte) (AuthSession, error) {
	if len(tokenHash) != 32 {
		return AuthSession{}, ErrInvalidAuthInput
	}
	var session AuthSession
	var rolesJSON []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT token_hash, identity_id, roles, expires_at FROM auth.sessions
		WHERE token_hash = $1 AND expires_at > CURRENT_TIMESTAMP`, tokenHash).
		Scan(&session.TokenHash, &session.IdentityID, &rolesJSON, &session.ExpiresAt)
	if errors.Is(err, sql.ErrNoRows) {
		return AuthSession{}, ErrAuthNotFound
	}
	if err != nil {
		return AuthSession{}, fmt.Errorf("get authentication session: %w", err)
	}
	if err := json.Unmarshal(rolesJSON, &session.Roles); err != nil {
		return AuthSession{}, fmt.Errorf("decode authentication session roles: %w", err)
	}
	return session, nil
}

// DeleteSession is idempotent. Database failures are still returned; callers
// must not report successful server-side logout when deletion fails.
func (s *AuthStore) DeleteSession(ctx context.Context, tokenHash []byte) error {
	if len(tokenHash) != 32 {
		return ErrInvalidAuthInput
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM auth.sessions WHERE token_hash = $1`, tokenHash); err != nil {
		return fmt.Errorf("delete authentication session: %w", err)
	}
	return nil
}
