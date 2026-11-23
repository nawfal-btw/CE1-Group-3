package application

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var ErrInvalidIncident = domain.ErrInvalidIncident
var ErrUnauthorized = errors.New("unauthorized")

type Identity struct {
	SubjectID uuid.UUID
	Roles     []string
}

type Option[T any] struct {
	Value *T
	Set   bool
}

type IncidentInput struct {
	Type        Option[int16]
	Status      Option[int8]
	Priority    Option[int8]
	Location    LocationInput
	Description Option[string]
}

type Incident struct {
	ID          uuid.UUID
	Type        Option[int16]
	Status      Option[int8]
	Priority    Option[int8]
	Location    LocationInput
	Description Option[string]
}

type LocationInput struct {
	Latitude            Option[float64]
	Longitude           Option[float64]
	Address             Option[string]
	LocationDescription Option[string]
}

type IncidentRepository interface {
	PostIncident(context.Context, *domain.Incident) error
	PatchIncident(context.Context, *domain.Incident) error
	GetIncident(context.Context, domain.IncidentID) (*domain.Incident, error)
}

type IncidentAuthorizer interface {
	AuthorizePostIncident(context.Context, Identity) error
	AuthorizePatchIncident(context.Context, Identity, uuid.UUID) error
	AuthorizeGetIncident(context.Context, Identity, uuid.UUID) error
}

type IncidentService struct {
	incidents  IncidentRepository
	authorizer IncidentAuthorizer
	logger     *slog.Logger
	tracer     trace.Tracer
	latency    metric.Float64Histogram
	authorize  metric.Int64Counter
	changes    metric.Int64Counter
}

func NewIncidentService(incidents IncidentRepository, authorizer IncidentAuthorizer, logger *slog.Logger, tracer trace.Tracer, meter metric.Meter) (IncidentService, error) {
	latency, err := meter.Float64Histogram("incident_use_case_duration_seconds", metric.WithUnit("s"))
	if err != nil {
		return IncidentService{}, err
	}
	authorize, err := meter.Int64Counter("incident_authorization_total")
	if err != nil {
		return IncidentService{}, err
	}
	changes, err := meter.Int64Counter("incident_change_total")
	if err != nil {
		return IncidentService{}, err
	}
	return IncidentService{incidents: incidents, authorizer: authorizer, logger: logger, tracer: tracer, latency: latency, authorize: authorize, changes: changes}, nil
}

func (s IncidentService) PostIncident(ctx context.Context, identity Identity, input IncidentInput) (result *Incident, err error) {
	ctx, span := s.tracer.Start(ctx, "application.PostIncident")
	defer func(start time.Time) {
		endSpan(span, err)
		s.recordLatency(ctx, "post_incident", start)
		s.recordChange(ctx, "post_incident", outcomeFromError(err))
	}(time.Now())

	logger := observability.AddPersistentAttributes(ctx, s.logger, "application")
	if authorizationErr := s.authorizer.AuthorizePostIncident(ctx, identity); authorizationErr != nil {
		err = errors.Join(ErrUnauthorized, authorizationErr)
		s.recordAuthorization(ctx, "post_incident", "denied")
		logger.WarnContext(ctx, "incident post authorization denied", "error", authorizationErr)
		return nil, err
	}
	s.recordAuthorization(ctx, "post_incident", "success")

	domainCtx, domainSpan := s.tracer.Start(ctx, "domain.PostIncident")
	incident, domainErr := domain.PostIncident(toIncidentPatch(input))
	endSpan(domainSpan, domainErr)
	if domainErr != nil {
		logger := observability.AddPersistentAttributes(domainCtx, s.logger, "domain")
		logger.WarnContext(domainCtx, "incident validation failed", "error", domainErr)
		return nil, domainErr
	}
	if err = s.incidents.PostIncident(ctx, incident); err != nil {
		logger.ErrorContext(ctx, "incident storage failed", "incident_id", incident.ID, "error", err)
		return nil, err
	}
	logger.InfoContext(ctx, "incident posted", "incident_id", incident.ID)
	return fromDomainIncident(incident), nil
}

func (s IncidentService) PatchIncident(ctx context.Context, identity Identity, id uuid.UUID, input IncidentInput) (result *Incident, err error) {
	ctx, span := s.tracer.Start(ctx, "application.PatchIncident")
	defer func(start time.Time) {
		endSpan(span, err)
		s.recordLatency(ctx, "patch_incident", start)
		s.recordChange(ctx, "patch_incident", outcomeFromError(err))
	}(time.Now())

	logger := observability.AddPersistentAttributes(ctx, s.logger, "application")
	if authorizationErr := s.authorizer.AuthorizePatchIncident(ctx, identity, id); authorizationErr != nil {
		err = errors.Join(ErrUnauthorized, authorizationErr)
		s.recordAuthorization(ctx, "patch_incident", "denied")
		logger.WarnContext(ctx, "incident patch authorization denied", "incident_id", id, "error", authorizationErr)
		return nil, err
	}
	s.recordAuthorization(ctx, "patch_incident", "success")
	incident, err := s.incidents.GetIncident(ctx, domain.IncidentID(id))
	if err != nil {
		logger.WarnContext(ctx, "incident not available for patch", "incident_id", id, "error", err)
		return nil, err
	}
	domainCtx, domainSpan := s.tracer.Start(ctx, "domain.PatchIncident")
	incident, domainErr := domain.PatchIncident(toIncidentPatch(input), incident)
	endSpan(domainSpan, domainErr)
	if domainErr != nil {
		observability.AddPersistentAttributes(domainCtx, s.logger, "domain").WarnContext(domainCtx, "incident validation failed", "incident_id", id, "error", domainErr)
		return nil, domainErr
	}
	if err = s.incidents.PatchIncident(ctx, incident); err != nil {
		logger.ErrorContext(ctx, "incident update failed", "incident_id", id, "error", err)
		return nil, err
	}
	logger.InfoContext(ctx, "incident patched", "incident_id", id)
	return fromDomainIncident(incident), nil
}

func (s IncidentService) GetIncident(ctx context.Context, identity Identity, id uuid.UUID) (result *Incident, err error) {
	ctx, span := s.tracer.Start(ctx, "application.GetIncident")
	defer func(start time.Time) {
		endSpan(span, err)
		s.recordLatency(ctx, "get_incident", start)
	}(time.Now())

	logger := observability.AddPersistentAttributes(ctx, s.logger, "application")
	if authorizationErr := s.authorizer.AuthorizeGetIncident(ctx, identity, id); authorizationErr != nil {
		err = errors.Join(ErrUnauthorized, authorizationErr)
		s.recordAuthorization(ctx, "get_incident", "denied")
		logger.WarnContext(ctx, "incident get authorization denied", "incident_id", id, "error", authorizationErr)
		return nil, err
	}
	s.recordAuthorization(ctx, "get_incident", "success")
	incident, err := s.incidents.GetIncident(ctx, domain.IncidentID(id))
	if err != nil {
		logger.WarnContext(ctx, "incident not found", "incident_id", id, "error", err)
		return nil, err
	}
	domainCtx, domainSpan := s.tracer.Start(ctx, "domain.GetIncident")
	incident, domainErr := domain.GetIncident(incident)
	endSpan(domainSpan, domainErr)
	if domainErr != nil {
		observability.AddPersistentAttributes(domainCtx, s.logger, "domain").ErrorContext(domainCtx, "stored incident failed validation", "incident_id", id, "error", domainErr)
		return nil, domainErr
	}
	return fromDomainIncident(incident), nil
}

func (s IncidentService) recordLatency(ctx context.Context, operation string, start time.Time) {
	s.latency.Record(ctx, time.Since(start).Seconds(), metric.WithAttributes(attribute.String("operation", operation)))
}

func (s IncidentService) recordAuthorization(ctx context.Context, operation, outcome string) {
	s.authorize.Add(ctx, 1, metric.WithAttributes(attribute.String("operation", operation), attribute.String("outcome", outcome)))
}

func (s IncidentService) recordChange(ctx context.Context, operation, outcome string) {
	s.changes.Add(ctx, 1, metric.WithAttributes(attribute.String("operation", operation), attribute.String("outcome", outcome)))
}

func outcomeFromError(err error) string {
	if err == nil {
		return "success"
	}
	if errors.Is(err, ErrUnauthorized) {
		return "denied"
	}
	return "failed"
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

func toIncidentPatch(input IncidentInput) domain.IncidentPatch {
	return domain.IncidentPatch{
		Type:        toDomainOption(input.Type, func(value int16) domain.IncidentType { return domain.IncidentType(value) }),
		Status:      toDomainOption(input.Status, func(value int8) domain.IncidentStatus { return domain.IncidentStatus(value) }),
		Priority:    toDomainOption(input.Priority, func(value int8) domain.IncidentPriority { return domain.IncidentPriority(value) }),
		Description: toDomainOption(input.Description, func(value string) string { return value }),
		Location: domain.Location{
			Latitude:            toDomainOption(input.Location.Latitude, func(value float64) float64 { return value }),
			Longitude:           toDomainOption(input.Location.Longitude, func(value float64) float64 { return value }),
			Address:             toDomainOption(input.Location.Address, func(value string) string { return value }),
			LocationDescription: toDomainOption(input.Location.LocationDescription, func(value string) string { return value }),
		},
	}
}

func toDomainOption[Input, Output any](option Option[Input], convert func(Input) Output) domain.Option[Output] {
	result := domain.Option[Output]{Set: option.Set}
	if option.Value != nil {
		value := convert(*option.Value)
		result.Value = &value
	}
	return result
}

func fromDomainIncident(incident *domain.Incident) *Incident {
	return &Incident{
		ID:          uuid.UUID(incident.ID),
		Type:        fromDomainOption(incident.Type, func(value domain.IncidentType) int16 { return int16(value) }),
		Status:      fromDomainOption(incident.Status, func(value domain.IncidentStatus) int8 { return int8(value) }),
		Priority:    fromDomainOption(incident.Priority, func(value domain.IncidentPriority) int8 { return int8(value) }),
		Description: fromDomainOption(incident.Description, func(value string) string { return value }),
		Location: LocationInput{
			Latitude:            fromDomainOption(incident.Location.Latitude, func(value float64) float64 { return value }),
			Longitude:           fromDomainOption(incident.Location.Longitude, func(value float64) float64 { return value }),
			Address:             fromDomainOption(incident.Location.Address, func(value string) string { return value }),
			LocationDescription: fromDomainOption(incident.Location.LocationDescription, func(value string) string { return value }),
		},
	}
}

func fromDomainOption[Domain, Output any](option domain.Option[Domain], convert func(Domain) Output) Option[Output] {
	result := Option[Output]{Set: option.Set}
	if option.Value != nil {
		value := convert(*option.Value)
		result.Value = &value
	}
	return result
}
