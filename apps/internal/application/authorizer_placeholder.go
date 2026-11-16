package application

import (
	"context"

	"github.com/google/uuid"
)

// PlaceholderAuthorizer is a temporary authorizer used until role-based permissions exist.
// It allows every operation for every caller.
//
// TODO(auth team): replace with a real authorizer that checks identity.Roles
// (e.g. only dispatchers may post incidents) and pass it to NewIncidentService in cmd/api/main.go.
type PlaceholderAuthorizer struct{}

func (PlaceholderAuthorizer) AuthorizePostIncident(context.Context, Identity) error {
	return nil
}

func (PlaceholderAuthorizer) AuthorizePatchIncident(context.Context, Identity, uuid.UUID) error {
	return nil
}

func (PlaceholderAuthorizer) AuthorizeGetIncident(context.Context, Identity, uuid.UUID) error {
	return nil
}

// Compile-time check: the build fails if PlaceholderAuthorizer stops matching IncidentAuthorizer.
var _ IncidentAuthorizer = PlaceholderAuthorizer{}
