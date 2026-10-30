# Responsibility
Interops with the API. Ensures that API calls1
## EIDO Based
IncidentComponent
	IncidentCanonID
	IncidentInternalID
	IncidentTypeInternalCode
	IncidentTypeInternalText
	IncidentTypeCanonCode
	IncidentTypeCanonText
	IncidentStatusInternalCode
	IncidentStatusCanonCode
	IncidentPriorityInternalText
	IncidentPriorityCanon
	IncidentOwner
	CallRef
	CallBackRef
	DispatchRef
	EmergencyResourceRef [(Resource1)]
	AgentRef
	AgencyRef [(*Agency, Relation), (*Agency2, AnotherRelation)]
	LocationRef
	PersonRef
	VehicleRef
	DocumentID (Maybe same as IncidentCanonID, or a list of IDs used to refer to this incident)
	
CallComponent
	CallQueueID
	CallTypePrimary
	CallDirection
	CallStartTime
	CallAnswerTime
	CallState
	AgentRef
	CallBackRef
	LocationRef
	PersonRef
	
CallBackComponent
	CallBackInfoURI
	
DispatchComponent
	AgencyRef
	AgentRef
	EmergencyResourceRef 
	
ResourceComponent
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
	
AgencyComponent
	AgencyContactInfo
	AgencyType
	
AgentComponent
	AgentContactInfo
	AgentRoleCanon
	AgencyRef

LocationComponent
	LocationGIS
	CrossStreet
	IntersectingStreet
	LocationDescription
	CellTowerSectorID
	
PersonComponent
	PersonRelationType
	PersonRelationTime
	CallRef
	LocationRef
	CallBackRef
	
VehicleComponent
	VehicleRelationType
	VehicleRelationTime
	LocationRef

Incident
- ID
- Type
- Location (Longitude, Lattitude)
- Severity
- Status_id
- Created_at
- Description
- The person who called


Status
- Status_ID
- Incident_ID
- Last_modified
- Status:

Resource
- ID
- Type
- Organization
- Location
- Availability
- Capacity


Dispatcher:
- See available resources
- Attach resources to incident
	- w/ Suggestions based on location
- track resources 
- track incidents
- get info from dispatched units


## NW
Call-taker -> Dispatcher

