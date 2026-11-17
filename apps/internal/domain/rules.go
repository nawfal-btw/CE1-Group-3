package domain

// This file is the single source of truth for the values an incident may take.
// For now the catalogues are written in Go. Later they can be loaded from a versioned
// rules file without changing the functions below or anything that calls them.

// CatalogEntry is one allowed value: the code stored in the database and a readable label.
type CatalogEntry struct {
	Code  int
	Label string
}

// RulesVersion identifies this set of catalogues. Increase it whenever a catalogue changes,
// so clients know their cached copy is outdated.
const RulesVersion = 1

// Test values: they only exist to exercise the architecture end to end.
var incidentTypes = []CatalogEntry{
	{Code: 1, Label: "Type 1"},
	{Code: 2, Label: "Type 2"},
	{Code: 3, Label: "Type 3"},
}

var incidentPriorities = []CatalogEntry{
	{Code: 1, Label: "Priority 1"},
	{Code: 2, Label: "Priority 2"},
	{Code: 3, Label: "Priority 3"},
}

// IncidentTypes returns a copy of the allowed incident types, so callers cannot modify the catalogue.
func IncidentTypes() []CatalogEntry {
	return append([]CatalogEntry(nil), incidentTypes...)
}

// IncidentPriorities returns a copy of the allowed incident priorities.
func IncidentPriorities() []CatalogEntry {
	return append([]CatalogEntry(nil), incidentPriorities...)
}

func IsValidType(t IncidentType) bool {
	return inCatalog(incidentTypes, int(t))
}

func IsValidPriority(p IncidentPriority) bool {
	return inCatalog(incidentPriorities, int(p))
}

// inCatalog reports whether code is one of the entries' codes.
func inCatalog(entries []CatalogEntry, code int) bool {
	for _, entry := range entries {
		if entry.Code == code {
			return true
		}
	}
	return false
}
