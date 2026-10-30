// This file translates incident HTTP requests and responses. It validates request shape and
// invokes incident use cases, without executing SQL or applying lifecycle rules directly.
package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	app "github.com/nawfal-btw/CE1-Group-3/apps2/internal/application"
)

func getIncidentByID(w http.ResponseWriter, r *http.Request){
	idPath := r.PathValue("id")
	id, err := strconv.Atoi(idPath)
	if err != nil {
		log.Println(err)
	}
	// build application req struct
	// call application method?
	application.

}
func getIncidents(w http.ResponseWriter, r *http.Request){
	application.
}
func postIncident(w http.ResponseWriter, r *http.Request){
	decoder := json.NewDecoder(r.Body)
	app.IncidentService.

	err := decoder.Decode(&incident)
	if err != nil {
		log.Println(err)
		return
	}
	application.
	

}
