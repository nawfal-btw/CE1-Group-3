package application

import (
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
)

// IncidentOptions lists every value a client may choose when filling in an incident.
// Clients use it to build their form, so the allowed values are defined only in the domain.
type IncidentOptions struct {
	Version    int
	Types      []OptionEntry
	Priorities []OptionEntry
}

type OptionEntry struct {
	Code  int
	Label string
}

// GetIncidentOptions returns the domain's catalogues. They are not secret, so no authorization is needed.
func (s IncidentService) GetIncidentOptions() IncidentOptions {
	return IncidentOptions{
		Version:    domain.RulesVersion,
		Types:      toOptionEntries(domain.IncidentTypes()),
		Priorities: toOptionEntries(domain.IncidentPriorities()),
	}
}

func toOptionEntries(entries []domain.CatalogEntry) []OptionEntry {
	result := make([]OptionEntry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, OptionEntry{Code: entry.Code, Label: entry.Label})
	}
	return result
}
