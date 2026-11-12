package apihttp

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
)

// Field records whether a JSON property was present. A JSON null value is
// represented by Set=true and Value=nil.
type Field[T any] struct {
	Value *T
	Set   bool
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
	service       application.IncidentService
	authenticator IncidentAuthenticator
	logger        *slog.Logger
}

// IncidentAuthenticator authenticates HTTP requests. It must not make
// role-dependent authorization decisions.
type IncidentAuthenticator interface {
	AuthenticateIncident(*http.Request) (application.Identity, error)
}

func NewIncidentHandler(service application.IncidentService, authenticator IncidentAuthenticator, logger *slog.Logger) IncidentHandler {
	return IncidentHandler{service: service, authenticator: authenticator, logger: logger.With("layer", "transport")}
}

func (h IncidentHandler) PostIncident(w http.ResponseWriter, r *http.Request) {
	h.logger.DebugContext(r.Context(), "received post incident request")
	identity, ok := h.authenticateIncident(w, r)
	if !ok {

		return
	}
	input, ok := decodeIncidentInput(h.logger, w, r)
	if !ok {
		return
	}
	incident, err := h.service.PostIncident(r.Context(), identity, input)
	if err != nil {
		status := writeError(w, err)
		if status == http.StatusInternalServerError {
			h.logger.ErrorContext(r.Context(), "incident post failed", "error", err, "status", status)
		} else {
			h.logger.WarnContext(r.Context(), "incident post rejected", "error", err, "status", status)
		}
		return
	}

	h.logger.InfoContext(r.Context(), "incident created", "incident_id", incident.ID, "status", http.StatusCreated)
	writeJSON(w, http.StatusCreated, toIncidentResponse(incident))
}

func (h IncidentHandler) PatchIncident(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "received patch incident request")
	identity, ok := h.authenticateIncident(w, r)
	if !ok {
		return
	}
	id, ok := parseIncidentID(w, r)
	if !ok {
		return
	}
	input, ok := decodeIncidentInput(h.logger, w, r)
	if !ok {
		return
	}
	incident, err := h.service.PatchIncident(r.Context(), identity, id, input)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toIncidentResponse(incident))
}

func (h IncidentHandler) GetIncident(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "received get incident request")
	identity, ok := h.authenticateIncident(w, r)
	if !ok {
		return
	}
	id, ok := parseIncidentID(w, r)
	if !ok {
		return
	}
	incident, err := h.service.GetIncident(r.Context(), identity, id)
	if err != nil {
		writeError(w, err)
		return
	}
	h.logger.InfoContext(r.Context(), "incident retrieved", "incident_id", id, "status", http.StatusOK)
	writeJSON(w, http.StatusOK, toIncidentResponse(incident))
}

func (h IncidentHandler) authenticateIncident(w http.ResponseWriter, r *http.Request) (application.Identity, bool) {
	identity, err := h.authenticator.AuthenticateIncident(r)
	if err != nil {
		h.logger.WarnContext(r.Context(), "incident authentication failed", "error", err, "status", http.StatusUnauthorized)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return application.Identity{}, false
	}
	return identity, true
}

func decodeIncidentInput(logger *slog.Logger, w http.ResponseWriter, r *http.Request) (application.IncidentInput, bool) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var request incidentRequest
	if err := decoder.Decode(&request); err != nil {
		status := http.StatusBadRequest
		http.Error(w, err.Error(), status)
		logger.WarnContext(r.Context(), "incident json rejected", "error", err, "status", status)
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

func toApplicationField[T any](field Field[T]) application.Field[T] {
	return application.Field[T]{Value: field.Value, Set: field.Set}
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

func writeError(w http.ResponseWriter, err error) int {
	status := http.StatusInternalServerError
	if errors.Is(err, application.ErrInvalidIncident) {
		status = http.StatusBadRequest
	}
	if errors.Is(err, application.ErrUnauthorized) {
		status = http.StatusForbidden
	}
	http.Error(w, err.Error(), status)
	return status
}
