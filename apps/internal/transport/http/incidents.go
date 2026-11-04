// This file translates incident HTTP requests and responses. It validates request shape and
// invokes incident use cases, without executing SQL or applying lifecycle rules directly.
package http

import (
	"encoding/json"
	"github.com/google/uuid"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"log"
	"net/http"
	"strconv"
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

func getIncidentByID(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")
	_, err := strconv.Atoi(idPath)
	if err != nil {
		log.Println(err)
	}
	// build application req struct
	// call application method?
}
func (iHandler IncidentHandler) postIncident(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var iBuf IncidentBuffer
	err := decoder.Decode(&iBuf)
	if err != nil {
		log.Println(err)
		return
	}
	iPatch := application.PatchIncidentCommand{
		Type:     &iBuf.Type,
		Status:   iBuf.Status,
		Priority: iBuf.Priority,
	}
	iHandler.service.CreateIncident()

}
