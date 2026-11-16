package http

import (
	"net/http"

	"github.com/google/uuid"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
)

// PlaceholderAuthenticator is a temporary authenticator used until login is implemented.
// It does not check the request at all: every caller is treated as the same test dispatcher.
//
// TODO(auth team): replace with a real authenticator (e.g. reading a token or client
// certificate from the request) and pass it to NewIncidentHandler in cmd/api/main.go.
type PlaceholderAuthenticator struct{}

// placeholderDispatcherID is the fixed user ID given to every request.
var placeholderDispatcherID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

func (PlaceholderAuthenticator) AuthenticateIncident(*http.Request) (application.Identity, error) {
	return application.Identity{
		SubjectID: placeholderDispatcherID,
		Roles:     []string{"dispatcher"},
	}, nil
}

// Compile-time check: the build fails if PlaceholderAuthenticator stops matching the interface.
var _ IncidentAuthorizer = PlaceholderAuthenticator{}
