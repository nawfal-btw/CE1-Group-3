package main

import (
	"net/http"
)

func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /incident/{id}", app.getIncidentByID)
	mux.HandleFunc("GET /incident", app.getAllIncidents)
	mux.HandleFunc("POST /incident", app.postIncident)
	return mux
}
