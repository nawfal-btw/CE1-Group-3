package main

import (
	"encoding/json"
	_ "github.com/jackc/pgx/v5/stdlib"
	core "github.com/nawfal-btw/CE1-Group-3/apps/core"
	"log"
	"net/http"
	"strconv"
)

func (app *app) getIncidentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Println("this bs", err)
	}

	incident, err := app.queries.GetIncidentByID(r.Context(), int64(id))
	if err != nil {
		log.Println("this bs", err)
	}
	incident_json := core.Incident{
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
		Type: core.Field[core.IncidentType]{Value: core.IncidentType(incident.Type), Initialized: true},
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
		ID:   core.Field[core.IncidentID]{Value: core.IncidentID(incident.ID), Initialized: true},
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(incident_json); err != nil {
		log.Println("this bs", err)
	}

}
