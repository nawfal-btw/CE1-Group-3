package http

import (
	"net/http"

	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
)

// incidentOptionsResponse is the JSON shape of GET /v1/incident/options, e.g.
// {"version":1,"types":[{"code":1,"label":"Type 1"}],"priorities":[...]}
type incidentOptionsResponse struct {
	Version    int                   `json:"version"`
	Types      []optionEntryResponse `json:"types"`
	Priorities []optionEntryResponse `json:"priorities"`
}

type optionEntryResponse struct {
	Code  int    `json:"code"`
	Label string `json:"label"`
}

// GetIncidentOptions returns the values the client may offer in the incident form.
func (h IncidentHandler) GetIncidentOptions(w http.ResponseWriter, r *http.Request) {
	options := h.incidentService.GetIncidentOptions()
	writeJSON(w, http.StatusOK, incidentOptionsResponse{
		Version:    options.Version,
		Types:      toOptionEntryResponses(options.Types),
		Priorities: toOptionEntryResponses(options.Priorities),
	})
}

func toOptionEntryResponses(entries []application.OptionEntry) []optionEntryResponse {
	result := make([]optionEntryResponse, 0, len(entries))
	for _, entry := range entries {
		result = append(result, optionEntryResponse{Code: entry.Code, Label: entry.Label})
	}
	return result
}
