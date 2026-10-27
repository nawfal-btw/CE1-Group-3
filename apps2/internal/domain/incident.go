// This file defines the incident entity and its permitted lifecycle operations. It contains
// business invariants, but no HTTP, database, or serialization concerns.
package domain

import (
	"time"

	app "github.com/nawfal-btw/CE1-Group-3/apps2/internal/application"
)

type (
	IncidentID       int64
	IncidentType     int16
	IncidentPriority int8
	IncidentStatus   int8
)

const (
	PLANNING    IncidentStatus = 0
	IN_PROGRESS IncidentStatus = 1
	RESOLVED    IncidentStatus = 2
	CANCELLED   IncidentStatus = 3
)
const (
	_NONEMERG IncidentType = 0
	_911H     IncidentType = 1
	_911OPN   IncidentType = 2
	_ABDOM    IncidentType = 3
	_ABDUCTF  IncidentType = 4
	_ABDUCTNF IncidentType = 5
	_ABUSE    IncidentType = 6
	_ACCIND   IncidentType = 7
	_ACTSHOOT IncidentType = 8
	_ADMIN    IncidentType = 9
	_AIR      IncidentType = 10
	_AIRINC   IncidentType = 11
	_ALIEN    IncidentType = 12
	_ALLERGIC IncidentType = 13
	_ALMA     IncidentType = 14
	_ALMB     IncidentType = 15
	_ALMH     IncidentType = 16
)

type Field[T any] struct {
	Value       T
	Initialized bool
}
type Incident struct {
	ID          Field[IncidentID]
	Type        Field[IncidentType]
	Status      Field[IncidentStatus]
	Priority    Field[IncidentPriority]
	Location    Field[LocationComponent]
	Description Field[string]

	CreatedAtUTC Field[time.Time]
	UpdatedAtUTC Field[time.Time]
	ClosedAtUTC  Field[time.Time]

	// Owner        AgencyRef
	// Calls        []CallComponent
	// Allocations  []Allocation
	// ExternalRefs []LocalIncident
}

type LocationGIS struct {
	Latitude  float64
	Longitude float64
}
type LocationComponent struct {
	GIS                 Field[LocationGIS]
	Address             string
	ProximateLandmarks  []string
	LocationDescription string
}

func (i Incident) NewIncident(iID IncidentID, iType IncidentType, iStatus IncidentStatus, iPrio IncidentPriority, iLoc LocationComponent) {
	incident := Incident{
		ID,
	}

}
