package main

import "errors"

type IncidentType int16

type IncidentID int64
type CallID int64
type AllocationID int64
type AgencyID int64
type AgentID int64
type ResourceID int64

type IncidentTypeCode int16
type IncidentStatusCode int16
type IncidentPriority int8

type IncidentStatus int8

const (
	IncidentStatusNew IncidentStatus = iota
	IncidentStatusActive
	IncidentStatusSuspended
	IncidentStatusClosed
)

type AllocationStatus int8

const (
	AllocationStatusAssigned AllocationStatus = iota
	AllocationStatusEnRoute
	AllocationStatusArrived
	AllocationStatusReleased
)

type CallType int8
type CallDirection int8
type CallStatus int8

var ErrInvalidIncidentTransition = errors.New("invalid incident status transition")

type AgencyRef struct{ AgencyID AgencyID }
type AgentRef struct{ AgentID AgentID }
type ResourceRefs struct{ ResourceIDs []ResourceID }
type LocationRef struct{ LocationID string }
type CallRef struct{ CallID CallID }
type CallBackRef struct{ CallBackURI string }

type LocationGIS struct {
	Latitude  float64
	Longitude float64
}

type CrossStreet struct{ Name string }
type IntersectingStreet struct{ Name string }
type LocationDescription struct{ Text string }
type CellTowerSectorID struct{ Value string }

type AgencyContactInfo struct{ ContactURI string }
type AgencyRelationType struct{ Value string }
type AgencyType struct{ Value string }
type AgentContactInfo struct{ ContactURI string }
type AgentRoleCanon struct{ Value string }
type PersonRelationType struct{ Value string }
type PersonRelationTime struct{ At string }
type VehicleRelationType struct{ Value string }
type VehicleRelationTime struct{ At string }

type ResourceTypeCanon struct{ Value string }
type ResourceName struct{ Value string }
type ResourceStatus struct{ Value string }
type ResourceTypeInternal struct{ Value string }
type ResourceAttribute struct{ Value string }
type DestinationRef struct{ LocationID string }
type ETA struct{ At string }
