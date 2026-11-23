package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/observability"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Field records whether a JSON property was present. A JSON null value is
// represented by Set=true and Value=nil.
type Field[T any] struct {
	Value *T
	Set   bool
}

type incidentRequest struct {
	Type        Field[int16]    `json:"type"`
	Status      Field[int8]     `json:"status"`
	Priority    Field[int8]     `json:"priority"`
	Description Field[string]   `json:"description"`
	Location    locationRequest `json:"location"`
}

type locationRequest struct {
	Latitude            Field[float64] `json:"latitude"`
	Longitude           Field[float64] `json:"longitude"`
	Address             Field[string]  `json:"address"`
	LocationDescription Field[string]  `json:"locationDescription"`
}

type incidentResponse struct {
	ID          uuid.UUID        `json:"id"`
	Type        *int16           `json:"type"`
	Status      *int8            `json:"status"`
	Priority    *int8            `json:"priority"`
	Location    locationResponse `json:"location"`
	Description *string          `json:"description"`
}

type locationResponse struct {
	Latitude            *float64 `json:"latitude"`
	Longitude           *float64 `json:"longitude"`
	Address             *string  `json:"address"`
	LocationDescription *string  `json:"locationDescription"`
}

type IncidentHandler struct {
	incidentService application.IncidentService
	authorizer      IncidentAuthorizer
	logger          *slog.Logger
	tracer          trace.Tracer
}

func (field *Field[T]) UnmarshalJSON(data []byte) error {
	field.Set = true
	if string(data) == "null" {
		field.Value = nil
		return nil
	}
	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	field.Value = &value
	return nil
}

func NewIncidentHandler(service application.IncidentService, authenticator IncidentAuthorizer, logger *slog.Logger, tracer trace.Tracer) IncidentHandler {
	return IncidentHandler{incidentService: service, authorizer: authenticator, logger: logger, tracer: tracer}
}

func (h IncidentHandler) PostIncident(w http.ResponseWriter, r *http.Request) {
	ctx, span, logger := h.startRequest(r, "transport.PostIncident")
	var err error
	defer func() { endSpan(span, err) }()
	r = r.WithContext(ctx)
	identity, ok := h.authenticateIncident(w, r)
	if !ok {
		return
	}
	// ctx = observability.WithUserID(ctx, identity.SubjectID.String())
	r = r.WithContext(ctx)
	logger = observability.AddPersistentAttributes(ctx, h.logger, "transport")
	input, ok := decodeIncidentInput(w, r)
	if !ok {
		logger.WarnContext(ctx, "incident request body rejected")
		return
	}
	incident, err := h.incidentService.PostIncident(ctx, identity, input)
	if err != nil {
		logger.WarnContext(ctx, "incident post request failed", "error", err)
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toIncidentResponse(incident))
}

func (h IncidentHandler) PatchIncident(w http.ResponseWriter, r *http.Request) {
	ctx, span, logger := h.startRequest(r, "transport.PatchIncident")
	var err error
	defer func() { endSpan(span, err) }()
	r = r.WithContext(ctx)
	identity, ok := h.authenticateIncident(w, r)
	if !ok {
		return
	}
	// ctx = observability.WithUserID(ctx, identity.SubjectID.String())
	r = r.WithContext(ctx)
	logger = observability.AddPersistentAttributes(ctx, h.logger, "transport")
	id, ok := parseIncidentID(w, r)
	if !ok {
		logger.WarnContext(ctx, "incident patch request has invalid ID")
		return
	}
	input, ok := decodeIncidentInput(w, r)
	if !ok {
		logger.WarnContext(ctx, "incident patch request body rejected", "incident_id", id)
		return
	}
	incident, err := h.incidentService.PatchIncident(ctx, identity, id, input)
	if err != nil {
		logger.WarnContext(ctx, "incident patch request failed", "incident_id", id, "error", err)
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toIncidentResponse(incident))
}

func (h IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	ctx, span, logger := h.startRequest(r, "transport.GetIncident")
	var err error
	defer func() { endSpan(span, err) }()
	r = r.WithContext(ctx)
	identity, ok := h.authenticateIncident(w, r)
	if !ok {
		return
	}
	// ctx = observability.WithUserID(ctx, identity.SubjectID.String())
	r = r.WithContext(ctx)
	logger = observability.AddPersistentAttributes(ctx, h.logger, "transport")
	id, ok := parseIncidentID(w, r)
	if !ok {
		logger.WarnContext(ctx, "incident get request has invalid ID")
		return
	}
	incident, err := h.incidentService.GetIncident(ctx, identity, id)
	if err != nil {
		logger.WarnContext(ctx, "incident get request failed", "incident_id", id, "error", err)
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toIncidentResponse(incident))
}

func (h IncidentHandler) authenticateIncident(w http.ResponseWriter, r *http.Request) (application.Identity, bool) {
	identity, err := h.authorizer.AuthenticateIncident(r)
	if err != nil {
		observability.AddPersistentAttributes(r.Context(), h.logger, "transport").WarnContext(r.Context(), "incident authentication failed", "error", err)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return application.Identity{}, false
	}
	return identity, true
}

type contextKey struct{}

func (h IncidentHandler) startRequest(r *http.Request, spanName string) (context.Context, trace.Span, *slog.Logger) {
	ctx := r.Context()
	ctx, span := h.tracer.Start(ctx, spanName)

	metadata := observability.GetMetadataFromContext(ctx)
	metadata.Endpoint = r.Pattern
	metadata.TraceID = trace.SpanContextFromContext(ctx).TraceID().String()

	ctx = observability.AddMetadataToContext(ctx, metadata)

	return ctx, span, observability.AddPersistentAttributes(ctx, h.logger, "transport")
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}

func decodeIncidentInput(w http.ResponseWriter, r *http.Request) (application.IncidentInput, bool) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request incidentRequest
	if err := decoder.Decode(&request); err != nil {
		writeError(w, err)
		return application.IncidentInput{}, false
	}
	return toIncidentInput(request), true
}

func parseIncidentID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, err)
		return uuid.Nil, false
	}
	return id, true
}

func toIncidentInput(request incidentRequest) application.IncidentInput {
	return application.IncidentInput{
		Type:        toApplicationField(request.Type),
		Status:      toApplicationField(request.Status),
		Priority:    toApplicationField(request.Priority),
		Description: toApplicationField(request.Description),
		Location: application.LocationInput{
			Latitude:            toApplicationField(request.Location.Latitude),
			Longitude:           toApplicationField(request.Location.Longitude),
			Address:             toApplicationField(request.Location.Address),
			LocationDescription: toApplicationField(request.Location.LocationDescription),
		},
	}
}

func toApplicationField[T any](field Field[T]) application.Option[T] {
	return application.Option[T]{Value: field.Value, Set: field.Set}
}

func toIncidentResponse(incident *application.Incident) incidentResponse {
	return incidentResponse{
		ID:          incident.ID,
		Type:        incident.Type.Value,
		Status:      incident.Status.Value,
		Priority:    incident.Priority.Value,
		Description: incident.Description.Value,
		Location: locationResponse{
			Latitude:            incident.Location.Latitude.Value,
			Longitude:           incident.Location.Longitude.Value,
			Address:             incident.Location.Address.Value,
			LocationDescription: incident.Location.LocationDescription.Value,
		},
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if errors.Is(err, application.ErrInvalidIncident) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, application.ErrUnauthorized) {
		status = http.StatusForbidden
	}
	http.Error(w, err.Error(), status)
}
