package domain

import "testing"

func TestIsValidType(t *testing.T) {
	if !IsValidType(1) {
		t.Error("type 1 is in the catalogue, expected valid")
	}
	if IsValidType(999) {
		t.Error("type 999 is not in the catalogue, expected invalid")
	}
}

func TestIsValidPriority(t *testing.T) {
	if !IsValidPriority(1) {
		t.Error("priority 1 is in the catalogue, expected valid")
	}
	if IsValidPriority(0) {
		t.Error("priority 0 is not in the catalogue, expected invalid")
	}
}

func TestPostIncidentRejectsUnknownType(t *testing.T) {
	unknown := IncidentType(999)
	_, err := PostIncident(IncidentPatch{Type: Option[IncidentType]{Value: &unknown, Set: true}})
	if err == nil {
		t.Fatal("expected an error for an unknown incident type")
	}
}
