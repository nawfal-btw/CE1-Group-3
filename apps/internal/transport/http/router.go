// This file registers versioned HTTP routes and connects them to handlers. It defines the API
// surface, but delegates request behavior to handlers and application services.
package http

import (
	"net/http"
)

type Handler struct {
	incidentHandler IncidentHandler
}

func (h Handler) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /incident/{id}", h.incidentHandler.getIncident)
	mux.HandleFunc("POST /incident", h.incidentHandler.postIncident)
	mux.HandleFunc("PATCH /incident/{id}", h.incidentHandler.patchIncident)
	return mux
}
