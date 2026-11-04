// This file translates incident HTTP requests and responses. It validates request shape and
// invokes incident use cases, without executing SQL or applying lifecycle rules directly.
package http

import (
	"encoding/json"
	"github.com/google/uuid"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"log"
	"net/http"
)

type IncidentBuffer struct {
	ID          uuid.UUID                  `json:"id"`
	Type        int16                      `json:"type"`
	Status      int8                       `json:"status"`
	Priority    int8                       `json:"priority"`
	Location    application.LocationBuffer `json:"location"`
	Description string                     `json:"description"`
	// Owner        AgencyRef
}
type IncidentHandler struct {
	service application.IncidentService
}

func NewIncidentHandler(s application.IncidentService) *IncidentHandler {
	return &IncidentHandler{service: s}
}

func (h IncidentHandler) postIncident(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var iBuf IncidentBuffer
	err := decoder.Decode(&iBuf)
	if err != nil {
		log.Println(err)
		return
	}
	iPatch := application.PatchIncidentTransport{
		ID:          nil,
		Type:        &iBuf.Type,
		Status:      &iBuf.Status,
		Priority:    &iBuf.Priority,
		Location:    &iBuf.Location,
		Description: &iBuf.Description,
	}
	h.service.CreateIncident(r.Context(), iPatch)
}
func (h IncidentHandler) patchIncident(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var iBuf IncidentBuffer
	err := decoder.Decode(&iBuf)
	if err != nil {
		log.Println(err)
		return
	}
	iPatch := application.PatchIncidentTransport{
		ID:          nil,
		Type:        &iBuf.Type,
		Status:      &iBuf.Status,
		Priority:    &iBuf.Priority,
		Location:    &iBuf.Location,
		Description: &iBuf.Description,
	}
	h.service.CreateIncident(r.Context(), iPatch)
}
func (h IncidentHandler) getIncident(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		log.Println(err)
	}
	h.service.GetIncident(r.Context(), id)
}
