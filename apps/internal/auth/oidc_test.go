package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
)

// All issuer requests stay on local test servers. Keys exist only in memory.
type testIssuer struct {
	server            *httptest.Server
	key               *rsa.PrivateKey
	mu                sync.Mutex
	keys              []jose.JSONWebKey
	metadata          map[string]any
	discoveryRedirect string
	keysRedirect      string
	keysStatus        int
	keyRequests       int
}

func newTestIssuer(t *testing.T) *testIssuer {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	f := &testIssuer{key: key}
	f.keys = []jose.JSONWebKey{{Key: &key.PublicKey, KeyID: "first", Algorithm: "RS256", Use: "sig"}}
	f.server = httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/.well-known/openid-configuration":
			if f.discoveryRedirect != "" {
				http.Redirect(w, r, f.discoveryRedirect, http.StatusFound)
				return
			}
			metadata := map[string]any{"issuer": f.server.URL, "jwks_uri": f.server.URL + "/keys", "id_token_signing_alg_values_supported": []string{"RS256"}}
			for key, value := range f.metadata {
				metadata[key] = value
			}
			_ = json.NewEncoder(w).Encode(metadata)
		case "/keys":
			f.keyRequests++
			if f.keysRedirect != "" {
				http.Redirect(w, r, f.keysRedirect, http.StatusFound)
				return
			}
			if f.keysStatus != 0 {
				w.WriteHeader(f.keysStatus)
				return
			}
			_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: f.keys})
		default:
			http.NotFound(w, r)
		}
	}))
	f.server.Config.ErrorLog = log.New(io.Discard, "", 0) // Expected TLS failures are tested below.
	f.server.StartTLS()
	t.Cleanup(f.server.Close)
	return f
}

func (f *testIssuer) config() OIDCConfig {
	roots := x509.NewCertPool()
	roots.AddCert(f.server.Certificate())
	return OIDCConfig{IssuerURL: f.server.URL, ClientID: "ce1-bff-dev", APIAudience: "ce1-api",
		HumanRoles: []string{"administrator", "dispatcher", "operator"}, RootCAs: roots}
}

func (f *testIssuer) validator(t *testing.T) *OIDC {
	t.Helper()
	validator, err := NewOIDC(context.Background(), f.config())
	if err != nil {
		t.Fatal(err)
	}
	return validator
}

func (f *testIssuer) claims() (map[string]any, map[string]any) {
	now := time.Now().Unix()
	id := map[string]any{"iss": f.server.URL, "sub": "opaque-human-subject", "aud": "ce1-bff-dev",
		"iat": now - 60, "exp": now + 120, "nonce": "transaction-nonce"}
	access := map[string]any{"iss": f.server.URL, "sub": "opaque-human-subject", "aud": "ce1-api",
		"iat": now - 60, "exp": now + 60, "azp": "ce1-bff-dev", "typ": "Bearer",
		"realm_access": map[string]any{"roles": []string{"dispatcher", "unrelated-role", "dispatcher", "responder"}}}
	return id, access
}

func signTestToken(t *testing.T, key *rsa.PrivateKey, kid string, algorithm jose.SignatureAlgorithm, claims map[string]any) string {
	t.Helper()
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: algorithm, Key: key},
		(&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", kid))
	if err != nil {
		t.Fatal(err)
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		t.Fatal(err)
	}
	signed, err := signer.Sign(payload)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := signed.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func (f *testIssuer) validate(t *testing.T, validator *OIDC, id, access map[string]any) (ValidatedIdentity, error) {
	t.Helper()
	nonce := sha256.Sum256([]byte("transaction-nonce"))
	return validator.Validate(context.Background(), signTestToken(t, f.key, "first", jose.RS256, id),
		signTestToken(t, f.key, "first", jose.RS256, access), nonce[:])
}

func TestOIDCValidation(t *testing.T) {
	f := newTestIssuer(t)
	validator := f.validator(t)
	t.Run("valid identity roles and access expiry", func(t *testing.T) {
		id, access := f.claims()
		got, err := f.validate(t, validator, id, access)
		if err != nil || got.Issuer != f.server.URL || got.Subject != "opaque-human-subject" ||
			!reflect.DeepEqual(got.Roles, []string{"dispatcher"}) || !got.ExpiresAt.Equal(time.Unix(access["exp"].(int64), 0)) {
			t.Fatalf("incorrect validated identity: %+v, %v", got, err)
		}
	})
	t.Run("ID expiry bounds identity", func(t *testing.T) {
		id, access := f.claims()
		id["exp"] = time.Now().Unix() + 45
		got, err := f.validate(t, validator, id, access)
		if err != nil || !got.ExpiresAt.Equal(time.Unix(id["exp"].(int64), 0)) {
			t.Fatalf("incorrect expiry: %v", err)
		}
	})
	t.Run("multiple audiences with authorized party and small clock skew", func(t *testing.T) {
		id, access := f.claims()
		id["aud"], id["azp"] = []string{"ce1-bff-dev", "trusted-secondary"}, "ce1-bff-dev"
		id["nbf"], access["iat"] = time.Now().Unix()+20, time.Now().Unix()+20
		if _, err := f.validate(t, validator, id, access); err != nil {
			t.Fatal(err)
		}
	})

	for name, change := range map[string]func(map[string]any, map[string]any){
		"wrong ID issuer":                      func(id, access map[string]any) { id["iss"] = "https://other.invalid" },
		"wrong access issuer":                  func(id, access map[string]any) { access["iss"] = "https://other.invalid" },
		"wrong ID audience":                    func(id, access map[string]any) { id["aud"] = "other-client" },
		"wrong access audience":                func(id, access map[string]any) { access["aud"] = "ce1-bff-dev" },
		"missing access audience":              func(id, access map[string]any) { delete(access, "aud") },
		"expired ID":                           func(id, access map[string]any) { id["exp"] = time.Now().Unix() - 1 },
		"expired access":                       func(id, access map[string]any) { access["exp"] = time.Now().Unix() - 1 },
		"missing ID expiry":                    func(id, access map[string]any) { delete(id, "exp") },
		"missing access expiry":                func(id, access map[string]any) { delete(access, "exp") },
		"wrong nonce":                          func(id, access map[string]any) { id["nonce"] = "wrong" },
		"missing nonce":                        func(id, access map[string]any) { delete(id, "nonce") },
		"wrong ID client":                      func(id, access map[string]any) { id["azp"] = "other-client" },
		"multiple ID audiences without client": func(id, access map[string]any) { id["aud"] = []string{"ce1-bff-dev", "other"} },
		"wrong access client":                  func(id, access map[string]any) { access["azp"] = "other-client" },
		"missing access client":                func(id, access map[string]any) { delete(access, "azp") },
		"ID token type in access slot":         func(id, access map[string]any) { access["typ"] = "ID" },
		"missing access type":                  func(id, access map[string]any) { delete(access, "typ") },
		"only responder role": func(id, access map[string]any) {
			access["realm_access"] = map[string]any{"roles": []string{"responder"}}
		},
		"missing roles":            func(id, access map[string]any) { delete(access, "realm_access") },
		"empty roles":              func(id, access map[string]any) { access["realm_access"] = map[string]any{"roles": []string{}} },
		"malformed roles":          func(id, access map[string]any) { access["realm_access"] = map[string]any{"roles": "dispatcher"} },
		"mismatching subjects":     func(id, access map[string]any) { access["sub"] = "other-human" },
		"empty subjects":           func(id, access map[string]any) { id["sub"], access["sub"] = "", "" },
		"missing ID issued at":     func(id, access map[string]any) { delete(id, "iat") },
		"missing access issued at": func(id, access map[string]any) { delete(access, "iat") },
		"future ID issued at":      func(id, access map[string]any) { id["iat"] = time.Now().Unix() + 60 },
		"future access issued at":  func(id, access map[string]any) { access["iat"] = time.Now().Unix() + 60 },
		"future ID not before":     func(id, access map[string]any) { id["nbf"] = time.Now().Unix() + 60 },
		"future access not before": func(id, access map[string]any) { access["nbf"] = time.Now().Unix() + 60 },
		"ID issued at after expiry": func(id, access map[string]any) {
			now := time.Now().Unix()
			id["exp"], id["iat"] = now+20, now+25
		},
		"access issued at after expiry": func(id, access map[string]any) {
			now := time.Now().Unix()
			access["exp"], access["iat"] = now+20, now+25
		},
		"ID not before after expiry": func(id, access map[string]any) {
			now := time.Now().Unix()
			id["exp"], id["nbf"] = now+20, now+25
		},
		"access not before after expiry": func(id, access map[string]any) {
			now := time.Now().Unix()
			access["exp"], access["nbf"] = now+20, now+25
		},
		"malformed not before": func(id, access map[string]any) { access["nbf"] = "tomorrow" },
		"wrong access hash":    func(id, access map[string]any) { id["at_hash"] = "wrong" },
	} {
		t.Run(name, func(t *testing.T) {
			id, access := f.claims()
			change(id, access)
			got, err := f.validate(t, validator, id, access)
			if !errors.Is(err, ErrInvalidTokens) || !reflect.DeepEqual(got, ValidatedIdentity{}) {
				t.Fatalf("invalid token pair returned identity: %+v, %v", got, err)
			}
		})
	}
	t.Run("valid access hash", func(t *testing.T) {
		id, access := f.claims()
		rawAccess := signTestToken(t, f.key, "first", jose.RS256, access)
		hash := sha256.Sum256([]byte(rawAccess))
		id["at_hash"] = base64.RawURLEncoding.EncodeToString(hash[:16])
		nonce := sha256.Sum256([]byte("transaction-nonce"))
		if _, err := validator.Validate(context.Background(), signTestToken(t, f.key, "first", jose.RS256, id), rawAccess, nonce[:]); err != nil {
			t.Fatal(err)
		}
	})
}

func TestOIDCSignaturesAndBinding(t *testing.T) {
	f := newTestIssuer(t)
	validator := f.validator(t)
	id, access := f.claims()
	rawID, rawAccess := signTestToken(t, f.key, "first", jose.RS256, id), signTestToken(t, f.key, "first", jose.RS256, access)
	wrongKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	nonce := sha256.Sum256([]byte("transaction-nonce"))
	for _, tc := range []struct {
		name, id, access string
		nonce            []byte
	}{
		{"wrong ID signature", signTestToken(t, wrongKey, "first", jose.RS256, id), rawAccess, nonce[:]},
		{"wrong access signature", rawID, signTestToken(t, wrongKey, "first", jose.RS256, access), nonce[:]},
		{"wrong ID algorithm", signTestToken(t, f.key, "first", jose.RS512, id), rawAccess, nonce[:]},
		{"wrong access algorithm", rawID, signTestToken(t, f.key, "first", jose.RS512, access), nonce[:]},
		{"unsigned ID", "eyJhbGciOiJub25lIn0.e30.", rawAccess, nonce[:]},
		{"malformed access", rawID, "not-a-jwt", nonce[:]},
		{"swapped tokens", rawAccess, rawID, nonce[:]},
		{"missing nonce expectation", rawID, rawAccess, nil},
		{"malformed nonce expectation", rawID, rawAccess, nonce[:31]},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validator.Validate(context.Background(), tc.id, tc.access, tc.nonce)
			if !errors.Is(err, ErrInvalidTokens) || !reflect.DeepEqual(got, ValidatedIdentity{}) {
				t.Fatalf("accepted invalid tokens: %v", err)
			}
		})
	}
}

func TestOIDCConfigurationAndTLS(t *testing.T) {
	f := newTestIssuer(t)
	for name, change := range map[string]func(*OIDCConfig){
		"HTTP issuer":           func(c *OIDCConfig) { c.IssuerURL = "http://localhost" },
		"missing host":          func(c *OIDCConfig) { c.IssuerURL = "https:///realm" },
		"credentials in issuer": func(c *OIDCConfig) { c.IssuerURL = "https://user:secret@localhost" },
		"query in issuer":       func(c *OIDCConfig) { c.IssuerURL += "?query=value" },
		"empty query in issuer": func(c *OIDCConfig) { c.IssuerURL += "?" },
		"fragment in issuer":    func(c *OIDCConfig) { c.IssuerURL += "#fragment" },
		"missing client":        func(c *OIDCConfig) { c.ClientID = "" },
		"missing API audience":  func(c *OIDCConfig) { c.APIAudience = "" },
		"same audiences":        func(c *OIDCConfig) { c.APIAudience = c.ClientID },
		"missing human roles":   func(c *OIDCConfig) { c.HumanRoles = nil },
		"empty human role":      func(c *OIDCConfig) { c.HumanRoles = []string{""} },
		"untrusted certificate": func(c *OIDCConfig) { c.RootCAs = x509.NewCertPool() },
	} {
		t.Run(name, func(t *testing.T) {
			config := f.config()
			change(&config)
			if got, err := NewOIDC(context.Background(), config); err == nil || got != nil {
				t.Fatal("accepted invalid configuration or TLS trust")
			}
		})
	}
	for name, metadata := range map[string]map[string]any{
		"issuer mismatch":           {"issuer": "https://other.invalid"},
		"HTTP signing keys":         {"jwks_uri": "http://localhost/keys"},
		"cross origin signing keys": {"jwks_uri": "https://other.invalid/keys"},
		"missing signing keys":      {"jwks_uri": ""},
	} {
		t.Run(name, func(t *testing.T) {
			f.mu.Lock()
			f.metadata = metadata
			f.mu.Unlock()
			if got, err := NewOIDC(context.Background(), f.config()); err == nil || got != nil {
				t.Fatal("accepted invalid discovery")
			}
		})
	}
	f.mu.Lock()
	f.metadata = nil
	f.mu.Unlock()
	// A redirected response must never reach this target, including a downgrade.
	var targetMu sync.Mutex
	targetRequests := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { targetMu.Lock(); targetRequests++; targetMu.Unlock() }))
	t.Cleanup(target.Close)
	t.Run("discovery redirect rejected", func(t *testing.T) {
		f.mu.Lock()
		f.discoveryRedirect = target.URL
		f.mu.Unlock()
		if _, err := NewOIDC(context.Background(), f.config()); err == nil {
			t.Fatal("accepted discovery redirect")
		}
		f.mu.Lock()
		f.discoveryRedirect = ""
		f.mu.Unlock()
	})
	t.Run("JWKS redirect rejected", func(t *testing.T) {
		validator := f.validator(t)
		f.mu.Lock()
		f.keysRedirect = target.URL
		f.mu.Unlock()
		id, access := f.claims()
		if _, err := f.validate(t, validator, id, access); !errors.Is(err, ErrInvalidTokens) {
			t.Fatalf("accepted JWKS redirect: %v", err)
		}
		f.mu.Lock()
		f.keysRedirect = ""
		f.mu.Unlock()
	})
	targetMu.Lock()
	defer targetMu.Unlock()
	if targetRequests != 0 {
		t.Fatal("HTTP client followed a forbidden redirect")
	}
}

func TestOIDCKeyRotationAndProviderFailure(t *testing.T) {
	f := newTestIssuer(t)
	validator := f.validator(t)
	id, access := f.claims()
	if _, err := f.validate(t, validator, id, access); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	initialRequests := f.keyRequests
	f.mu.Unlock()
	if initialRequests != 1 {
		t.Fatalf("want one shared JWKS fetch, got %d", initialRequests)
	}
	if _, err := f.validate(t, validator, id, access); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	cachedRequests := f.keyRequests
	f.mu.Unlock()
	if cachedRequests != initialRequests {
		t.Fatal("signing keys were not cached")
	}
	rotated, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	nonce := sha256.Sum256([]byte("transaction-nonce"))
	rawID, rawAccess := signTestToken(t, rotated, "second", jose.RS256, id), signTestToken(t, rotated, "second", jose.RS256, access)
	f.mu.Lock()
	f.keysStatus = http.StatusServiceUnavailable
	f.mu.Unlock()
	if _, err := validator.Validate(context.Background(), rawID, rawAccess, nonce[:]); !errors.Is(err, ErrInvalidTokens) {
		t.Fatalf("provider failure accepted unknown key: %v", err)
	}
	f.mu.Lock()
	f.keysStatus = 0
	f.keys = []jose.JSONWebKey{{Key: &rotated.PublicKey, KeyID: "second", Algorithm: "RS256", Use: "sig"}}
	beforeRotation := f.keyRequests
	f.mu.Unlock()
	if _, err := validator.Validate(context.Background(), rawID, rawAccess, nonce[:]); err != nil {
		t.Fatal(err)
	}
	f.mu.Lock()
	afterRotation := f.keyRequests
	f.mu.Unlock()
	if afterRotation != beforeRotation+1 {
		t.Fatal("rotated keys were not shared between verifiers")
	}
}
