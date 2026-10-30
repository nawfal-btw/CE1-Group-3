package core

import "time"

type (
	IncidentID         int64
	IncidentType       int16
	IncidentTypeCode   int16
	IncidentStatusCode int16
	IncidentPriority   int8
	IncidentStatus     int8
)

const (
	PLANNING    IncidentStatus = 0
	IN_PROGRESS IncidentStatus = 1
	RESOLVED    IncidentStatus = 2
	CANCELLED   IncidentStatus = 3
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

// type AgencyID int64
// type LocalIncident struct {
// 	AgencyID   AgencyID
// 	SystemName string
// 	IncidentID string
// 	TypeCode   string
// 	TypeText   string
// 	StatusCode string
// }

// type (
// 	CallID        int64
// 	CallType      int8
// 	CallDirection int8
// 	CallStatus    int8
// )
// type CallComponent struct {
// 	ID             CallID
// 	Type           CallType
// 	Direction      CallDirection
// 	CallStartTime  time.Time
// 	CallAnswerTime *time.Time
// 	CallEndTime    *time.Time
// 	Status         CallStatus
// 	CallBack       string
// 	Location       *LocationComponent
// 	Person         *PersonComponent
// }

// type ResourceComponent struct {
// 	ResourceTypeCanon
// 	ResourceName
// 	ResourceStatus
// 	ResourceTypeInternal
// 	ResourceAttribute
// 	AgentRef
// 	AgencyRef
// 	LocationRef
// 	DestinationRef
// 	ETA
// }
//
// type Allocation struct {
// 	ID           AllocationID
// 	ResourceID   ResourceID
// 	AgencyID     AgencyID
// 	DispatcherID AgentID
// 	Status       AllocationStatus
// 	AssignedAt   time.Time
// 	ArrivedAt    *time.Time
// 	ReleasedAt   *time.Time
// 	Destination  *LocationComponent
// 	ETA          *time.Time
// }
//
// // Required
// type AgencyContactInfo struct {
// 	establishmentName string
// 	phoneNumber       string
// 	email             string
// }
// type AgencyComponent struct {
// 	ContactInfo  AgencyContactInfo
// 	RelationType AgencyRelationType
// 	Type         AgencyType
// }
//
// // Conditional
// type AgentComponent struct {
// 	AgentContactInfo
// 	AgentRoleCanon
// 	AgencyRef
// }

type LocationGIS struct {
	Latitude  float64
	Longitude float64
}
type LocationComponent struct {
	GIS                 LocationGIS
	Address             string
	ProximateLandmarks  []string
	LocationDescription string
}

// type PersonComponent struct {
// 	PersonRelationType
// 	PersonRelationTime
// 	CallRef
// 	LocationRef
// 	CallBackRef
// }
//
// type VehicleComponent struct {
// 	VehicleRelationType
// 	VehicleRelationTime
// 	LocationRef
// }
