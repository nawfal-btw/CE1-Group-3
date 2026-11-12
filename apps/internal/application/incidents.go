package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
)

var ErrInvalidIncident = domain.ErrInvalidIncident
var ErrUnauthorized = errors.New("unauthorized")

// Identity is the authenticated caller information needed for application
// authorization. Authentication mechanisms belong to the transport layer.
type Identity struct {
	SubjectID uuid.UUID
	Roles     []string
}

// Field keeps the caller's intent intact while data crosses into the domain.
type Field[T any] struct {
	Value *T
	Set   bool
}

type IncidentInput struct {
	Type        Field[int16]
	Status      Field[int8]
	Priority    Field[int8]
	Location    LocationInput
	Description Field[string]
}

type Incident struct {
	ID          uuid.UUID
	Type        Field[int16]
	Status      Field[int8]
	Priority    Field[int8]
	Location    LocationInput
	Description Field[string]
}

type LocationInput struct {
	Latitude            Field[float64]
	Longitude           Field[float64]
	Address             Field[string]
	LocationDescription Field[string]
}

type IncidentRepository interface {
	PostIncident(context.Context, *domain.Incident) error
	PatchIncident(context.Context, *domain.Incident) error
	GetIncident(context.Context, domain.IncidentID) (*domain.Incident, error)
}

// IncidentAuthorizer authorizes each incident use-case. Its implementation
// applies the role-dependent policy outside of the service orchestration.
type IncidentAuthorizer interface {
	AuthorizePostIncident(context.Context, Identity) error
	AuthorizePatchIncident(context.Context, Identity, uuid.UUID) error
	AuthorizeGetIncident(context.Context, Identity, uuid.UUID) error
}

type IncidentService struct {
	incidents  IncidentRepository
	authorizer IncidentAuthorizer
	logger     *slog.Logger
	domainLog  *slog.Logger
}

func NewIncidentService(incidents IncidentRepository, authorizer IncidentAuthorizer, logger *slog.Logger) IncidentService {
	return IncidentService{
		incidents:  incidents,
		authorizer: authorizer,
		logger:     logger.With("layer", "application"),
		domainLog:  logger.With("layer", "domain"),
	}
}

func (s IncidentService) PostIncident(ctx context.Context, identity Identity, input IncidentInput) (*Incident, error) {
	s.logger.DebugContext(ctx, "posting incident", "subject_id", identity.SubjectID)
	if err := s.authorizer.AuthorizePostIncident(ctx, identity); err != nil {
		s.logger.WarnContext(ctx, "incident post authorization failed", "error", err)
		return nil, err
	}
	incident, err := domain.PostIncident(ctx, s.domainLog, toIncidentPatch(input))
	if err != nil {
		return nil, err
	}
	if err := s.incidents.PostIncident(ctx, incident); err != nil {
		return nil, err
	}
	return fromDomainIncident(incident), nil
}

func (s IncidentService) PatchIncident(ctx context.Context, identity Identity, id uuid.UUID, input IncidentInput) (*Incident, error) {
	s.logger.InfoContext(ctx, "patching incident", "incident_id", id, "subject_id", identity.SubjectID)
	if err := s.authorizer.AuthorizePatchIncident(ctx, identity, id); err != nil {
		s.logger.WarnContext(ctx, "incident patch authorization failed", "incident_id", id, "error", err)
		return nil, err
	}
	incident, err := s.incidents.GetIncident(ctx, domain.IncidentID(id))
	if err != nil {
		return nil, err
	}
	incident, err = domain.PatchIncident(ctx, s.domainLog, toIncidentPatch(input), incident)
	if err != nil {
		return nil, err
	}
	if err := s.incidents.PatchIncident(ctx, incident); err != nil {
		return nil, err
	}
	return fromDomainIncident(incident), nil
}

func (s IncidentService) GetIncident(ctx context.Context, identity Identity, id uuid.UUID) (*Incident, error) {
	s.logger.InfoContext(ctx, "getting incident", "incident_id", id, "subject_id", identity.SubjectID)
	if err := s.authorizer.AuthorizeGetIncident(ctx, identity, id); err != nil {
		s.logger.WarnContext(ctx, "incident get authorization failed", "incident_id", id, "error", err)
		return nil, err
	}
	incident, err := s.incidents.GetIncident(ctx, domain.IncidentID(id))
	if err != nil {
		return nil, err
	}
	incident, err = domain.GetIncident(ctx, s.domainLog, incident)
	if err != nil {
		return nil, err
	}
	return fromDomainIncident(incident), nil
}

func toIncidentPatch(input IncidentInput) domain.IncidentPatch {
	return domain.IncidentPatch{
		Type:        toDomainField(input.Type, func(value int16) domain.IncidentType { return domain.IncidentType(value) }),
		Status:      toDomainField(input.Status, func(value int8) domain.IncidentStatus { return domain.IncidentStatus(value) }),
		Priority:    toDomainField(input.Priority, func(value int8) domain.IncidentPriority { return domain.IncidentPriority(value) }),
		Description: toDomainField(input.Description, func(value string) string { return value }),
		Location: domain.Location{
			Latitude:            toDomainField(input.Location.Latitude, func(value float64) float64 { return value }),
			Longitude:           toDomainField(input.Location.Longitude, func(value float64) float64 { return value }),
			Address:             toDomainField(input.Location.Address, func(value string) string { return value }),
			LocationDescription: toDomainField(input.Location.LocationDescription, func(value string) string { return value }),
		},
	}
}

func toDomainField[Input, Output any](field Field[Input], convert func(Input) Output) domain.Option[Output] {
	result := domain.Option[Output]{Set: field.Set}
	if field.Value != nil {
		value := convert(*field.Value)
		result.Value = &value
	}
	return result
}

func fromDomainIncident(incident *domain.Incident) *Incident {
	return &Incident{
		ID:          uuid.UUID(incident.ID),
		Type:        fromDomainField(incident.Type, func(value domain.IncidentType) int16 { return int16(value) }),
		Status:      fromDomainField(incident.Status, func(value domain.IncidentStatus) int8 { return int8(value) }),
		Priority:    fromDomainField(incident.Priority, func(value domain.IncidentPriority) int8 { return int8(value) }),
		Description: fromDomainField(incident.Description, func(value string) string { return value }),
		Location: LocationInput{
			Latitude:            fromDomainField(incident.Location.Latitude, func(value float64) float64 { return value }),
			Longitude:           fromDomainField(incident.Location.Longitude, func(value float64) float64 { return value }),
			Address:             fromDomainField(incident.Location.Address, func(value string) string { return value }),
			LocationDescription: fromDomainField(incident.Location.LocationDescription, func(value string) string { return value }),
		},
	}
}

func fromDomainField[Domain, Output any](field domain.Option[Domain], convert func(Domain) Output) Field[Output] {
	result := Field[Output]{Set: field.Set}
	if field.Value != nil {
		value := convert(*field.Value)
		result.Value = &value
	}
	return result
}
