package main

import (
	"encoding/json"
	"log"
	"net/http"

	// "os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	db "github.com/nawfal-btw/CE1-Group-3/apps/sql"
)

func (app *app) getAllIncidents(w http.ResponseWriter, r *http.Request) {
	incidentsDB, err := app.queries.GetAllIncidents(r.Context())
	if err != nil {
		log.Println("this bs", err)
	}
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	if err != nil {
		log.Println("ffff", err)
	}
	incidents := make([]incidentCanon, 0, len(incidentsDB))

	for _, incident := range incidentsDB {
		incidents = append(incidents, incidentCanon{
			ID:        incident.ID,
			Latitude:  incident.Latitude,
			Longitude: incident.Longitude,
		})

	}
	if err := encoder.Encode(incidents); err != nil {
		log.Println("fdsafs", err)

	}
}

func (app *app) getIncidentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("this bs", err)
	}
	incidentDB, err := app.queries.GetIncidentByID(r.Context(), int64(id))
	if err != nil {
		log.Println("this bs", err)
	}
	incident_json := incidentCanon{
		ID:        incidentDB.ID,
		Latitude:  incidentDB.Latitude,
		Longitude: incidentDB.Longitude,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(incident_json); err != nil {
		log.Println("this bs", err)
	}

}
func (app *app) postIncident(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var incident db.AddIncidentParams

	if err := decoder.Decode(&incident); err != nil {
		log.Println("post failed")
		return
	}
	app.queries.AddIncident(r.Context(), incident)
}
