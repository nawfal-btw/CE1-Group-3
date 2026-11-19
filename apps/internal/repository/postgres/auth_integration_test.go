//go:build integration

package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"flag"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var authPostgres = flag.Bool("auth-postgres", false, "run against an explicitly approved, pre-provisioned authentication test database")

func TestAuthStorePostgres(t *testing.T) {
	if !*authPostgres {
		t.Skip("requires -auth-postgres and a separately approved PostgreSQL test environment")
	}
	dsn := os.Getenv("CE1_AUTH_TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("set CE1_AUTH_TEST_DATABASE_URL privately; no database is created by this test")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	open := func() *sql.DB {
		db, err := sql.Open("pgx", dsn)
		if err != nil {
			t.Fatal("invalid test database configuration (URL withheld)")
		}
		t.Cleanup(func() { _ = db.Close() })
		db.SetMaxOpenConns(8)
		var database string
		if err := db.QueryRowContext(ctx, "SELECT current_database()").Scan(&database); err != nil {
			t.Fatal("cannot connect to the approved test database (connection details withheld)")
		}
		if !strings.HasSuffix(database, "_auth_test") {
			t.Fatal("refusing fixture writes: database name must end in _auth_test; never use Keycloak or the application database")
		}
		return db
	}
	// Separate pools model independent Go processes, not a shared transaction.
	dbA, dbB := open(), open()
	a, b := NewAuthStore(dbA), NewAuthStore(dbB)
	cleanup := func(t *testing.T, query string, key any) {
		t.Helper()
		t.Cleanup(func() {
			cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cleanupCancel()
			if _, err := dbA.ExecContext(cleanupCtx, query, key); err != nil {
				t.Error("could not remove this test's authentication fixture")
			}
		})
	}
	identity := func(t *testing.T) uuid.UUID {
		t.Helper()
		id, err := a.GetOrCreateIdentity(ctx, "https://issuer.invalid/"+uuid.NewString(), "opaque-subject")
		if err != nil {
			t.Fatalf("create identity fixture: %v", err)
		}
		cleanup(t, "DELETE FROM auth.identities WHERE id = $1", id)
		return id
	}
	login := func(t *testing.T) AuthLoginTransaction {
		t.Helper()
		v := authTestLogin()
		v.StateHash = authTestHash(uuid.NewString())
		if err := a.CreateLoginTransaction(ctx, v); err != nil {
			t.Fatal(err)
		}
		cleanup(t, "DELETE FROM auth.login_transactions WHERE state_hash = $1", v.StateHash)
		return v
	}
	session := func(t *testing.T, id uuid.UUID, roles []string) AuthSession {
		t.Helper()
		v := authTestSession()
		v.IdentityID, v.Roles, v.TokenHash = id, roles, authTestHash(uuid.NewString())
		if err := a.CreateSession(ctx, v); err != nil {
			t.Fatal(err)
		}
		cleanup(t, "DELETE FROM auth.sessions WHERE token_hash = $1", v.TokenHash)
		return v
	}

	t.Run("concurrent identities and issuer isolation", func(t *testing.T) {
		issuer := "https://issuer.invalid/" + uuid.NewString()
		otherIssuer := "https://other.invalid/" + uuid.NewString()
		cleanup(t, "DELETE FROM auth.identities WHERE issuer = $1", issuer)
		cleanup(t, "DELETE FROM auth.identities WHERE issuer = $1", otherIssuer)
		const subject = "opaque-subject-'not-a-uuid"
		type result struct {
			id  uuid.UUID
			err error
		}
		results := make(chan result, 16)
		start := make(chan struct{})
		for i := 0; i < 16; i++ {
			store := a
			if i%2 != 0 {
				store = b
			}
			go func() {
				<-start
				id, err := store.GetOrCreateIdentity(ctx, issuer, subject)
				results <- result{id, err}
			}()
		}
		close(start)
		var first uuid.UUID
		for i := 0; i < 16; i++ {
			v := <-results
			if v.err != nil {
				t.Errorf("identity creation failed: %v", v.err)
				continue
			}
			if first == uuid.Nil {
				first = v.id
			}
			if v.id == uuid.Nil || v.id != first {
				t.Error("concurrent calls returned different identities")
			}
		}
		other, err := b.GetOrCreateIdentity(ctx, otherIssuer, subject)
		if err != nil || other == first {
			t.Fatalf("issuer isolation failed: %v", err)
		}
		otherSubject, err := b.GetOrCreateIdentity(ctx, issuer, subject+"-other")
		if err != nil || otherSubject == first {
			t.Fatalf("subject isolation failed: %v", err)
		}
		again, err := b.GetOrCreateIdentity(ctx, issuer, subject)
		if err != nil || again != first {
			t.Fatalf("identity was not stable: %v", err)
		}
	})

	t.Run("browser binding concurrent consumption and replay", func(t *testing.T) {
		v := login(t)
		if err := b.CreateLoginTransaction(ctx, v); err == nil {
			t.Fatal("duplicate state was accepted")
		}
		if _, err := b.ConsumeLoginTransaction(ctx, v.StateHash, authTestHash("wrong-browser")); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("wrong browser: %v", err)
		}
		results := make(chan error, 16)
		start := make(chan struct{})
		var wg sync.WaitGroup
		for i := 0; i < 16; i++ {
			store := a
			if i%2 != 0 {
				store = b
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-start
				got, err := store.ConsumeLoginTransaction(ctx, v.StateHash, v.BrowserHash)
				if err == nil && (!bytes.Equal(got.NonceHash, v.NonceHash) || got.PKCEVerifier != v.PKCEVerifier) {
					err = errors.New("consumed transaction lost its nonce/verifier")
				}
				results <- err
			}()
		}
		close(start)
		wg.Wait()
		close(results)
		successes := 0
		for err := range results {
			if err == nil {
				successes++
			} else if !errors.Is(err, ErrAuthNotFound) {
				t.Error(err)
			}
		}
		if successes != 1 {
			t.Fatalf("got %d successful consumers; want 1", successes)
		}
		if _, err := b.ConsumeLoginTransaction(ctx, v.StateHash, v.BrowserHash); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("replayed state: %v", err)
		}
	})

	t.Run("transaction expiration", func(t *testing.T) {
		v := login(t)
		if _, err := dbA.ExecContext(ctx, `UPDATE auth.login_transactions
			SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes',
			    expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute'
			WHERE state_hash = $1`, v.StateHash); err != nil {
			t.Fatal(err)
		}
		if _, err := b.ConsumeLoginTransaction(ctx, v.StateHash, v.BrowserHash); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("expired transaction: %v", err)
		}
		v.StateHash = authTestHash(uuid.NewString())
		v.ExpiresAt = time.Now().Add(-time.Hour)
		if err := b.CreateLoginTransaction(ctx, v); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("past expiry: %v", err)
		}
	})

	t.Run("shared sessions duplicate protection logout and expiration", func(t *testing.T) {
		id := identity(t)
		roles := []string{"dispatcher", "role-with-'quotes-å"}
		v := session(t, id, roles)
		got, err := b.GetSession(ctx, v.TokenHash)
		if err != nil || got.IdentityID != id || !reflect.DeepEqual(got.Roles, roles) || !bytes.Equal(got.TokenHash, v.TokenHash) ||
			got.ExpiresAt.Sub(v.ExpiresAt) > time.Microsecond || v.ExpiresAt.Sub(got.ExpiresAt) > time.Microsecond {
			t.Fatalf("cross-pool session lookup failed: %v", err)
		}
		duplicate := v
		duplicate.Roles = []string{"administrator"}
		if err := b.CreateSession(ctx, duplicate); err == nil {
			t.Fatal("duplicate session replaced an existing session")
		}
		got, err = a.GetSession(ctx, v.TokenHash)
		if err != nil || !reflect.DeepEqual(got.Roles, roles) {
			t.Fatal("duplicate creation modified roles")
		}
		if err := b.DeleteSession(ctx, v.TokenHash); err != nil {
			t.Fatal(err)
		}
		if _, err := a.GetSession(ctx, v.TokenHash); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("session survived logout: %v", err)
		}
		if err := a.DeleteSession(ctx, v.TokenHash); err != nil {
			t.Fatalf("repeated logout: %v", err)
		}
		if _, err := a.GetSession(ctx, authTestHash(uuid.NewString())); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("missing session: %v", err)
		}
		empty := session(t, id, nil)
		got, err = b.GetSession(ctx, empty.TokenHash)
		if err != nil || got.Roles == nil || len(got.Roles) != 0 {
			t.Fatalf("empty roles: %v", err)
		}
		if _, err := dbA.ExecContext(ctx, `UPDATE auth.sessions
			SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes',
			    expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute'
			WHERE token_hash = $1`, empty.TokenHash); err != nil {
			t.Fatal(err)
		}
		if _, err := a.GetSession(ctx, empty.TokenHash); !errors.Is(err, ErrAuthNotFound) {
			t.Fatalf("expired session: %v", err)
		}
		expired := authTestSession()
		expired.IdentityID, expired.ExpiresAt = id, time.Now().Add(-time.Hour)
		if err := a.CreateSession(ctx, expired); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("past session expiry: %v", err)
		}
		unknown := authTestSession()
		if err := b.CreateSession(ctx, unknown); err == nil {
			t.Fatal("session accepted unknown identity")
		}
	})

	t.Run("atomic replacement and account switching", func(t *testing.T) {
		previousID, nextID := identity(t), identity(t)
		previous := session(t, previousID, []string{"dispatcher"})
		next := authTestSession()
		next.TokenHash, next.IdentityID = authTestHash(uuid.NewString()), nextID
		cleanup(t, "DELETE FROM auth.sessions WHERE token_hash = $1", next.TokenHash)
		if err := a.ReplaceSession(ctx, next, previous.TokenHash); err != nil {
			t.Fatal(err)
		}
		if _, err := b.GetSession(ctx, previous.TokenHash); !errors.Is(err, ErrAuthNotFound) {
			t.Fatal("old session survived replacement")
		}
		got, err := b.GetSession(ctx, next.TokenHash)
		if err != nil || got.IdentityID != nextID || !reflect.DeepEqual(got.Roles, next.Roles) {
			t.Fatal("replacement did not preserve fresh account identity/roles")
		}
	})

	for _, mode := range []string{"no previous", "missing previous", "revoked previous", "expired previous"} {
		t.Run("replacement/"+mode, func(t *testing.T) {
			id := identity(t)
			var previousHash []byte
			if mode == "missing previous" {
				previousHash = authTestHash(uuid.NewString())
			}
			if mode == "revoked previous" {
				previousHash = session(t, id, nil).TokenHash
				if err := b.DeleteSession(ctx, previousHash); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "expired previous" {
				previousHash = session(t, id, nil).TokenHash
				if _, err := dbA.ExecContext(ctx, `UPDATE auth.sessions SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes',
					expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute' WHERE token_hash = $1`, previousHash); err != nil {
					t.Fatal(err)
				}
			}
			next := authTestSession()
			next.TokenHash, next.IdentityID = authTestHash(uuid.NewString()), id
			cleanup(t, "DELETE FROM auth.sessions WHERE token_hash = $1", next.TokenHash)
			if err := a.ReplaceSession(ctx, next, previousHash); err != nil {
				t.Fatal(err)
			}
			if _, err := b.GetSession(ctx, next.TokenHash); err != nil {
				t.Fatal("new session missing")
			}
			if len(previousHash) > 0 {
				var count int
				if err := dbB.QueryRowContext(ctx, "SELECT count(*) FROM auth.sessions WHERE token_hash = $1", previousHash).Scan(&count); err != nil || count != 0 {
					t.Fatal("previous row was not removed")
				}
			}
		})
	}

	for _, mode := range []string{"duplicate new hash", "unknown identity", "expired new session"} {
		t.Run("replacement rollback/"+mode, func(t *testing.T) {
			id := identity(t)
			previous := session(t, id, nil)
			next := authTestSession()
			next.TokenHash, next.IdentityID = authTestHash(uuid.NewString()), id
			if mode == "duplicate new hash" {
				next.TokenHash = session(t, id, nil).TokenHash
			}
			if mode == "unknown identity" {
				next.IdentityID = uuid.New()
			}
			if mode == "expired new session" {
				next.ExpiresAt = time.Now().Add(-time.Hour)
			}
			if err := a.ReplaceSession(ctx, next, previous.TokenHash); err == nil {
				t.Fatal("invalid insertion succeeded")
			}
			if _, err := b.GetSession(ctx, previous.TokenHash); err != nil {
				t.Fatal("insertion failure removed previous session")
			}
			if mode != "duplicate new hash" {
				if _, err := b.GetSession(ctx, next.TokenHash); !errors.Is(err, ErrAuthNotFound) {
					t.Fatal("failed insertion left new session")
				}
			} else {
				if _, err := b.GetSession(ctx, next.TokenHash); err != nil {
					t.Fatal("duplicate insertion changed the existing conflicting session")
				}
			}
		})
	}

	for _, sameNewHash := range []bool{false, true} {
		name := "concurrent independent replacements"
		if sameNewHash {
			name = "concurrent duplicate new hash"
		}
		t.Run(name, func(t *testing.T) {
			id := identity(t)
			previous := session(t, id, nil)
			const contenders = 16
			next := make([]AuthSession, contenders)
			sharedHash := authTestHash(uuid.NewString())
			for i := range next {
				next[i] = authTestSession()
				next[i].IdentityID, next[i].TokenHash = id, authTestHash(uuid.NewString())
				if sameNewHash {
					next[i].TokenHash = sharedHash
				}
				cleanup(t, "DELETE FROM auth.sessions WHERE token_hash = $1", next[i].TokenHash)
			}
			start := make(chan struct{})
			results := make(chan error, contenders)
			for i := range next {
				store := a
				if i%2 != 0 {
					store = b
				}
				go func() { <-start; results <- store.ReplaceSession(ctx, next[i], previous.TokenHash) }()
			}
			close(start)
			successes := 0
			for range contenders {
				if err := <-results; err == nil {
					successes++
				} else if !sameNewHash {
					t.Error(err)
				}
			}
			expected := contenders
			if sameNewHash {
				expected = 1
			}
			if successes != expected {
				t.Fatalf("got %d commits, want %d", successes, expected)
			}
			if _, err := b.GetSession(ctx, previous.TokenHash); !errors.Is(err, ErrAuthNotFound) {
				t.Fatal("old session survived concurrent replacements")
			}
			for _, value := range next {
				if _, err := b.GetSession(ctx, value.TokenHash); err != nil {
					t.Fatal("committed successor missing")
				}
			}
		})
	}

	t.Run("database constraints reject invalid direct writes", func(t *testing.T) {
		id := identity(t)
		if _, err := dbA.ExecContext(ctx, `INSERT INTO auth.sessions (token_hash, identity_id, roles, expires_at)
			VALUES ($1, $2, '[]'::jsonb, CURRENT_TIMESTAMP + INTERVAL '5 minutes')`, []byte{1}, id); err == nil {
			t.Error("database accepted a short session hash")
		}
		if _, err := dbA.ExecContext(ctx, `INSERT INTO auth.sessions (token_hash, identity_id, roles, expires_at)
			VALUES ($1, $2, 'null'::jsonb, CURRENT_TIMESTAMP + INTERVAL '5 minutes')`, authTestHash(uuid.NewString()), id); err == nil {
			t.Error("database accepted non-array roles")
		}
	})
}
