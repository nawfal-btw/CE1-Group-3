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
	PRIORITYMIN IncidentPriority = iota
	_1
	_2
	_3
	_4
	_5
	PRIORITYMAX
)
const (
	STATUSMIN IncidentStatus = iota
	PLANNING
	IN_PROGRESS
	RESOLVED
	CANCELLED
	STATUSMAX
)
const (
	TYPEMIN IncidentType = iota
	_NONEMERG
	_911H
	_911OPN
	_ABDOM
	_ABDUCTF
	_ABDUCTNF
	_ABUSE
	_ACCIND
	_ACTSHOOT
	_ADMIN
	_AIR
	_AIRINC
	_ALIEN
	_ALLERGIC
	_ALMA
	_ALMB
	_ALMH
	TYPEMAX
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
	Location    LocationComponent
	Description string

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

func NewIncident(iType IncidentType, iStatus IncidentStatus, iPriority IncidentPriority, iLocation LocationComponent, iDescription string) (*Incident, error) {
	var incident Incident
	if incident.IsValidType(iType) {
		incident.Type.Initialized = true
		incident.Type.Value = iType
	} else {
		incident.Type.Initialized = false
	}
	if incident.IsValidStatus(iStatus) {
		incident.Status.Initialized = true
		incident.Type.Value = iType
	} else {
		incident.Status.Initialized = false
	}
	if incident.IsValidPriority(iPriority) {
		incident.Priority.Initialized = true
		incident.Priority.Value = iPriority
	} else {
		incident.Priority.Initialized = false
	}
	if incident.IsValidLocationGIS(iLocation) {
		incident.Location.GIS.Initialized = true
		incident.Location.GIS.Value.Longitude = iLocation.GIS.Value.Longitude
		incident.Location.GIS.Value.Latitude = iLocation.GIS.Value.Latitude
	} else {
		incident.Priority.Initialized = false
	}
	return &incident, nil
}

func (i Incident) IsValidType(iType IncidentType) bool {
	if iType <= TYPEMIN || iType >= TYPEMAX {
		return false
	}
	return true
}
func (i Incident) IsValidStatus(iStatus IncidentStatus) bool {
	if iStatus <= STATUSMIN || iStatus >= STATUSMAX {
		return false
	}
	return true
}
func (i Incident) IsValidPriority(iPriority IncidentPriority) bool {
	if iPriority <= PRIORITYMIN || iPriority >= PRIORITYMAX {
		return false
	}
	return true
}

func (i Incident) IsValidLocationGIS(iLocation LocationComponent) bool {
	if iLocation.GIS.Value.Latitude < 0.0 || iLocation.GIS.Value.Longitude < 0.0 {
		return false
	}
	if iLocation.GIS.Value.Latitude > 180.0 || iLocation.GIS.Value.Longitude > 180.0 {
		return false
	}
	return true
}
