package http

import (
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"net/http"
)

// IncidentAuthenticator authenticates HTTP requests. It must not make
// role-dependent authorization decisions.
type IncidentAuthenticator interface {
	AuthenticateIncident(*http.Request) (application.Identity, error)
}

func 
