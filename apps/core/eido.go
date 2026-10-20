package main

import "time"

type (
	IncidentStatus int16
	IncidentType   int16
)
type GIS struct {
	Latitude  float64
	Longitude float64
}

const (
	IncOpen IncidentStatus = iota
	IncInProgress
	IncClosed
	IncCancelled
)

type (
	ID       int64
	Code     int16
	Priority int8
)

type IncidentInternal struct {
	AgencyRef                    AgencyComponent
	InternalIncidentID           string
	InternalIncidentTypeCode     string
	InternalIncidentTypeText     string
	InternalIncidentStatusCode   string
	InternalIncidentPriorityText string
}
type IncidentCanon struct {
	// IncidentTypeInternalCode int16
	// IncidentTypeInternalText string
	CanonIncidentID         int64
	CanonIncidentTypeCode   int16
	CanonIncidentTypeText   string
	CanonIncidentStatusCode int16
	CanonIncidentPriority   int16
	IncidentInternal        []IncidentInternal
	IncidentOwner           *AgencyComponent
	CallRefs                []*CallComponent
	DispatcherRefs          []*DispatcherComponent
	AllocationRef           []*Allocation
	// DocumentID
}

// Unsure, all individuals that call?j
type CallComponent struct {
	// CallQueueID
	CallTypePrimary
	CallDirection
	CallStartTime
	CallAnswerTime
	CallStatus
	AgentRef
	CallBackRef
	LocationRef
	PersonRef
}
type CallBackComponent struct {
	CallBackInformationURI string
}
type DispatcherComponent struct {
	AgencyRef
	AgentRef
	ResourceRefs
}
type ResourceComponent struct {
	ResourceTypeCanon
	ResourceName
	ResourceStatus
	ResourceTypeInternal
	ResourceAttribute
	AgentRef
	AgencyRef
	LocationRef
	DestinationRef
	ETA
}

// Required
type AgencyComponent struct {
	AgencyContactInfo
	AgencyRelationType
	AgencyType
}

// Conditional
type AgentComponent struct {
	AgentContactInfo
	AgentRoleCanon
	AgencyRef
}
type LocationComponent struct {
	LocationGIS
	CrossStreet
	IntersectingStreet
	LocationDescription
	CellTowerSectorID
}
type PersonComponent struct {
	PersonRelationType
	PersonRelationTime
	CallRef
	LocationRef
	CallBackRef
}
type VehicleComponent struct {
	VehicleRelationType
	VehicleRelationTime
	LocationRef
}

// Link component maybe
