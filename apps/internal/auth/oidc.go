// Package auth provides human OIDC login and token validation independently of
// application authorization policies.
package auth

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var ErrInvalidTokens = errors.New("invalid OIDC token pair")

// OIDCConfig is supplied by the composition root. RootCAs selects trusted roots;
// nil uses system trust. Validation never needs a client secret.
type OIDCConfig struct {
	IssuerURL   string
	ClientID    string
	APIAudience string
	HumanRoles  []string
	RootCAs     *x509.CertPool
}

// ValidatedIdentity is ready for issuer/subject-to-UUID mapping later. Roles are
// an admitted human-role snapshot, not permission to operate on any incident.
type ValidatedIdentity struct {
	Issuer    string
	Subject   string
	Roles     []string
	ExpiresAt time.Time
}

type OIDC struct {
	idVerifier     *oidc.IDTokenVerifier
	accessVerifier *oidc.IDTokenVerifier
	issuer         string
	clientID       string
	humanRoles     map[string]bool
	endpoint       oauth2.Endpoint
	httpClient     *http.Client
}

// NewOIDC discovers one exact HTTPS issuer. Verifiers share the provider's JWKS
// cache; keys are fetched on demand and refreshed by go-oidc for key rotation.
func NewOIDC(ctx context.Context, config OIDCConfig) (*OIDC, error) {
	issuer, err := httpsURL(config.IssuerURL)
	if err != nil || config.ClientID == "" || config.APIAudience == "" ||
		config.ClientID == config.APIAudience || len(config.HumanRoles) == 0 {
		return nil, errors.New("invalid OIDC configuration")
	}
	roles := make(map[string]bool, len(config.HumanRoles))
	for _, role := range config.HumanRoles {
		if role == "" {
			return nil, errors.New("empty OIDC human role")
		}
		roles[role] = true
	}
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12, RootCAs: config.RootCAs}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	provider, err := oidc.NewProvider(oidc.ClientContext(ctx, client), config.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("discover OIDC issuer: %w", err)
	}
	var metadata struct {
		JWKSURL string `json:"jwks_uri"`
	}
	if err := provider.Claims(&metadata); err != nil {
		return nil, errors.New("invalid OIDC discovery metadata")
	}
	jwks, err := httpsURL(metadata.JWKSURL)
	if err != nil || jwks.Host != issuer.Host {
		return nil, errors.New("OIDC signing keys must use the issuer's HTTPS origin")
	}
	return &OIDC{
		idVerifier:     provider.Verifier(&oidc.Config{ClientID: config.ClientID, SupportedSigningAlgs: []string{oidc.RS256}}),
		accessVerifier: provider.Verifier(&oidc.Config{ClientID: config.APIAudience, SupportedSigningAlgs: []string{oidc.RS256}}),
		issuer:         config.IssuerURL,
		clientID:       config.ClientID,
		humanRoles:     roles,
		endpoint:       provider.Endpoint(),
		httpClient:     client,
	}, nil
}

type tokenClaims struct {
	AuthorizedParty string `json:"azp"`
	Type            string `json:"typ"`
	NotBefore       *int64 `json:"nbf"`
	RealmAccess     struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

// Validate accepts tokens from the same future confidential-client code exchange.
// expectedNonceHash must come from the consumed, browser-bound login transaction,
// never from request input. No token or raw library token diagnostic is logged
// or included in returned errors, and no partial identity is returned on failure.
func (o *OIDC) Validate(ctx context.Context, rawIDToken, rawAccessToken string, expectedNonceHash []byte) (ValidatedIdentity, error) {
	if len(expectedNonceHash) != sha256.Size {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	id, err := o.idVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return ValidatedIdentity{}, fmt.Errorf("%w: ID token verification", ErrInvalidTokens)
	}
	// For Keycloak JWT access tokens, reuse the library's signature/issuer/time
	// checks with the API audience, then enforce the access-token profile below.
	access, err := o.accessVerifier.Verify(ctx, rawAccessToken)
	if err != nil {
		return ValidatedIdentity{}, fmt.Errorf("%w: access token verification", ErrInvalidTokens)
	}
	var idClaims, accessClaims tokenClaims
	if id.Claims(&idClaims) != nil || access.Claims(&accessClaims) != nil {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	now := time.Now()
	if !validTokenTimes(id, idClaims.NotBefore, now) || !validTokenTimes(access, accessClaims.NotBefore, now) {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	// Keep literal issuer equality, including providers for which go-oidc has
	// historical issuer-compatibility exceptions.
	if id.Issuer != o.issuer || access.Issuer != o.issuer || id.Subject == "" || id.Subject != access.Subject ||
		(idClaims.AuthorizedParty != "" && idClaims.AuthorizedParty != o.clientID) ||
		(len(id.Audience) > 1 && idClaims.AuthorizedParty != o.clientID) ||
		accessClaims.AuthorizedParty != o.clientID || accessClaims.Type != "Bearer" {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	nonceHash := sha256.Sum256([]byte(id.Nonce))
	if id.Nonce == "" || subtle.ConstantTimeCompare(nonceHash[:], expectedNonceHash) != 1 {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	// This library method checks only at_hash, not the access JWT's signature.
	if id.AccessTokenHash != "" && id.VerifyAccessToken(rawAccessToken) != nil {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	var roles []string
	for _, role := range accessClaims.RealmAccess.Roles {
		if o.humanRoles[role] && !slices.Contains(roles, role) {
			roles = append(roles, role)
		}
	}
	if len(roles) == 0 {
		return ValidatedIdentity{}, ErrInvalidTokens
	}
	expiresAt := id.Expiry
	if access.Expiry.Before(expiresAt) {
		expiresAt = access.Expiry
	}
	return ValidatedIdentity{Issuer: id.Issuer, Subject: id.Subject, Roles: roles, ExpiresAt: expiresAt}, nil
}

func validTokenTimes(token *oidc.IDToken, notBefore *int64, now time.Time) bool {
	// go-oidc permits five minutes of nbf skew. This profile permits only 30s
	// for nbf/iat, requires iat, and never extends the signed expiration.
	latest := now.Add(30 * time.Second)
	// Clock tolerance cannot make inconsistent signed validity times acceptable.
	if token.IssuedAt.After(token.Expiry) || (notBefore != nil && time.Unix(*notBefore, 0).After(token.Expiry)) {
		return false
	}
	return now.Before(token.Expiry) && !token.IssuedAt.IsZero() &&
		!token.IssuedAt.After(latest) &&
		(notBefore == nil || !time.Unix(*notBefore, 0).After(latest))
}

func httpsURL(raw string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil ||
		parsed.RawQuery != "" || parsed.ForceQuery || parsed.Fragment != "" || parsed.Opaque != "" {
		return nil, errors.New("expected an HTTPS URL without credentials, query or fragment")
	}
	return parsed, nil
}
