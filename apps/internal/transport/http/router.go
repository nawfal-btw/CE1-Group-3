// This file registers versioned HTTP routes and connects them to handlers. It defines the API
// surface, but delegates request behavior to handlers and application services.
package http

import (
	"net/http"
)

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /incident", getIncidents)
	mux.HandleFunc("GET /incident/{id}" getIncidentByID)
	mux.HandleFunc("POST /incident" postIncident)
	mux.HandleFunc("PATCH /incident/{id}" patchIncidentByID)
	return mux
}
