package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	"golang.org/x/oauth2"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrAuthUnavailable = errors.New("authentication unavailable")
)

const (
	loginCookie   = "__Host-ce1-login"
	sessionCookie = "__Host-ce1-session"
	loginLifetime = 5 * time.Minute
	sessionLimit  = 30 * time.Minute
)

// LoginConfig comes from the composition root, not HTTP input. Never log this
// configuration: ClientSecret is used only for confidential-client exchange.
type LoginConfig struct {
	PublicOrigin string
	ClientSecret string
}

type LoginHandler struct {
	oidc   *OIDC
	store  *postgres.AuthStore
	oauth  oauth2.Config
	origin *url.URL
}

func NewLoginHandler(validator *OIDC, store *postgres.AuthStore, config LoginConfig) (*LoginHandler, error) {
	origin, err := httpsURL(config.PublicOrigin)
	if err != nil || origin.Path != "" || config.ClientSecret == "" || validator == nil || store == nil {
		return nil, errors.New("invalid OIDC login configuration")
	}
	issuer, err := httpsURL(validator.issuer)
	if err != nil {
		return nil, errors.New("invalid OIDC login issuer")
	}
	endpoint := validator.endpoint
	for _, raw := range []string{endpoint.AuthURL, endpoint.TokenURL} {
		parsed, err := httpsURL(raw)
		if err != nil || parsed.Host != issuer.Host {
			return nil, errors.New("OIDC login endpoints must use the issuer's HTTPS origin")
		}
	}
	// Explicit authentication avoids OAuth2's automatic style-probing retry of a
	// one-use code. The secret is sent in HTTP Basic auth, never in a URL.
	endpoint.AuthStyle = oauth2.AuthStyleInHeader
	return &LoginHandler{
		oidc: validator, store: store, origin: origin,
		oauth: oauth2.Config{ClientID: validator.clientID, ClientSecret: config.ClientSecret,
			Endpoint: endpoint, RedirectURL: config.PublicOrigin + "/auth/callback", Scopes: []string{oidc.ScopeOpenID}},
	}, nil
}

// Routes can later be mounted by the existing application. No server, session
// authenticator, incident policy or trusted-proxy behavior is installed here.
func (h *LoginHandler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", h.Login)
	mux.HandleFunc("GET /auth/callback", h.Callback)
	mux.HandleFunc("GET /auth/me", h.Me)
	mux.HandleFunc("POST /auth/logout", h.Logout)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Referrer-Policy", "no-referrer")
		mux.ServeHTTP(w, r)
	})
}

// RequireSameOrigin guards future cookie-authenticated mutation routes. It does
// not authenticate the caller and is deliberately not applied to Routes().
func (h *LoginHandler) RequireSameOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origins := r.Header.Values("Origin")
		if r.TLS == nil || r.Host != h.origin.Host ||
			len(origins) != 1 || origins[0] != h.origin.String() {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *LoginHandler) Me(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequest(w, r, http.MethodGet) {
		return
	}
	session, err := h.session(r)
	if err != nil {
		status := http.StatusServiceUnavailable
		if errors.Is(err, ErrUnauthenticated) {
			status = http.StatusUnauthorized
		}
		http.Error(w, http.StatusText(status), status)
		return
	}
	if session.Roles == nil {
		session.Roles = []string{}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// Explicit fields keep persistence and provider identifiers out of responses.
	// A response write failure cannot be repaired with a second HTTP response.
	_ = json.NewEncoder(w).Encode(struct {
		IdentityID string   `json:"identityId"`
		Roles      []string `json:"roles"`
	}{session.IdentityID.String(), session.Roles})
}

// AuthenticateIncident supplies human session identity, not incident permissions
// or authorized device identity. Production wiring must enforce those separately.
func (h *LoginHandler) AuthenticateIncident(r *http.Request) (application.Identity, error) {
	if r.TLS == nil || r.Host != h.origin.Host {
		return application.Identity{}, ErrUnauthenticated
	}
	session, err := h.session(r)
	if err != nil {
		return application.Identity{}, err
	}
	return application.Identity{SubjectID: session.IdentityID, Roles: session.Roles}, nil
}

// session shares cookie validation and read-only lookup with /auth/me. Raw store
// errors never escape this boundary, including through an error wrapping chain.
func (h *LoginHandler) session(r *http.Request) (postgres.AuthSession, error) {
	identifier, err := secretCookie(r, sessionCookie)
	if err != nil || identifier == "" {
		return postgres.AuthSession{}, ErrUnauthenticated
	}
	session, err := h.store.GetSession(r.Context(), hashSecret(identifier))
	if errors.Is(err, postgres.ErrAuthNotFound) {
		return postgres.AuthSession{}, ErrUnauthenticated
	}
	if err != nil {
		return postgres.AuthSession{}, ErrAuthUnavailable
	}
	return session, nil
}

func (h *LoginHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequest(w, r, http.MethodPost) {
		return
	}
	if origins := r.Header.Values("Origin"); len(origins) != 1 || origins[0] != h.origin.String() {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	identifier, err := secretCookie(r, sessionCookie)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if identifier != "" {
		if err := h.store.DeleteSession(r.Context(), hashSecret(identifier)); err != nil {
			// An acknowledgement failure may follow deletion; do not claim success.
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		setSecretCookie(w, sessionCookie, "", time.Unix(1, 0))
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequest(w, r, http.MethodPost) {
		return
	}
	if origins := r.Header.Values("Origin"); len(origins) != 1 || origins[0] != h.origin.String() {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	// GenerateVerifier provides 32 cryptographically random bytes encoded as
	// base64url. That entropy is also appropriate for the other opaque secrets.
	state, nonce := oauth2.GenerateVerifier(), oauth2.GenerateVerifier()
	binding, verifier := oauth2.GenerateVerifier(), oauth2.GenerateVerifier()
	expiresAt := time.Now().Add(loginLifetime)
	transaction := postgres.AuthLoginTransaction{
		StateHash: hashSecret(state), NonceHash: hashSecret(nonce), BrowserHash: hashSecret(binding),
		PKCEVerifier: verifier, ExpiresAt: expiresAt,
	}
	if err := h.store.CreateLoginTransaction(r.Context(), transaction); err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	setSecretCookie(w, loginCookie, binding, expiresAt)
	location := h.oauth.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier), oidc.Nonce(nonce), oauth2.SetAuthURLParam("response_mode", "query"))
	http.Redirect(w, r, location, http.StatusSeeOther)
}

func (h *LoginHandler) Callback(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequest(w, r, http.MethodGet) {
		return
	}
	query, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || !validCallback(query, h.oidc.issuer) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	binding, err := secretCookie(r, loginCookie)
	previous, previousErr := secretCookie(r, sessionCookie)
	if err != nil || binding == "" || previousErr != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	transaction, err := h.store.ConsumeLoginTransaction(r.Context(), hashSecret(query.Get("state")), hashSecret(binding))
	if err != nil {
		status := http.StatusServiceUnavailable
		if errors.Is(err, postgres.ErrAuthNotFound) {
			status = http.StatusBadRequest
		}
		http.Error(w, http.StatusText(status), status)
		return
	}
	// Even a consumed transaction leaves its cookie untouched on failure: an
	// in-flight callback must not erase a subsequently started pending login.
	if query.Has("error") {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := oidc.ClientContext(r.Context(), h.oidc.httpClient)
	tokens, err := h.oauth.Exchange(ctx, query.Get("code"), oauth2.VerifierOption(transaction.PKCEVerifier))
	if err != nil {
		// OAuth errors can contain provider response bodies. Do not expose/log them.
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	rawID, ok := tokens.Extra("id_token").(string)
	if !ok || !strings.EqualFold(tokens.TokenType, "Bearer") {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		return
	}
	identity, err := h.oidc.Validate(r.Context(), rawID, tokens.AccessToken, transaction.NonceHash)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	id, err := h.store.GetOrCreateIdentity(r.Context(), identity.Issuer, identity.Subject)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	identifier := oauth2.GenerateVerifier()
	expiresAt := time.Now().Add(sessionLimit)
	if identity.ExpiresAt.Before(expiresAt) {
		expiresAt = identity.ExpiresAt
	}
	session := postgres.AuthSession{TokenHash: hashSecret(identifier), IdentityID: id, Roles: identity.Roles, ExpiresAt: expiresAt}
	var previousHash []byte
	if previous != "" {
		previousHash = hashSecret(previous)
	}
	// Issue cookies only after the store acknowledges the atomic commit. A
	// commit error is not retried, even when its outcome is uncertain.
	if err := h.store.ReplaceSession(r.Context(), session, previousHash); err != nil {
		http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
		return
	}
	setSecretCookie(w, loginCookie, "", time.Unix(1, 0))
	setSecretCookie(w, sessionCookie, identifier, expiresAt)
	http.Redirect(w, r, h.origin.String()+"/", http.StatusSeeOther)
}

func (h *LoginHandler) checkRequest(w http.ResponseWriter, r *http.Request, method string) bool {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Referrer-Policy", "no-referrer")
	// Forwarded headers are deliberately insufficient; trusted ingress TLS
	// termination needs a separate deployment contract.
	if r.TLS == nil || r.Host != h.origin.Host {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return false
	}
	if r.Method != method {
		w.Header().Set("Allow", method)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func validCallback(query url.Values, issuer string) bool {
	for key, values := range query {
		switch key {
		case "state", "code", "error", "error_description", "error_uri", "iss", "session_state":
		default:
			return false
		}
		if len(values) != 1 {
			return false
		}
	}
	if !validSecret(query.Get("state")) || query.Has("code") == query.Has("error") ||
		(query.Has("code") && query.Get("code") == "") || (query.Has("error") && query.Get("error") == "") ||
		(query.Has("iss") && query.Get("iss") != issuer) ||
		(!query.Has("error") && (query.Has("error_description") || query.Has("error_uri"))) {
		return false
	}
	return true
}

func hashSecret(secret string) []byte {
	hash := sha256.Sum256([]byte(secret))
	return hash[:]
}

func validSecret(secret string) bool {
	if len(secret) != 43 {
		return false
	}
	decoded, err := base64.RawURLEncoding.Strict().DecodeString(secret)
	return err == nil && len(decoded) == 32
}

func secretCookie(r *http.Request, name string) (string, error) {
	var value string
	// Inspect raw pairs: net/http can silently discard malformed values, hiding
	// an invalid duplicate beside an otherwise valid authentication cookie.
	for _, header := range r.Header.Values("Cookie") {
		for _, pair := range strings.Split(header, ";") {
			key, candidate, _ := strings.Cut(strings.TrimSpace(pair), "=")
			if strings.TrimSpace(key) == name {
				if value != "" || !validSecret(candidate) {
					return "", errors.New("invalid authentication cookie")
				}
				value = candidate
			}
		}
	}
	return value, nil
}

func setSecretCookie(w http.ResponseWriter, name, value string, expiresAt time.Time) {
	maxAge := int(time.Until(expiresAt).Seconds())
	if value == "" {
		maxAge = -1
	}
	http.SetCookie(w, &http.Cookie{Name: name, Value: value, Path: "/", Secure: true, HttpOnly: true,
		SameSite: http.SameSiteLaxMode, Expires: expiresAt.UTC(), MaxAge: maxAge})
}
