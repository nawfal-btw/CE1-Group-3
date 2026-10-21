package main

import "time"

type Incident struct {
	ID          IncidentID
	Type        IncidentType
	Status      IncidentStatus
	Priority    IncidentPriority
	Location    LocationComponent
	Description string

	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  *time.Time

	Owner         AgencyRef
	Calls         []CallComponent
	Allocations   []Allocation
	ExternalRefs  []ExternalIncident
	StatusHistory []StatusChange
}
type ExternalIncident struct {
	AgencyID   AgencyID
	SystemName string
	IncidentID string
	TypeCode   string
	TypeText   string
	StatusCode string
}

type IncidentInternal struct {
	AgencyRef                    AgencyComponent
	InternalIncidentID           string
	InternalIncidentTypeCode     string
	InternalIncidentTypeText     string
	InternalIncidentStatusCode   string
	InternalIncidentPriorityText string
}

type StatusChange struct {
	From      IncidentStatus
	To        IncidentStatus
	ChangedAt time.Time
	ChangedBy AgentID
	Reason    string
}

type CallComponent struct {
	ID             CallID
	Type           CallType
	Direction      CallDirection
	CallStartTime  time.Time
	CallAnswerTime *time.Time
	CallEndTime    *time.Time
	Status         CallStatus
	CallBack       *CallBackComponent
	Location       *LocationComponent
	Person         *PersonComponent
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

type Allocation struct {
	ID           AllocationID
	ResourceID   ResourceID
	AgencyID     AgencyID
	DispatcherID AgentID
	Status       AllocationStatus
	AssignedAt   time.Time
	ArrivedAt    *time.Time
	ReleasedAt   *time.Time
	Destination  *LocationComponent
	ETA          *time.Time
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
