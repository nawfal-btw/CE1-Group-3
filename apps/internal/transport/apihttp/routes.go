// This file registers versioned HTTP routes and connects them to handlers. It defines the API
// surface, but delegates request behavior to handlers and application services.
package apihttp

import (
	"net/http"
)

type AuxHandler struct {
}
type Handler struct {
	incidentHandler IncidentHandler
	auxHandler      AuxHandler
}

func NewHandler(incidents IncidentHandler) *Handler {
	return &Handler{
		incidentHandler: incidents,
		auxHandler:      AuxHandler{},
	}
}

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/incident/{id}", h.incidentHandler.GetIncident)
	mux.HandleFunc("POST /v1/incident", h.incidentHandler.PostIncident)
	mux.HandleFunc("PATCH /v1/incident/{id}", h.incidentHandler.PatchIncident)

	mux.HandleFunc("GET /v1/checkhealth", h.auxHandler.checkhealth)
	return mux
}

func (h AuxHandler) checkhealth(w http.ResponseWriter, r *http.Request) {

}
