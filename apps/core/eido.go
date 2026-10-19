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

type IncidentComponent struct {
	// IncidentTypeInternalCode int16
	// IncidentTypeInternalText string
	IncidentTypeCommonRegistry   int16
	IncidentStatusCommonRegistry int16
	AgentReference               int32
	LocationReference            LocationComponent
	PersonReference              PersonComponent
	VehicleReference             VehicleComponent
	DocumentIdentification       int32
	ReportNumberType             int16
	IncidentCommonPriorityNumber int8
}

// Unsure, all individuals that call?j
type CallComponent struct {
	StandardPrimaryCallType int16
	Direction               int8
	CallStartTime           time.Time
	AnswerStartTime         time.Time
	CallStateRegistry       int8
	AgentReference          AgentComponent
	CallbackReference       CallbackComponent
	// Location of call-site
	LocationReference LocationComponent
	PersonReference   PersonComponent
}
type CallbackComponent struct {
	CallbackInformationURI string
}
type DispatchComponent struct {
	AgencyRef            AgencyComponent
	AgentRef             AgentComponent
	EmergencyResourceRef EmergencyResourceComponent
}
type EmergencyResourceComponent struct {
	EmergencyResourceTypeCommonRegistry int8
	ResourceAttributeRegistry           int8
	EmergencyResourceName               string
	PrimaryUnitStatusRegistry
	SecondaryUnitStatusRegistry
	AgentRef        AgentComponent
	PersonRef       PersonComponent
	VehicleRef      VehicleComponent
	UnitLocationRef LocationComponent
	ETA             time.Time
}

// Required
type AgencyComponent struct {
	AgencyContactInfo             string
	AgencyRoleDescriptionRegistry string
	AgencyType                    string
	IncidentOwner                 bool
}

// Conditional
type AgentComponent struct {
	AgentContactInfo  string
	AgentRoleRegistry int
	AgentAgencyRef    int
}
type LocationComponent struct {
	LocationTypeDescriptionRegistry int8
	LocationByValue                 GIS
	LocationDescriptionText         string
	CrossStreetByValue              string
	IntersectingStreetByValue       string
	CellTowerSectorID               int32
}
type PersonComponent struct {
	PersonIncidentRoleRegistry int8
	CallIdentifier             CallComponent
	LocationRef                LocationComponent
	CallbackRef                CallbackComponent
}
type VehicleComponent struct {
	VehicleRelationshipType int8
	LocationRef             LocationComponent
}

// Link component maybe
