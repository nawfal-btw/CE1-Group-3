package http

import (
	"context"
	"net/http"

	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
)

type identityContextKey struct{}

// WithIdentity stores the identity established by authentication middleware.
// Credential validation belongs to the concrete authentication adapter.
func WithIdentity(r *http.Request, identity application.Identity) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), identityContextKey{}, identity))
}

func identityFromContext(r *http.Request) (application.Identity, bool) {
	identity, ok := r.Context().Value(identityContextKey{}).(application.Identity)
	return identity, ok
}
