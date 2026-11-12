package domain

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
)

var ErrInvalidIncident = errors.New("invalid incident")

type (
	IncidentID       uuid.UUID
	IncidentType     int16
	IncidentPriority int8
	IncidentStatus   int8
)

// Option separates an omitted property from a property explicitly set to its zero value.
// A nil Value denotes the property's zero value. Set=false means the property is
// uninitialized on an Incident or must not be applied by a patch.
type Option[T any] struct {
	Value *T
	Set   bool
}

type Incident struct {
	ID          IncidentID
	Type        Option[IncidentType]
	Status      Option[IncidentStatus]
	Priority    Option[IncidentPriority]
	Location    Location
	Description Option[string]
}

type IncidentPatch struct {
	Type        Option[IncidentType]
	Status      Option[IncidentStatus]
	Priority    Option[IncidentPriority]
	Location    Location
	Description Option[string]
}

type Location struct {
	Latitude            Option[float64]
	Longitude           Option[float64]
	Address             Option[string]
	LocationDescription Option[string]
}

// ValidateIncident is the single place for incident field validation. Value-specific
// rules can be added to the validate functions as the accepted values are defined.
func ValidateIncident(ctx context.Context, logger *slog.Logger, incident *Incident) error {
	logger.DebugContext(ctx, "validating incident")
	if incident == nil {
		logger.WarnContext(ctx, "incident validation failed", "reason", "incident is nil")
		return ErrInvalidIncident
	}
	if uuid.UUID(incident.ID) == uuid.Nil {
		logger.WarnContext(ctx, "incident validation failed", "reason", "incident ID is required")
		return errors.Join(ErrInvalidIncident, errors.New("incident ID is required"))
	}
	if incident.Type.Set && incident.Type.Value != nil {
		if err := validateIncidentType(*incident.Type.Value); err != nil {
			return errors.Join(ErrInvalidIncident, err)
		}
	}
	if incident.Status.Set && incident.Status.Value != nil {
		if err := validateIncidentStatus(*incident.Status.Value); err != nil {
			return errors.Join(ErrInvalidIncident, err)
		}
	}
	if incident.Priority.Set && incident.Priority.Value != nil {
		if err := validateIncidentPriority(*incident.Priority.Value); err != nil {
			return errors.Join(ErrInvalidIncident, err)
		}
	}
	if err := validateLocation(incident.Location); err != nil {
		logger.WarnContext(ctx, "incident validation failed", "reason", err)
		return errors.Join(ErrInvalidIncident, err)
	}
	return nil
}

func PostIncident(ctx context.Context, logger *slog.Logger, patch IncidentPatch) (*Incident, error) {
	logger.DebugContext(ctx, "posting incident")
	incident := &Incident{
		ID:          IncidentID(uuid.New()),
		Type:        patch.Type,
		Status:      patch.Status,
		Priority:    patch.Priority,
		Location:    patch.Location,
		Description: patch.Description,
	}
	if err := ValidateIncident(ctx, logger, incident); err != nil {
		return nil, err
	}
	return incident, nil
}

func PatchIncident(ctx context.Context, logger *slog.Logger, patch IncidentPatch, incident *Incident) (*Incident, error) {
	logger.InfoContext(ctx, "patching incident")
	if err := ValidateIncident(ctx, logger, incident); err != nil {
		return nil, err
	}

	// Copies current incident and overrides only set values
	patched := *incident
	if patch.Type.Set {
		patched.Type = patch.Type
	}
	if patch.Status.Set {
		patched.Status = patch.Status
	}
	if patch.Priority.Set {
		patched.Priority = patch.Priority
	}
	if patch.Location.Latitude.Set {
		patched.Location.Latitude = patch.Location.Latitude
	}
	if patch.Location.Longitude.Set {
		patched.Location.Longitude = patch.Location.Longitude
	}
	if patch.Location.Address.Set {
		patched.Location.Address = patch.Location.Address
	}
	if patch.Location.LocationDescription.Set {
		patched.Location.LocationDescription = patch.Location.LocationDescription
	}
	if patch.Description.Set {
		patched.Description = patch.Description
	}

	if err := ValidateIncident(ctx, logger, &patched); err != nil {
		return nil, err
	}
	return &patched, nil
}

func GetIncident(ctx context.Context, logger *slog.Logger, incident *Incident) (*Incident, error) {
	logger.InfoContext(ctx, "getting incident")
	if err := ValidateIncident(ctx, logger, incident); err != nil {
		return nil, err
	}
	return incident, nil
}

func validateIncidentType(IncidentType) error {
	return nil
}

func validateIncidentStatus(IncidentStatus) error {
	return nil
}

func validateIncidentPriority(IncidentPriority) error {
	return nil
}

func validateLocation(location Location) error {
	if location.Latitude.Set != location.Longitude.Set {
		return errors.New("latitude and longitude must be set together")
	}
	if location.Latitude.Set && (location.Latitude.Value == nil || location.Longitude.Value == nil) {
		return errors.New("latitude and longitude must both have values")
	}
	return nil
}
