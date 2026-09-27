package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	db "github.com/nawfal-btw/CE1-Group-3/apps/sql"
)

func (app *app) getIncidentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		os.Exit(1)
	}
	incidentDB, err := app.queries.GetIncidentByID(r.Context(), int64(id))
	if err != nil {
		os.Exit(1)
	}
	incident_json := incidentCanon{
		ID:        incidentDB.ID,
		Latitude:  incidentDB.Latitude,
		Longitude: incidentDB.Longitude,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(incident_json); err != nil {
		os.Exit(1)
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
