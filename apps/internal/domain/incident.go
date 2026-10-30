// This file defines the incident entity and its permitted lifecycle operations. It contains
// business invariants, but no HTTP, database, or serialization concerns.
package domain

import (
	"errors"
	"time"
)

type (
	IncidentID       int64
	IncidentType     int16
	IncidentPriority int8
	IncidentStatus   int8
)

// This should not be hardcoded but derived from a file
const (
	PRIORITYMIN IncidentPriority = iota
	_1
	_2
	_3
	_4
	_5
	PRIORITYMAX
)

// This should not be hardcoded but derived from a file
const (
	STATUSMIN IncidentStatus = iota
	PLANNING
	IN_PROGRESS
	RESOLVED
	CANCELLED
	STATUSMAX
)

// This should not be hardcoded but derived from a file
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
	Value T
	Init  bool
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
	Version string
}
type IncidentPatch struct {
	Type        Field[IncidentType]
	Status      Field[IncidentStatus]
	Priority    Field[IncidentPriority]
	Location    LocationComponent
	Description string

	// Owner        AgencyRef
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

func NewIncident(p IncidentPatch) (*Incident, error) {
	var i Incident
	// is it init
	if p.Type.Init {
		// is it valid -> Y: Valid
		if !i.IsValidType(p.Type.Value) {
			return &Incident{}, errors.New("Invalid Type")
		}
	}
	if p.Status.Init {
		// is it valid -> Y: Valid
		if !i.IsValidStatus(p.Status.Value) {
			return &Incident{}, errors.New("Invalid Status")
		}
	}
	if p.Priority.Init {
		// is it valid -> Y: Valid
		if !i.IsValidPriority(p.Priority.Value) {
			return &Incident{}, errors.New("Invalid Priority")
		}
	}
	if p.Location.GIS.Init {
		// is it valid -> Y: Valid
		if !i.IsValidLocationGIS(p.Location) {
			return &Incident{}, errors.New("Invalid Location")
		}
	}
	i.Type = p.Type
	i.Status = p.Status
	i.Priority = p.Priority
	i.Location = p.Location
	i.Description = p.Description
	return &i, nil
}
func (i *Incident) PatchIncident(p IncidentPatch) error {
	if p.Type.Init {
		if !i.IsValidType(p.Type.Value) {
			return errors.New("Invalid Type")
		}
	}
	if p.Status.Init {
		if !i.IsValidStatus(p.Status.Value) {
			return errors.New("Invalid Status")
		}
		if isValid, err := i.IsValidTransition(p.Status); !isValid {
			return err
		}
	}
	if p.Priority.Init {
		if !i.IsValidPriority(p.Priority.Value) {
			return errors.New("Invalid Priority")
		}
	}
	if p.Location.GIS.Init {
		if !i.IsValidLocationGIS(p.Location) {
			return errors.New("Invalid Location")
		}
	}
	i.Type = p.Type
	i.Status = p.Status
	i.Priority = p.Priority
	i.Location = p.Location
	i.Description = p.Description
	return nil
}
func (i Incident) IsValidTransition(p any) (bool, error) {
	switch t := p.(type) {
	case IncidentStatus:
		isValid, err := i.IsValidStatusTransition(t)
		return isValid, err
	default:
		return false, errors.New("Unsupported transition check")
	}

}
func (i Incident) IsValidStatusTransition(s IncidentStatus) (bool, error) {
	switch i.Status.Value {
	case PLANNING:
		if s == RESOLVED {
			return false, errors.New("Invalid Status Transition PLANNING -> RESOLVED")
		}
	case IN_PROGRESS:
		if s == PLANNING {
			return false, errors.New("Invalid Status Transition IN PROGRESS -> PLANNING")
		}
	case CANCELLED:
		if s == RESOLVED {
			return false, errors.New("Invalid Status Transition IN CANCELLED -> RESOLVED")
		}
	case RESOLVED:
		if s == CANCELLED {
			return false, errors.New("Invalid Status Transition IN RESOLVED -> CANCELLED")
		}
	default:
		return false, errors.New("Invalid Status slipped through...")
	}
	return true, nil

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
