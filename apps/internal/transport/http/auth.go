package http

import (
	"net/http"

	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
)

// IncidentAuthorizer authenticates HTTP requests. It must not make
// role-dependent authorization decisions.
type IncidentAuthorizer interface {
	AuthenticateIncident(*http.Request) (application.Identity, error)
}
