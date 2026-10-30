package core

import ()

func (i *Incident) VerifyStatusTransition(p IncidentStatus) bool {
	switch i.Status.Value {
	case PLANNING:
		return p == IN_PROGRESS || p == CANCELLED
	case IN_PROGRESS:
		return p == CANCELLED || p == RESOLVED
	case CANCELLED:
		return p == PLANNING || p == IN_PROGRESS
	case RESOLVED:
		return p == PLANNING || p == IN_PROGRESS
	default:
		return false
	}
}
