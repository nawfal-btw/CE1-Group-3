# Tables
Incidents
	Colums:
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
		AgencyRef 
		LocationRef
		PersonRef
		VehicleRef
		DocumentID (Maybe same as IncidentCanonID, or a list of IDs used to refer to this incident)
CallComponent
	Columns:
		CallQueueID
		CallTypePrimary
		CallDirection
		CallStartTime
		CallAnswerTime
		CallState
		AgentRef
		LocationRef
		PersonRef

AgencyComponent
	Columns:
		AgencyContactInfo
		AgencyType

AgentComponent
	Columns:
		AgentContactInfo
		AgentRoleCanon
		AgencyRef

DispatcherComponent
	Columns:
		AgencyRef
		AgentRef

LocationComponent
	Columns:
		LocationGIS
		CrossStreet
		IntersectingStreet
		LocationDescription
		CellTowerSectorID

PersonComponent
	Columns:
		PersonRelationType
		PersonRelationTime
		CallRef
		LocationRef
		CallBackRef
	
VehicleComponent
	Columns:
		VehicleRelationType
		VehicleRelationTime
		LocationRef

ResourceComponent
	Columns:
		ResourceTypeCanon
		ResourceName
		ResourceStatus
		ResourceTypeInternal
		ResourceAttribute
		AgentRef
		AgencyRef
		LocationRef (Only when availability)
		VehicleRef

Allocation
	Columns:
		IncidentCanonID
		ResourceRef
		VehicleRef
		AgentRef
		status
		
		