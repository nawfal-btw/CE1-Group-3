package core

import (
	"time"
)

type Patch[T any] struct {
	Value T
	Apply bool
}

type IncidentPatch struct {
	Type        Patch[IncidentType]
	Status      Patch[IncidentStatus]
	Priority    Patch[IncidentPriority]
	Location    Patch[LocationComponent]
	Description Patch[string]

	UpdatedAtUTC Patch[time.Time]
	ClosedAtUTC  Patch[time.Time]

	// Owner Patch[AgencyRef]
}

func GenIncident() *Incident {
	// Default initializes to zero i.e. correct for Field.Initialized
	return &Incident{}
}
func (i *Incident) PatchMulIncident(p IncidentPatch) {
	if p.Type.Apply {
		i.PatchType(p.Type.Value)
	}
	if p.Status.Apply {
		i.PatchStatus(p.Status.Value)
	}
	if p.Priority.Apply {
		i.PatchPriority(p.Priority.Value)
	}
	if p.Location.Apply {
		i.PatchLocation(p.Location.Value)
	}
	if p.Description.Apply {
		i.PatchDescription(p.Description.Value)
	}
	if p.UpdatedAtUTC.Apply {
		i.PatchUpdatedAtUTC(p.UpdatedAtUTC.Value)
	}
	if p.ClosedAtUTC.Apply {
		i.PatchClosedAtUTC(p.ClosedAtUTC.Value)
	}
}
func (i *Incident) PatchType(t IncidentType) {
	i.Type.Value = t
	i.Type.Initialized = true
}
func (i *Incident) PatchStatus(s IncidentStatus) {
	i.Status.Value = s
	i.Status.Initialized = true
}
func (i *Incident) PatchPriority(p IncidentPriority) {
	i.Priority.Value = p
	i.Priority.Initialized = true
}
func (i *Incident) PatchLocation(l LocationComponent) {
	i.Location.Value = l
	i.Location.Initialized = true
}

func (i *Incident) PatchDescription(d string) {
	i.Description.Value = d
	i.Description.Initialized = true
}

func (i *Incident) PatchUpdatedAtUTC(t time.Time) {
	i.UpdatedAtUTC.Value = t
	i.UpdatedAtUTC.Initialized = true
}

func (i *Incident) PatchClosedAtUTC(t time.Time) {
	i.ClosedAtUTC.Value = t
	i.ClosedAtUTC.Initialized = true
}
