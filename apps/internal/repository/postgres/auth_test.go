package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func authTestHash(value string) []byte {
	hash := sha256.Sum256([]byte(value))
	return hash[:]
}

func authTestLogin() AuthLoginTransaction {
	return AuthLoginTransaction{
		StateHash: authTestHash("state"), BrowserHash: authTestHash("browser"),
		NonceHash: authTestHash("nonce"), PKCEVerifier: strings.Repeat("a", 43),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

func authTestSession() AuthSession {
	return AuthSession{
		TokenHash: authTestHash("session"), IdentityID: uuid.New(),
		Roles: []string{"dispatcher"}, ExpiresAt: time.Now().Add(5 * time.Minute),
	}
}

// This deliberately has no database handle: invalid input must be rejected
// before storage access. These tests do not exercise PostgreSQL behavior.
func TestAuthStoreValidation(t *testing.T) {
	store := NewAuthStore(nil)
	ctx := context.Background()
	for _, input := range [][2]string{{"", "subject"}, {"issuer", ""}} {
		if id, err := store.GetOrCreateIdentity(ctx, input[0], input[1]); !errors.Is(err, ErrInvalidAuthInput) || id != uuid.Nil {
			t.Fatalf("invalid identity: id=%v error=%v", id, err)
		}
	}
	for _, hash := range [][]byte{nil, make([]byte, 31), make([]byte, 33)} {
		if _, err := store.GetSession(ctx, hash); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("invalid session hash: %v", err)
		}
		if err := store.DeleteSession(ctx, hash); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("invalid logout hash: %v", err)
		}
		if _, err := store.ConsumeLoginTransaction(ctx, hash, authTestHash("browser")); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("invalid state hash: %v", err)
		}
		if _, err := store.ConsumeLoginTransaction(ctx, authTestHash("state"), hash); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("invalid browser hash: %v", err)
		}
	}
	for name, change := range map[string]func(*AuthLoginTransaction){
		"state":          func(v *AuthLoginTransaction) { v.StateHash = nil },
		"browser":        func(v *AuthLoginTransaction) { v.BrowserHash = nil },
		"nonce":          func(v *AuthLoginTransaction) { v.NonceHash = nil },
		"short verifier": func(v *AuthLoginTransaction) { v.PKCEVerifier = strings.Repeat("a", 42) },
		"long verifier":  func(v *AuthLoginTransaction) { v.PKCEVerifier = strings.Repeat("a", 129) },
		"expiry":         func(v *AuthLoginTransaction) { v.ExpiresAt = time.Time{} },
	} {
		t.Run("login/"+name, func(t *testing.T) {
			input := authTestLogin()
			change(&input)
			if err := store.CreateLoginTransaction(ctx, input); !errors.Is(err, ErrInvalidAuthInput) {
				t.Fatalf("expected invalid input: %v", err)
			}
		})
	}
	for name, change := range map[string]func(*AuthSession){
		"hash":     func(v *AuthSession) { v.TokenHash = nil },
		"identity": func(v *AuthSession) { v.IdentityID = uuid.Nil },
		"expiry":   func(v *AuthSession) { v.ExpiresAt = time.Time{} },
	} {
		t.Run("session/"+name, func(t *testing.T) {
			input := authTestSession()
			change(&input)
			if err := store.CreateSession(ctx, input); !errors.Is(err, ErrInvalidAuthInput) {
				t.Fatalf("expected invalid input: %v", err)
			}
			if err := store.ReplaceSession(ctx, input, authTestHash("previous")); !errors.Is(err, ErrInvalidAuthInput) {
				t.Fatalf("expected invalid replacement input: %v", err)
			}
		})
	}
	input := authTestSession()
	for _, previous := range [][]byte{make([]byte, 31), make([]byte, 33), input.TokenHash} {
		if err := store.ReplaceSession(ctx, input, previous); !errors.Is(err, ErrInvalidAuthInput) {
			t.Fatalf("expected invalid previous hash: %v", err)
		}
	}
}

var errAuthTestDatabase = errors.New("test database unavailable")

// A standard-library connector that always fails, without network access or
// SQL emulation. It verifies error propagation, not SQL or concurrency semantics.
type authFailingConnector struct{}

func (authFailingConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, errAuthTestDatabase
}
func (authFailingConnector) Driver() driver.Driver { return authFailingDriver{} }

type authFailingDriver struct{}

func (authFailingDriver) Open(string) (driver.Conn, error) { return nil, errAuthTestDatabase }

func TestAuthStoreDatabaseFailures(t *testing.T) {
	for _, mode := range []string{"unavailable", "closed", "canceled"} {
		t.Run(mode, func(t *testing.T) {
			db := sql.OpenDB(authFailingConnector{})
			t.Cleanup(func() { _ = db.Close() })
			ctx := context.Background()
			if mode == "closed" {
				if err := db.Close(); err != nil {
					t.Fatal(err)
				}
			}
			if mode == "canceled" {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}
			store := NewAuthStore(db)
			operations := map[string]func() error{
				"identity": func() error {
					id, err := store.GetOrCreateIdentity(ctx, "issuer", "opaque-subject")
					if id != uuid.Nil {
						t.Error("failed identity lookup returned an identity")
					}
					return err
				},
				"create login": func() error { return store.CreateLoginTransaction(ctx, authTestLogin()) },
				"consume login": func() error {
					v, err := store.ConsumeLoginTransaction(ctx, authTestHash("state"), authTestHash("browser"))
					if !reflect.DeepEqual(v, AuthLoginTransaction{}) {
						t.Error("failed consumption returned transaction data")
					}
					return err
				},
				"create session":  func() error { return store.CreateSession(ctx, authTestSession()) },
				"replace session": func() error { return store.ReplaceSession(ctx, authTestSession(), authTestHash("previous")) },
				"get session": func() error {
					v, err := store.GetSession(ctx, authTestHash("session"))
					if !reflect.DeepEqual(v, AuthSession{}) {
						t.Error("failed session lookup returned session data")
					}
					return err
				},
				"delete session": func() error { return store.DeleteSession(ctx, authTestHash("session")) },
			}
			for name, operation := range operations {
				t.Run(name, func(t *testing.T) {
					err := operation()
					if err == nil || errors.Is(err, ErrAuthNotFound) || errors.Is(err, ErrInvalidAuthInput) {
						t.Fatalf("expected storage failure, got %v", err)
					}
					if mode == "unavailable" && !errors.Is(err, errAuthTestDatabase) {
						t.Errorf("lost database cause: %v", err)
					}
					if mode == "canceled" && !errors.Is(err, context.Canceled) {
						t.Errorf("lost context cancellation: %v", err)
					}
				})
			}
		})
	}
}
