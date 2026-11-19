package auth

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"flag"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	"golang.org/x/oauth2"
)

var loginPostgres = flag.Bool("auth-postgres", false, "allow login tests to write unique fixtures in the approved authentication test database")

const testOrigin = "https://bff.example.test"

func loginRequest(method, path string, cookies ...*http.Cookie) *http.Request {
	r := httptest.NewRequest(method, testOrigin+path, nil)
	r.Header.Set("Origin", testOrigin)
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	return r
}

func responseCookie(response *httptest.ResponseRecorder, name string) *http.Cookie {
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func checkCookie(t *testing.T, cookie *http.Cookie, deleted bool) {
	t.Helper()
	if cookie == nil {
		t.Fatal("missing expected cookie")
	}
	if !cookie.Secure || !cookie.HttpOnly || cookie.SameSite != http.SameSiteLaxMode || cookie.Path != "/" || cookie.Domain != "" {
		t.Fatal("incorrect authentication cookie attributes")
	}
	if deleted {
		if cookie.Value != "" || cookie.MaxAge != -1 {
			t.Fatal("binding cookie was not cleared on success")
		}
	} else if !validSecret(cookie.Value) || cookie.MaxAge <= 0 || !cookie.Expires.After(time.Now()) {
		t.Fatal("incorrect opaque cookie value or expiry")
	}
}

// sql.OpenDB does not connect; closing the pool guarantees guard tests cannot
// write to a database or accidentally make an outbound connection.
func closedLoginDB(t *testing.T) *sql.DB {
	t.Helper()
	config, err := pgx.ParseConfig("postgres://unused@127.0.0.1:1/unused?sslmode=disable")
	if err != nil {
		t.Fatal("invalid closed-pool fixture")
	}
	db := sql.OpenDB(stdlib.GetConnector(*config))
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	return db
}

func TestLoginGuards(t *testing.T) {
	f := newTestIssuer(t)
	f.mu.Lock()
	f.metadata = map[string]any{"authorization_endpoint": f.server.URL + "/authorize", "token_endpoint": f.server.URL + "/token"}
	f.mu.Unlock()
	validator := f.validator(t)
	store := postgres.NewAuthStore(closedLoginDB(t))
	h, err := NewLoginHandler(validator, store, LoginConfig{PublicOrigin: testOrigin, ClientSecret: "test-secret"})
	if err != nil {
		t.Fatal(err)
	}
	for name, change := range map[string]func(*LoginConfig){
		"HTTP origin":    func(c *LoginConfig) { c.PublicOrigin = "http://bff.example.test" },
		"origin path":    func(c *LoginConfig) { c.PublicOrigin += "/path" },
		"missing secret": func(c *LoginConfig) { c.ClientSecret = "" },
	} {
		t.Run(name, func(t *testing.T) {
			config := LoginConfig{PublicOrigin: testOrigin, ClientSecret: "test-secret"}
			change(&config)
			if _, err := NewLoginHandler(validator, store, config); err == nil {
				t.Fatal("accepted invalid login configuration")
			}
		})
	}
	for _, endpoint := range []string{"http://localhost/token", "https://other.invalid/token", ""} {
		copy := *validator
		copy.endpoint.TokenURL = endpoint
		if _, err := NewLoginHandler(&copy, store, LoginConfig{PublicOrigin: testOrigin, ClientSecret: "test-secret"}); err == nil {
			t.Fatal("accepted untrusted token endpoint")
		}
	}
	for name, tc := range map[string]struct {
		change func(*http.Request)
		status int
	}{
		"no TLS":                     {func(r *http.Request) { r.TLS = nil; r.Header.Set("X-Forwarded-Proto", "https") }, 403},
		"wrong host":                 {func(r *http.Request) { r.Host = "other.example.test" }, 403},
		"missing Origin":             {func(r *http.Request) { r.Header.Del("Origin") }, 403},
		"null Origin":                {func(r *http.Request) { r.Header.Set("Origin", "null") }, 403},
		"cross origin":               {func(r *http.Request) { r.Header.Set("Origin", "https://attacker.example.test") }, 403},
		"duplicate Origin":           {func(r *http.Request) { r.Header.Add("Origin", testOrigin) }, 403},
		"wrong method":               {func(r *http.Request) { r.Method = "GET" }, 405},
		"unavailable login database": {func(r *http.Request) {}, 503},
	} {
		t.Run(name, func(t *testing.T) {
			r := loginRequest("POST", "/auth/login")
			tc.change(r)
			w := httptest.NewRecorder()
			h.Login(w, r)
			if w.Code != tc.status || len(w.Result().Cookies()) != 0 || w.Header().Get("Location") != "" {
				t.Fatal("incorrect rejected login response")
			}
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("missing privacy headers")
			}
		})
	}
	state, binding := oauth2.GenerateVerifier(), oauth2.GenerateVerifier()
	for name, query := range map[string]string{
		"missing state":                 "code=code",
		"duplicate state":               "state=" + state + "&state=" + state + "&code=code",
		"duplicate code":                "state=" + state + "&code=a&code=b",
		"duplicate issuer":              "state=" + state + "&code=code&iss=a&iss=b",
		"duplicate error":               "state=" + state + "&error=a&error=b",
		"code and error":                "state=" + state + "&code=code&error=denied",
		"empty code and error":          "state=" + state + "&code=&error=denied",
		"empty error":                   "state=" + state + "&error=",
		"missing result":                "state=" + state,
		"bad encoding":                  "state=" + state + "&code=%zz",
		"wrong issuer":                  "state=" + state + "&code=code&iss=https://other.invalid",
		"error detail with code":        "state=" + state + "&code=code&error_description=detail",
		"untrusted completion redirect": "state=" + state + "&code=code&returnTo=https://other.invalid",
	} {
		t.Run(name, func(t *testing.T) {
			r := loginRequest("GET", "/auth/callback?"+query, &http.Cookie{Name: loginCookie, Value: binding})
			w := httptest.NewRecorder()
			h.Callback(w, r)
			if w.Code != 400 || len(w.Result().Cookies()) != 0 {
				t.Fatal("malformed callback accepted or changed a pending cookie")
			}
		})
	}
	for _, cookies := range [][]*http.Cookie{
		nil,
		{{Name: loginCookie, Value: "invalid"}},
		{{Name: loginCookie, Value: binding}, {Name: loginCookie, Value: binding}},
		{{Name: loginCookie, Value: binding}, {Name: sessionCookie, Value: "invalid"}},
	} {
		w := httptest.NewRecorder()
		h.Callback(w, loginRequest("GET", "/auth/callback?state="+state+"&code=code", cookies...))
		if w.Code != 400 || len(w.Result().Cookies()) != 0 {
			t.Fatal("invalid cookies accepted or cleared")
		}
	}
	w := httptest.NewRecorder()
	h.Callback(w, loginRequest("GET", "/auth/callback?state="+state+"&code=code", &http.Cookie{Name: loginCookie, Value: binding}))
	if w.Code != 503 || len(w.Result().Cookies()) != 0 {
		t.Fatal("callback database failure was not closed safely")
	}
}

func TestRequireSameOrigin(t *testing.T) {
	for _, publicOrigin := range []string{testOrigin, testOrigin + ":8443"} {
		t.Run(publicOrigin, func(t *testing.T) {
			origin, err := url.Parse(publicOrigin)
			if err != nil {
				t.Fatal(err)
			}
			// Nil store and OIDC dependencies ensure neither is used by the guard.
			h := &LoginHandler{origin: origin}
			for name, change := range map[string]func(*http.Request){
				"valid":                        func(r *http.Request) {},
				"valid without authentication": func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=invalid") },
				"missing Origin":               func(r *http.Request) { r.Header.Del("Origin") },
				"empty Origin":                 func(r *http.Request) { r.Header.Set("Origin", "") },
				"null Origin":                  func(r *http.Request) { r.Header.Set("Origin", "null") },
				"wrong Origin":                 func(r *http.Request) { r.Header.Set("Origin", "https://other.example.test") },
				"duplicate Origin":             func(r *http.Request) { r.Header.Add("Origin", publicOrigin) },
				"combined comma":               func(r *http.Request) { r.Header.Set("Origin", publicOrigin+", "+publicOrigin) },
				"combined space":               func(r *http.Request) { r.Header.Set("Origin", publicOrigin+" "+publicOrigin) },
				"malformed Origin":             func(r *http.Request) { r.Header.Set("Origin", "://%zz") },
				"Origin path":                  func(r *http.Request) { r.Header.Set("Origin", publicOrigin+"/") },
				"wrong Origin scheme":          func(r *http.Request) { r.Header.Set("Origin", "http://"+origin.Host) },
				"wrong Origin port":            func(r *http.Request) { r.Header.Set("Origin", "https://"+origin.Hostname()+":9443") },
				"non TLS":                      func(r *http.Request) { r.TLS = nil; r.Header.Set("X-Forwarded-Proto", "https") },
				"wrong host":                   func(r *http.Request) { r.Host = "other.example.test"; r.Header.Set("X-Forwarded-Host", origin.Host) },
				"wrong host port":              func(r *http.Request) { r.Host = origin.Hostname() + ":9443" },
				"Referer is insufficient":      func(r *http.Request) { r.Header.Del("Origin"); r.Header.Set("Referer", publicOrigin+"/") },
				"Forwarded is insufficient": func(r *http.Request) {
					r.TLS = nil
					r.Host = "internal.example.test"
					r.Header.Set("Forwarded", "proto=https;host="+origin.Host)
				},
			} {
				t.Run(name, func(t *testing.T) {
					for _, method := range []string{"POST", "PATCH"} {
						r := httptest.NewRequest(method, publicOrigin+"/v1/incident", nil)
						r.Header.Set("Origin", publicOrigin)
						change(r)
						headers := r.Header.Clone()
						calls := 0
						next := http.HandlerFunc(func(w http.ResponseWriter, received *http.Request) {
							calls++
							if received != r {
								t.Fatal("guard replaced the request")
							}
							w.WriteHeader(http.StatusNoContent)
						})
						w := httptest.NewRecorder()
						h.RequireSameOrigin(next).ServeHTTP(w, r)
						if strings.HasPrefix(name, "valid") {
							if calls != 1 || w.Code != 204 || w.Body.Len() != 0 {
								t.Fatal("valid request did not delegate exactly once")
							}
						} else if calls != 0 || w.Code != 403 || w.Body.String() != "Forbidden\n" {
							t.Fatal("invalid request reached downstream or lacked generic 403")
						}
						if len(w.Header().Values("Set-Cookie")) != 0 || !reflect.DeepEqual(r.Header, headers) {
							t.Fatal("guard changed cookies or request headers")
						}
					}
				})
			}
		})
	}
}

func TestAuthenticateIncidentGuards(t *testing.T) {
	origin, err := url.Parse(testOrigin)
	if err != nil {
		t.Fatal(err)
	}
	h := &LoginHandler{origin: origin, store: postgres.NewAuthStore(closedLoginDB(t))}
	secret := oauth2.GenerateVerifier()
	for name, tc := range map[string]struct {
		change func(*http.Request)
		want   error
	}{
		"missing":            {func(r *http.Request) { r.Header.Del("Cookie") }, ErrUnauthenticated},
		"malformed":          {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=invalid") }, ErrUnauthenticated},
		"empty":              {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=") }, ErrUnauthenticated},
		"duplicate":          {func(r *http.Request) { r.Header.Add("Cookie", sessionCookie+"="+secret) }, ErrUnauthenticated},
		"hidden duplicate":   {func(r *http.Request) { r.Header.Add("Cookie", sessionCookie+"=\"invalid") }, ErrUnauthenticated},
		"wrong host":         {func(r *http.Request) { r.Host = "other.example.test" }, ErrUnauthenticated},
		"non TLS":            {func(r *http.Request) { r.TLS = nil; r.Header.Set("X-Forwarded-Proto", "https") }, ErrUnauthenticated},
		"only pending login": {func(r *http.Request) { r.Header.Set("Cookie", loginCookie+"="+secret) }, ErrUnauthenticated},
		"body URL and bearer ignored": {func(r *http.Request) {
			r.Header.Del("Cookie")
			r.Header.Set("Authorization", "Bearer "+secret)
			r.URL.RawQuery = "session=" + secret + "&identityId=untrusted&roles=administrator"
			r.Body = io.NopCloser(strings.NewReader("session=" + secret))
		}, ErrUnauthenticated},
		"unavailable store": {func(r *http.Request) {}, ErrAuthUnavailable},
	} {
		t.Run(name, func(t *testing.T) {
			r := loginRequest("GET", "/v1/incident", &http.Cookie{Name: sessionCookie, Value: secret})
			tc.change(r)
			headers := r.Header.Clone()
			identity, err := h.AuthenticateIncident(r)
			if !errors.Is(err, tc.want) || err != tc.want || !reflect.DeepEqual(identity, application.Identity{}) {
				t.Fatal("incorrect sentinel or nonzero failed identity")
			}
			if errors.Unwrap(err) != nil || !reflect.DeepEqual(r.Header, headers) {
				t.Fatal("raw error chain or request cookie mutation")
			}
			if name == "body URL and bearer ignored" {
				body, err := io.ReadAll(r.Body)
				if err != nil || string(body) != "session="+secret {
					t.Fatal("adapter consumed request body")
				}
			}
		})
	}
}

func TestSessionGuards(t *testing.T) {
	origin, err := url.Parse(testOrigin)
	if err != nil {
		t.Fatal(err)
	}
	h := &LoginHandler{origin: origin, store: postgres.NewAuthStore(closedLoginDB(t))}
	secret := oauth2.GenerateVerifier()
	for _, endpoint := range []struct {
		path, method       string
		missing, malformed int
	}{
		{"/auth/me", "GET", 401, 401},
		{"/auth/logout", "POST", 204, 400},
	} {
		t.Run(endpoint.path, func(t *testing.T) {
			for name, tc := range map[string]struct {
				change func(*http.Request)
				status int
			}{
				"missing":             {func(r *http.Request) { r.Header.Del("Cookie") }, endpoint.missing},
				"malformed":           {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=invalid") }, endpoint.malformed},
				"empty":               {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=") }, endpoint.malformed},
				"no equals":           {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie) }, endpoint.malformed},
				"duplicate":           {func(r *http.Request) { r.Header.Add("Cookie", sessionCookie+"="+secret) }, endpoint.malformed},
				"malformed duplicate": {func(r *http.Request) { r.Header.Add("Cookie", sessionCookie+"=\"unterminated") }, endpoint.malformed},
				"malformed first":     {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"=\"bad; "+sessionCookie+"="+secret) }, endpoint.malformed},
				"noncanonical":        {func(r *http.Request) { r.Header.Set("Cookie", sessionCookie+"="+strings.Repeat("!", 43)) }, endpoint.malformed},
				"database failure":    {func(r *http.Request) {}, 503},
				"no TLS":              {func(r *http.Request) { r.TLS = nil; r.Header.Set("X-Forwarded-Proto", "https") }, 403},
				"wrong host":          {func(r *http.Request) { r.Host = "other.example.test" }, 403},
				"wrong method":        {func(r *http.Request) { r.Method = "PUT" }, 405},
				"URL and body credentials ignored": {func(r *http.Request) {
					r.Header.Del("Cookie")
					r.URL.RawQuery = "session=" + secret + "&identityId=" + uuid.NewString() + "&roles=administrator"
					r.Body = io.NopCloser(strings.NewReader("session=" + secret))
				}, endpoint.missing},
			} {
				t.Run(name, func(t *testing.T) {
					r := loginRequest(endpoint.method, endpoint.path, &http.Cookie{Name: sessionCookie, Value: secret})
					tc.change(r)
					w := httptest.NewRecorder()
					h.Routes().ServeHTTP(w, r)
					if w.Code != tc.status || len(w.Header().Values("Set-Cookie")) != 0 {
						t.Fatalf("status %d, expected %d, or unexpected cookie", w.Code, tc.status)
					}
					if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer" {
						t.Fatal("missing privacy headers")
					}
					if strings.Contains(w.Body.String(), secret) || strings.Contains(w.Body.String(), "sql:") {
						t.Fatal("sensitive response")
					}
				})
			}
		})
	}
	for _, origins := range [][]string{nil, {"null"}, {"https://other.example.test"}, {testOrigin, testOrigin}, {testOrigin + ", " + testOrigin}} {
		// Even an absent session must pass logout CSRF checks.
		r := loginRequest("POST", "/auth/logout")
		r.Header.Del("Origin")
		for _, origin := range origins {
			r.Header.Add("Origin", origin)
		}
		w := httptest.NewRecorder()
		h.Routes().ServeHTTP(w, r)
		if w.Code != 403 || len(w.Header().Values("Set-Cookie")) != 0 {
			t.Fatal("logout accepted invalid Origin")
		}
	}
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, loginRequest("GET", "/auth/unknown"))
	if w.Code != 404 || w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer" {
		t.Fatal("404 lacks privacy headers")
	}
}

type loginProvider struct {
	server    *httptest.Server
	issuer    string
	secret    string
	signer    jose.Signer
	key       jose.JSONWebKey
	mu        sync.Mutex
	codes     map[string]url.Values
	exchanges int
	mode      string
}

func newLoginProvider(t *testing.T) *loginProvider {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: key}, (&jose.SignerOptions{}).WithType("JWT").WithHeader("kid", "login"))
	if err != nil {
		t.Fatal(err)
	}
	f := &loginProvider{secret: "test:client secret+", signer: signer, codes: make(map[string]url.Values),
		key: jose.JSONWebKey{Key: &key.PublicKey, KeyID: "login", Use: "sig", Algorithm: "RS256"}}
	f.server = httptest.NewTLSServer(http.HandlerFunc(f.serveHTTP))
	f.mu.Lock()
	f.issuer = f.server.URL + "/realm/" + uuid.NewString()
	f.mu.Unlock()
	t.Cleanup(f.server.Close)
	return f
}

func (f *loginProvider) handler(t *testing.T, db *sql.DB) *LoginHandler {
	t.Helper()
	roots := x509.NewCertPool()
	roots.AddCert(f.server.Certificate())
	validator, err := NewOIDC(context.Background(), OIDCConfig{IssuerURL: f.issuer, ClientID: "ce1-bff-dev", APIAudience: "ce1-api", HumanRoles: []string{"dispatcher"}, RootCAs: roots})
	if err != nil {
		t.Fatal(err)
	}
	h, err := NewLoginHandler(validator, postgres.NewAuthStore(db), LoginConfig{PublicOrigin: testOrigin, ClientSecret: f.secret})
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func (f *loginProvider) serveHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.URL.Path == strings.TrimPrefix(f.issuer, f.server.URL)+"/.well-known/openid-configuration":
		_ = json.NewEncoder(w).Encode(map[string]any{"issuer": f.issuer, "jwks_uri": f.server.URL + "/keys", "authorization_endpoint": f.server.URL + "/authorize", "token_endpoint": f.server.URL + "/token"})
	case r.URL.Path == "/keys":
		_ = json.NewEncoder(w).Encode(jose.JSONWebKeySet{Keys: []jose.JSONWebKey{f.key}})
	case r.URL.Path == "/authorize":
		q := r.URL.Query()
		if q.Get("client_id") != "ce1-bff-dev" || q.Get("response_type") != "code" || q.Get("scope") != "openid" ||
			q.Get("redirect_uri") != testOrigin+"/auth/callback" || q.Get("response_mode") != "query" ||
			q.Get("code_challenge_method") != "S256" || !validSecret(q.Get("code_challenge")) || !validSecret(q.Get("state")) || !validSecret(q.Get("nonce")) {
			http.Error(w, "invalid authorization request", 400)
			return
		}
		code := oauth2.GenerateVerifier()
		f.codes[code] = q
		callback := url.Values{"state": {q.Get("state")}, "code": {code}, "iss": {f.issuer}, "session_state": {"test"}}
		http.Redirect(w, r, testOrigin+"/auth/callback?"+callback.Encode(), 303)
	case r.URL.Path == "/token":
		f.exchanges++
		if r.Method != "POST" || r.URL.RawQuery != "" || r.ParseForm() != nil {
			http.Error(w, "invalid exchange", 400)
			return
		}
		user, password, ok := r.BasicAuth()
		user, userErr := url.QueryUnescape(user)
		password, passwordErr := url.QueryUnescape(password)
		attempt, found := f.codes[r.Form.Get("code")]
		if !ok || userErr != nil || passwordErr != nil || user != "ce1-bff-dev" || password != f.secret ||
			r.Form.Has("client_secret") || r.Form.Has("client_id") || r.Form.Get("grant_type") != "authorization_code" ||
			!found || r.Form.Get("redirect_uri") != attempt.Get("redirect_uri") ||
			oauth2.S256ChallengeFromVerifier(r.Form.Get("code_verifier")) != attempt.Get("code_challenge") {
			http.Error(w, "invalid exchange", 400)
			return
		}
		delete(f.codes, r.Form.Get("code"))
		if f.mode == "provider error" {
			http.Error(w, "sensitive provider diagnostics", 400)
			return
		}
		now := time.Now().Unix()
		idExpiry, accessExpiry := now+120, now+60
		if f.mode == "long tokens" {
			idExpiry, accessExpiry = now+7200, now+3600
		}
		id := map[string]any{"iss": f.issuer, "sub": "login-human", "aud": "ce1-bff-dev", "iat": now - 1, "exp": idExpiry, "nonce": attempt.Get("nonce")}
		access := map[string]any{"iss": f.issuer, "sub": "login-human", "aud": "ce1-api", "azp": "ce1-bff-dev", "typ": "Bearer", "iat": now - 1, "exp": accessExpiry, "realm_access": map[string]any{"roles": []string{"dispatcher"}}}
		if f.mode == "invalid token" {
			id["nonce"] = "wrong-nonce"
		}
		rawID, err := f.sign(id)
		if err != nil {
			http.Error(w, "fixture signing failure", 500)
			return
		}
		rawAccess, err := f.sign(access)
		if err != nil {
			http.Error(w, "fixture signing failure", 500)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"access_token": rawAccess, "id_token": rawID, "token_type": "Bearer", "expires_in": 99999})
	default:
		http.NotFound(w, r)
	}
}

func (f *loginProvider) sign(claims map[string]any) (string, error) {
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signed, err := f.signer.Sign(payload)
	if err != nil {
		return "", err
	}
	return signed.CompactSerialize()
}

func (f *loginProvider) begin(t *testing.T, h *LoginHandler, db *sql.DB) (*http.Cookie, string, url.Values) {
	t.Helper()
	w := httptest.NewRecorder()
	h.Routes().ServeHTTP(w, loginRequest("POST", "/auth/login"))
	if w.Code != 303 {
		t.Fatal("login initiation failed")
	}
	binding := responseCookie(w, loginCookie)
	checkCookie(t, binding, false)
	location, err := url.Parse(w.Header().Get("Location"))
	if err != nil {
		t.Fatal("invalid authorization URL")
	}
	q := location.Query()
	t.Cleanup(func() {
		if _, err := db.Exec("DELETE FROM auth.login_transactions WHERE state_hash = $1", hashSecret(q.Get("state"))); err != nil {
			t.Error("login fixture cleanup failed")
		}
	})
	client := f.server.Client()
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	response, err := client.Get(location.String())
	if err != nil {
		t.Fatal("test authorization request failed")
	}
	defer response.Body.Close()
	if response.StatusCode != 303 {
		t.Fatal("authorization request failed state/nonce/PKCE checks")
	}
	callback, err := url.Parse(response.Header.Get("Location"))
	if err != nil {
		t.Fatal("invalid test callback URL")
	}
	return binding, callback.RequestURI(), q
}

// This connector injects a selected database error while forwarding other SQL
// to real PostgreSQL. It does not emulate persistence or alter schema/grants.
type loginFaultConnector struct {
	base   driver.Connector
	fail   func(string, []driver.NamedValue) bool
	commit func(driver.Tx) error
}

func (c loginFaultConnector) Driver() driver.Driver { return c.base.Driver() }
func (c loginFaultConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &loginFaultConn{Conn: conn, fail: c.fail, commit: c.commit}, nil
}

type loginFaultConn struct {
	driver.Conn
	fail   func(string, []driver.NamedValue) bool
	commit func(driver.Tx) error
}

func (c *loginFaultConn) BeginTx(ctx context.Context, options driver.TxOptions) (driver.Tx, error) {
	if options.Isolation != driver.IsolationLevel(sql.LevelReadCommitted) || options.ReadOnly {
		return nil, errors.New("expected a read-committed write transaction")
	}
	tx, err := c.Conn.(driver.ConnBeginTx).BeginTx(ctx, options)
	if err != nil || c.commit == nil {
		return tx, err
	}
	return &loginFaultTx{Tx: tx, commit: c.commit}, nil
}

// Commit hooks model both confirmed rollback and an acknowledgement lost after
// a real commit. The latter must never be described as a successful rollback.
type loginFaultTx struct {
	driver.Tx
	commit func(driver.Tx) error
}

func (tx *loginFaultTx) Commit() error { return tx.commit(tx.Tx) }

func (c *loginFaultConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.fail(query, args) {
		return nil, errors.New("injected database failure")
	}
	return c.Conn.(driver.ExecerContext).ExecContext(ctx, query, args)
}
func (c *loginFaultConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.fail(query, args) {
		return nil, errors.New("injected database failure")
	}
	return c.Conn.(driver.QueryerContext).QueryContext(ctx, query, args)
}

func TestLoginPostgres(t *testing.T) {
	if !*loginPostgres {
		t.Skip("requires -auth-postgres and the approved isolated database")
	}
	dsn := os.Getenv("CE1_AUTH_TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("set CE1_AUTH_TEST_DATABASE_URL privately; never use the Keycloak or application database")
	}
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("invalid test database configuration (URL withheld)")
	}
	config.ConnectTimeout = 3 * time.Second
	open := func(connector driver.Connector) *sql.DB {
		db := sql.OpenDB(connector)
		db.SetMaxOpenConns(8)
		t.Cleanup(func() { _ = db.Close() })
		return db
	}
	dbA, dbB := open(stdlib.GetConnector(*config)), open(stdlib.GetConnector(*config))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	var database string
	if err := dbA.QueryRowContext(ctx, "SELECT current_database()").Scan(&database); err != nil {
		t.Fatal("approved test database unavailable (connection details withheld)")
	}
	if !strings.HasSuffix(database, "_auth_test") {
		t.Fatal("refusing fixture writes outside a dedicated _auth_test database")
	}
	fixture := func(t *testing.T) (*loginProvider, *LoginHandler, *LoginHandler) {
		t.Helper()
		f := newLoginProvider(t)
		t.Cleanup(func() {
			cleanupCtx, stop := context.WithTimeout(context.Background(), 5*time.Second)
			defer stop()
			if _, err := dbA.ExecContext(cleanupCtx, `DELETE FROM auth.sessions WHERE identity_id IN
				(SELECT id FROM auth.identities WHERE issuer = $1)`, f.issuer); err != nil {
				t.Error("session fixture cleanup failed")
			}
			if _, err := dbA.ExecContext(cleanupCtx, "DELETE FROM auth.identities WHERE issuer = $1", f.issuer); err != nil {
				t.Error("identity fixture cleanup failed")
			}
		})
		return f, f.handler(t, dbA), f.handler(t, dbB)
	}
	counts := func(t *testing.T, issuer string) (int, int) {
		t.Helper()
		var identities, sessions int
		if err := dbA.QueryRowContext(ctx, "SELECT count(*) FROM auth.identities WHERE issuer = $1", issuer).Scan(&identities); err != nil {
			t.Fatal("identity count failed")
		}
		if err := dbA.QueryRowContext(ctx, `SELECT count(*) FROM auth.sessions s JOIN auth.identities i ON i.id = s.identity_id WHERE i.issuer = $1`, issuer).Scan(&sessions); err != nil {
			t.Fatal("session count failed")
		}
		return identities, sessions
	}
	previousSession := func(t *testing.T, f *loginProvider) *http.Cookie {
		t.Helper()
		store := postgres.NewAuthStore(dbA)
		id, err := store.GetOrCreateIdentity(ctx, f.issuer, "previous-human")
		if err != nil {
			t.Fatal("previous identity fixture failed")
		}
		value := oauth2.GenerateVerifier()
		if err := store.CreateSession(ctx, postgres.AuthSession{TokenHash: hashSecret(value), IdentityID: id, Roles: []string{"dispatcher"}, ExpiresAt: time.Now().Add(time.Minute)}); err != nil {
			t.Fatal("previous session fixture failed")
		}
		return &http.Cookie{Name: sessionCookie, Value: value}
	}

	t.Run("session adapter shared read only lookup and failures", func(t *testing.T) {
		f, a, b := fixture(t)
		cookie := previousSession(t, f)
		store := postgres.NewAuthStore(dbA)
		before, err := store.GetSession(ctx, hashSecret(cookie.Value))
		if err != nil {
			t.Fatal("adapter session fixture unavailable")
		}
		want := application.Identity{SubjectID: before.IdentityID, Roles: before.Roles}
		assertFailure := func(h *LoginHandler, r *http.Request, wantError error) {
			t.Helper()
			got, err := h.AuthenticateIncident(r)
			if !reflect.DeepEqual(got, application.Identity{}) || !errors.Is(err, wantError) || err != wantError || errors.Unwrap(err) != nil {
				t.Fatal("adapter failure did not return zero identity and sanitized sentinel")
			}
		}
		reads, writes := 0, 0
		readDB := open(loginFaultConnector{base: stdlib.GetConnector(*config), fail: func(query string, args []driver.NamedValue) bool {
			if !strings.HasPrefix(strings.TrimSpace(query), "SELECT") {
				writes++
				return true
			}
			reads++
			if len(args) != 1 {
				t.Error("incorrect session query arguments")
				return true
			}
			hash, ok := args[0].Value.([]byte)
			if !ok || !bytes.Equal(hash, hashSecret(cookie.Value)) {
				t.Error("session query did not use the identifier hash")
				return true
			}
			return false
		}})
		reader := f.handler(t, readDB)
		// Include request claims and unrelated malformed login cookies; only the
		// session cookie supplies identity. GET/POST/PATCH all use the same adapter.
		for _, h := range []*LoginHandler{a, b, reader} {
			for _, method := range []string{"GET", "POST", "PATCH"} {
				r := loginRequest(method, "/v1/incident?identityId=untrusted&roles=administrator", cookie, &http.Cookie{Name: loginCookie, Value: "invalid"})
				r.Header.Del("Origin")
				headers := r.Header.Clone()
				got, err := h.AuthenticateIncident(r)
				if err != nil || !reflect.DeepEqual(got, want) || !reflect.DeepEqual(r.Header, headers) {
					t.Fatal("adapter changed identity/roles or request cookies")
				}
			}
		}
		if reads != 3 || writes != 0 {
			t.Fatal("authentication performed writes or unexpected lookups")
		}
		after, err := store.GetSession(ctx, hashSecret(cookie.Value))
		if err != nil || !reflect.DeepEqual(before, after) {
			t.Fatal("authentication changed stored session")
		}
		_, sessions := counts(t, f.issuer)
		if sessions != 1 {
			t.Fatal("authentication created or deleted sessions")
		}

		canceled, cancel := context.WithCancel(ctx)
		cancel()
		r := loginRequest("GET", "/v1/incident", cookie).WithContext(canceled)
		assertFailure(b, r, ErrAuthUnavailable)
		failDB := open(loginFaultConnector{base: stdlib.GetConnector(*config), fail: func(query string, _ []driver.NamedValue) bool { return strings.Contains(query, "SELECT token_hash") }})
		assertFailure(f.handler(t, failDB), loginRequest("GET", "/v1/incident", cookie), ErrAuthUnavailable)
		unknown := &http.Cookie{Name: sessionCookie, Value: oauth2.GenerateVerifier()}
		assertFailure(a, loginRequest("GET", "/v1/incident", unknown), ErrUnauthenticated)
		if err := postgres.NewAuthStore(dbB).DeleteSession(ctx, hashSecret(cookie.Value)); err != nil {
			t.Fatal("adapter revocation fixture failed")
		}
		assertFailure(a, loginRequest("GET", "/v1/incident", cookie), ErrUnauthenticated)
		assertFailure(b, loginRequest("GET", "/v1/incident", cookie), ErrUnauthenticated)

		empty := &http.Cookie{Name: sessionCookie, Value: oauth2.GenerateVerifier()}
		if err := store.CreateSession(ctx, postgres.AuthSession{TokenHash: hashSecret(empty.Value), IdentityID: before.IdentityID, ExpiresAt: time.Now().Add(time.Minute)}); err != nil {
			t.Fatal("empty roles fixture failed")
		}
		got, err := b.AuthenticateIncident(loginRequest("GET", "/v1/incident", empty))
		if err != nil || got.SubjectID != before.IdentityID || !reflect.DeepEqual(got.Roles, []string{}) {
			t.Fatal("empty roles were not preserved")
		}
		if _, err := dbA.ExecContext(ctx, "UPDATE auth.sessions SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes', expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute' WHERE token_hash = $1", hashSecret(empty.Value)); err != nil {
			t.Fatal("expired adapter fixture failed")
		}
		assertFailure(a, loginRequest("GET", "/v1/incident", empty), ErrUnauthenticated)
	})

	t.Run("shared session me and logout lifecycle", func(t *testing.T) {
		f, a, b := fixture(t)
		cookie := previousSession(t, f)
		store := postgres.NewAuthStore(dbA)
		before, err := store.GetSession(ctx, hashSecret(cookie.Value))
		if err != nil {
			t.Fatal("session fixture unavailable")
		}
		pending := &http.Cookie{Name: loginCookie, Value: oauth2.GenerateVerifier()}
		request := func(h *LoginHandler, method, path string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
			t.Helper()
			r := loginRequest(method, path, cookies...)
			if method == "GET" {
				r.Header.Del("Origin")
			}
			w := httptest.NewRecorder()
			h.Routes().ServeHTTP(w, r)
			if w.Header().Get("Cache-Control") != "no-store" || w.Header().Get("Referrer-Policy") != "no-referrer" {
				t.Fatal("missing privacy headers")
			}
			return w
		}
		w := request(b, "GET", "/auth/me", cookie)
		expected := `{"identityId":"` + before.IdentityID.String() + `","roles":["dispatcher"]}` + "\n"
		if w.Code != 200 || w.Body.String() != expected || w.Header().Get("Content-Type") != "application/json; charset=utf-8" || len(w.Header().Values("Set-Cookie")) != 0 {
			t.Fatal("incorrect me response, sensitive extra fields or cookie mutation")
		}
		after, err := store.GetSession(ctx, hashSecret(cookie.Value))
		if err != nil || !after.ExpiresAt.Equal(before.ExpiresAt) {
			t.Fatal("me changed session lifetime")
		}
		for _, query := range []string{"SELECT token_hash", "DELETE FROM auth.sessions"} {
			faultDB := open(loginFaultConnector{base: stdlib.GetConnector(*config), fail: func(sql string, _ []driver.NamedValue) bool { return strings.Contains(sql, query) }})
			method, path := "GET", "/auth/me"
			if strings.HasPrefix(query, "DELETE") {
				method, path = "POST", "/auth/logout"
			}
			failed := request(f.handler(t, faultDB), method, path, cookie, pending)
			if failed.Code != 503 || failed.Body.String() != "Service Unavailable\n" || len(failed.Header().Values("Set-Cookie")) != 0 {
				t.Fatal("database failure leaked details or changed cookies")
			}
			if _, err := store.GetSession(ctx, hashSecret(cookie.Value)); err != nil {
				t.Fatal("failed deletion removed session")
			}
		}
		w = request(b, "POST", "/auth/logout", cookie, pending)
		if w.Code != 204 || w.Body.Len() != 0 || len(w.Result().Cookies()) != 1 || responseCookie(w, loginCookie) != nil {
			t.Fatal("logout changed pending login or returned a body")
		}
		checkCookie(t, responseCookie(w, sessionCookie), true)
		if _, err := store.GetSession(ctx, hashSecret(cookie.Value)); !errors.Is(err, postgres.ErrAuthNotFound) {
			t.Fatal("logout did not revoke session")
		}
		if w := request(a, "GET", "/auth/me", cookie); w.Code != 401 || len(w.Result().Cookies()) != 0 {
			t.Fatal("revoked session authenticated")
		}
		w = request(a, "POST", "/auth/logout", cookie, pending)
		if w.Code != 204 || len(w.Result().Cookies()) != 1 {
			t.Fatal("repeat logout not idempotent")
		}
		checkCookie(t, responseCookie(w, sessionCookie), true)
		w = request(a, "POST", "/auth/logout", pending)
		if w.Code != 204 || len(w.Result().Cookies()) != 0 {
			t.Fatal("absent-session logout changed pending cookie")
		}

		// Empty roles remain an array; browser-supplied identity/roles cannot override them.
		empty := &http.Cookie{Name: sessionCookie, Value: oauth2.GenerateVerifier()}
		if err := store.CreateSession(ctx, postgres.AuthSession{TokenHash: hashSecret(empty.Value), IdentityID: before.IdentityID, ExpiresAt: time.Now().Add(time.Minute)}); err != nil {
			t.Fatal("empty-role fixture failed")
		}
		w = request(b, "GET", "/auth/me?identityId=untrusted&roles=administrator", empty)
		if w.Code != 200 || w.Body.String() != `{"identityId":"`+before.IdentityID.String()+`","roles":[]}`+"\n" {
			t.Fatal("incorrect empty roles or trusted request claims")
		}
		// Expire only this fixture, without waiting or changing the schema.
		if _, err := dbA.ExecContext(ctx, "UPDATE auth.sessions SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes', expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute' WHERE token_hash = $1", hashSecret(empty.Value)); err != nil {
			t.Fatal("expiry fixture failed")
		}
		if w := request(b, "GET", "/auth/me", empty); w.Code != 401 || len(w.Result().Cookies()) != 0 {
			t.Fatal("expired session authenticated")
		}
		w = request(b, "POST", "/auth/logout", empty)
		if w.Code != 204 {
			t.Fatal("expired-session logout failed")
		}
		checkCookie(t, responseCookie(w, sessionCookie), true)
		unknown := &http.Cookie{Name: sessionCookie, Value: oauth2.GenerateVerifier()}
		if w := request(a, "GET", "/auth/me", unknown); w.Code != 401 {
			t.Fatal("unknown session authenticated")
		}
	})

	t.Run("cross handler login PKCE nonce binding and shared session", func(t *testing.T) {
		f, a, b := fixture(t)
		binding, callback, authorization := f.begin(t, a, dbA)
		var browserHash, nonceHash []byte
		var verifier string
		var created, expires time.Time
		if err := dbA.QueryRowContext(ctx, `SELECT browser_hash, nonce_hash, pkce_verifier, created_at, expires_at
			FROM auth.login_transactions WHERE state_hash = $1`, hashSecret(authorization.Get("state"))).Scan(&browserHash, &nonceHash, &verifier, &created, &expires); err != nil {
			t.Fatal("transaction was not persisted")
		}
		if !bytes.Equal(browserHash, hashSecret(binding.Value)) || !bytes.Equal(nonceHash, hashSecret(authorization.Get("nonce"))) ||
			oauth2.S256ChallengeFromVerifier(verifier) != authorization.Get("code_challenge") || expires.Sub(created) > 5*time.Minute || expires.Sub(created) < 295*time.Second {
			t.Fatal("transaction state, nonce, binding, S256 or lifetime mismatch")
		}
		unique := map[string]bool{authorization.Get("state"): true, authorization.Get("nonce"): true, binding.Value: true, verifier: true}
		if len(unique) != 4 {
			t.Fatal("login secrets were reused")
		}
		w := httptest.NewRecorder()
		b.Routes().ServeHTTP(w, loginRequest("GET", callback, binding))
		if w.Code != 303 || w.Header().Get("Location") != testOrigin+"/" {
			t.Fatal("callback did not complete with a fixed redirect")
		}
		session := responseCookie(w, sessionCookie)
		checkCookie(t, session, false)
		checkCookie(t, responseCookie(w, loginCookie), true)
		stored, err := postgres.NewAuthStore(dbA).GetSession(ctx, hashSecret(session.Value))
		if err != nil || len(stored.Roles) != 1 || stored.Roles[0] != "dispatcher" || stored.ExpiresAt.After(time.Now().Add(61*time.Second)) {
			t.Fatal("session was not shared or did not respect token validity")
		}
		id, err := postgres.NewAuthStore(dbB).GetOrCreateIdentity(ctx, f.issuer, "login-human")
		if err != nil || id != stored.IdentityID {
			t.Fatal("session did not use stable application UUID mapping")
		}
		me := httptest.NewRecorder()
		a.Routes().ServeHTTP(me, loginRequest("GET", "/auth/me", session))
		if me.Code != 200 || me.Body.String() != `{"identityId":"`+id.String()+`","roles":["dispatcher"]}`+"\n" {
			t.Fatal("callback session unavailable through another handler's me route")
		}
		f.mu.Lock()
		exchanges := f.exchanges
		f.mu.Unlock()
		if exchanges != 1 {
			t.Fatal("confidential-client exchange was not performed exactly once")
		}
		replay := httptest.NewRecorder()
		a.Callback(replay, loginRequest("GET", callback, binding))
		if replay.Code != 400 || len(replay.Result().Cookies()) != 0 {
			t.Fatal("callback replay accepted or cleared a pending cookie")
		}
		f.mu.Lock()
		exchanges = f.exchanges
		f.mu.Unlock()
		if exchanges != 1 {
			t.Fatal("replay reached token exchange")
		}
	})
	t.Run("concurrent duplicate callbacks", func(t *testing.T) {
		f, a, b := fixture(t)
		binding, callback, _ := f.begin(t, a, dbA)
		results := make(chan *httptest.ResponseRecorder, 16)
		start := make(chan struct{})
		for i := 0; i < 16; i++ {
			h := a
			if i%2 != 0 {
				h = b
			}
			go func() {
				<-start
				w := httptest.NewRecorder()
				h.Callback(w, loginRequest("GET", callback, binding))
				results <- w
			}()
		}
		close(start)
		successes := 0
		for i := 0; i < 16; i++ {
			w := <-results
			if w.Code == 303 {
				successes++
			} else if w.Code != 400 || len(w.Result().Cookies()) != 0 {
				t.Error("losing callback was not safely rejected")
			}
		}
		f.mu.Lock()
		exchanges := f.exchanges
		f.mu.Unlock()
		_, sessions := counts(t, f.issuer)
		if successes != 1 || exchanges != 1 || sessions != 1 {
			t.Fatal("duplicate callbacks produced extra exchanges or sessions")
		}
	})

	t.Run("failed callback preserves another pending login", func(t *testing.T) {
		f, a, b := fixture(t)
		_, callbackA, _ := f.begin(t, a, dbA)
		bindingB, callbackB, _ := f.begin(t, b, dbA)
		w := httptest.NewRecorder()
		b.Callback(w, loginRequest("GET", callbackA, bindingB))
		if w.Code != 400 || len(w.Result().Cookies()) != 0 {
			t.Fatal("mismatched callback changed pending login")
		}
		w = httptest.NewRecorder()
		a.Callback(w, loginRequest("GET", callbackB, bindingB))
		if w.Code != 303 {
			t.Fatal("pending browser-bound transaction was consumed by wrong callback")
		}
	})

	t.Run("expired login transaction", func(t *testing.T) {
		f, a, b := fixture(t)
		binding, callback, authorization := f.begin(t, a, dbA)
		if _, err := dbA.ExecContext(ctx, `UPDATE auth.login_transactions SET created_at = CURRENT_TIMESTAMP - INTERVAL '2 minutes',
			expires_at = CURRENT_TIMESTAMP - INTERVAL '1 minute' WHERE state_hash = $1`, hashSecret(authorization.Get("state"))); err != nil {
			t.Fatal("expiry fixture failed")
		}
		w := httptest.NewRecorder()
		b.Callback(w, loginRequest("GET", callback, binding))
		if w.Code != 400 || len(w.Result().Cookies()) != 0 {
			t.Fatal("expired transaction accepted or cookie cleared")
		}
		f.mu.Lock()
		exchanges := f.exchanges
		f.mu.Unlock()
		if exchanges != 0 {
			t.Fatal("expired transaction reached exchange")
		}
	})

	for _, mode := range []string{"invalid token", "provider error", "authorization denied"} {
		t.Run(mode, func(t *testing.T) {
			f, a, b := fixture(t)
			f.mu.Lock()
			f.mode = mode
			f.mu.Unlock()
			binding, callback, _ := f.begin(t, a, dbA)
			status := 401
			if mode == "provider error" {
				status = 502
			}
			if mode == "authorization denied" {
				status = 400
				parsed, _ := url.Parse(callback)
				q := parsed.Query()
				q.Del("code")
				q.Set("error", "access_denied")
				q.Set("error_description", "sensitive provider diagnostics")
				parsed.RawQuery = q.Encode()
				callback = parsed.String()
			}
			w := httptest.NewRecorder()
			b.Callback(w, loginRequest("GET", callback, binding))
			identities, sessions := counts(t, f.issuer)
			if w.Code != status || identities != 0 || sessions != 0 || len(w.Result().Cookies()) != 0 || strings.Contains(w.Body.String(), "sensitive") {
				t.Fatal("failed authentication created identity/session, exposed diagnostics or cleared cookies")
			}
			replay := httptest.NewRecorder()
			a.Callback(replay, loginRequest("GET", callback, binding))
			if replay.Code != 400 || len(replay.Result().Cookies()) != 0 {
				t.Fatal("failed exchange left a reusable transaction")
			}
			f.mu.Lock()
			exchanges := f.exchanges
			f.mu.Unlock()
			expected := 1
			if mode == "authorization denied" {
				expected = 0
			}
			if exchanges != expected {
				t.Fatal("exchange was retried or authorization error was exchanged")
			}
		})
	}

	t.Run("session replacement and thirty minute cap", func(t *testing.T) {
		f, a, b := fixture(t)
		f.mu.Lock()
		f.mode = "long tokens"
		f.mu.Unlock()
		previous := previousSession(t, f)
		binding, callback, _ := f.begin(t, a, dbA)
		before := time.Now()
		w := httptest.NewRecorder()
		b.Callback(w, loginRequest("GET", callback, binding, previous))
		if w.Code != 303 {
			t.Fatal("replacement callback failed")
		}
		cookie := responseCookie(w, sessionCookie)
		checkCookie(t, cookie, false)
		if cookie.Value == previous.Value {
			t.Fatal("session fixation: previous identifier reused")
		}
		store := postgres.NewAuthStore(dbA)
		if _, err := store.GetSession(ctx, hashSecret(previous.Value)); !errors.Is(err, postgres.ErrAuthNotFound) {
			t.Fatal("previous session was not revoked")
		}
		stored, err := store.GetSession(ctx, hashSecret(cookie.Value))
		if err != nil || stored.ExpiresAt.Before(before.Add(30*time.Minute)) || stored.ExpiresAt.After(time.Now().Add(30*time.Minute)) {
			t.Fatal("session did not respect thirty minute cap")
		}
		_, sessions := counts(t, f.issuer)
		if sessions != 1 {
			t.Fatal("replacement left multiple sessions")
		}
	})

	t.Run("failed consumption retains retryable transaction", func(t *testing.T) {
		f, a, _ := fixture(t)
		binding, callback, _ := f.begin(t, a, dbA)
		broken := f.handler(t, closedLoginDB(t))
		w := httptest.NewRecorder()
		broken.Callback(w, loginRequest("GET", callback, binding))
		if w.Code != 503 || len(w.Result().Cookies()) != 0 {
			t.Fatal("database failure issued or cleared cookies")
		}
		w = httptest.NewRecorder()
		a.Callback(w, loginRequest("GET", callback, binding))
		if w.Code != 303 {
			t.Fatal("failed consumption lost an unconsumed transaction")
		}
	})

	for _, mode := range []string{"identity failure", "session creation failure", "revocation failure", "commit rolled back", "commit acknowledgement lost"} {
		t.Run(mode, func(t *testing.T) {
			f, a, _ := fixture(t)
			previous := previousSession(t, f)
			binding, callback, _ := f.begin(t, a, dbA)
			var capturedNew []byte
			deleteAttempts := 0
			insertAttempts, commitAttempts := 0, 0
			beforeCommitConsistent := false
			w := httptest.NewRecorder()
			fault := func(query string, args []driver.NamedValue) bool {
				if strings.Contains(query, "INSERT INTO auth.identities") && mode == "identity failure" {
					return true
				}
				if strings.Contains(query, "INSERT INTO auth.sessions") {
					insertAttempts++
					capturedNew = append([]byte(nil), args[0].Value.([]byte)...)
					return mode == "session creation failure"
				}
				if strings.Contains(query, "DELETE FROM auth.sessions") {
					deleteAttempts++
					if bytes.Equal(args[0].Value.([]byte), hashSecret(previous.Value)) {
						var oldCount, newCount int
						err := dbA.QueryRowContext(ctx, `SELECT count(*) FILTER (WHERE token_hash = $1),
							count(*) FILTER (WHERE token_hash = $2) FROM auth.sessions WHERE token_hash IN ($1, $2)`,
							hashSecret(previous.Value), capturedNew).Scan(&oldCount, &newCount)
						beforeCommitConsistent = err == nil && oldCount == 1 && newCount == 0
						return mode == "revocation failure"
					}
					t.Error("replacement attempted a compensating session deletion")
				}
				return false
			}
			commit := func(tx driver.Tx) error {
				commitAttempts++
				var oldCount, newCount int
				if err := dbA.QueryRowContext(ctx, `SELECT count(*) FILTER (WHERE token_hash = $1),
					count(*) FILTER (WHERE token_hash = $2) FROM auth.sessions WHERE token_hash IN ($1, $2)`,
					hashSecret(previous.Value), capturedNew).Scan(&oldCount, &newCount); err != nil || oldCount != 1 || newCount != 0 {
					t.Error("uncommitted deletion/insertion leaked to another connection")
				}
				if len(w.Header().Values("Set-Cookie")) != 0 {
					t.Error("cookie issued before commit")
				}
				if mode == "commit acknowledgement lost" {
					if err := tx.Commit(); err != nil {
						return err
					}
					return errors.New("injected lost commit acknowledgement")
				}
				if err := tx.Rollback(); err != nil {
					return err
				}
				return errors.New("injected commit failure with confirmed rollback")
			}
			faultDB := sql.OpenDB(loginFaultConnector{base: stdlib.GetConnector(*config), fail: fault, commit: commit})
			t.Cleanup(func() { _ = faultDB.Close() })
			b := f.handler(t, faultDB)
			b.Callback(w, loginRequest("GET", callback, binding, previous))
			if w.Code != 503 || len(w.Result().Cookies()) != 0 {
				t.Fatal("replacement failure issued or cleared cookies")
			}
			_, oldErr := postgres.NewAuthStore(dbA).GetSession(ctx, hashSecret(previous.Value))
			_, sessions := counts(t, f.issuer)
			if mode == "commit acknowledgement lost" {
				_, newErr := postgres.NewAuthStore(dbA).GetSession(ctx, capturedNew)
				if !errors.Is(oldErr, postgres.ErrAuthNotFound) || newErr != nil || sessions != 1 {
					t.Fatal("lost acknowledgement was incorrectly treated as rollback")
				}
			} else {
				if oldErr != nil || sessions != 1 {
					t.Fatal("definite failure did not preserve previous session")
				}
				if len(capturedNew) > 0 {
					if _, err := postgres.NewAuthStore(dbA).GetSession(ctx, capturedNew); !errors.Is(err, postgres.ErrAuthNotFound) {
						t.Fatal("failed transaction left a new session")
					}
				}
			}
			if mode == "revocation failure" || strings.HasPrefix(mode, "commit ") {
				if !beforeCommitConsistent || insertAttempts != 1 || deleteAttempts != 1 {
					t.Fatal("transaction visibility or operation count incorrect")
				}
			} else if deleteAttempts != 0 {
				t.Fatal("revocation attempted after failed creation")
			}
			expectedCommits := 0
			if strings.HasPrefix(mode, "commit ") {
				expectedCommits = 1
			}
			if commitAttempts != expectedCommits || insertAttempts > 1 {
				t.Fatal("failed transaction was retried or committed unexpectedly")
			}
		})
	}
}
