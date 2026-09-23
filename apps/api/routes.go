package main

import (
	"net/http"
)

func routes() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /incident/{id}", getIncidentByID)
	mux.HandleFunc("GET /incident", getAllIncidents)
	mux.HandleFunc("POST /incident", postIncident)
}
