package postgres

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
)

// A fully filled incident must survive domain -> row -> domain unchanged.
func TestIncidentRowRoundTripFull(t *testing.T) {
	now := time.Date(2026, 11, 5, 12, 0, 0, 0, time.UTC)
	in := &domain.Incident{
		ID:       domain.IncidentID(uuid.New()),
		Type:     domain.Field[domain.IncidentType]{Value: 3, Init: true},
		Status:   domain.Field[domain.IncidentStatus]{Value: domain.PLANNING, Init: true},
		Priority: domain.Field[domain.IncidentPriority]{Value: 2, Init: true},
		Location: domain.LocationComponent{
			GIS:                 domain.Field[domain.LocationGIS]{Value: domain.LocationGIS{Latitude: 57.0488, Longitude: 9.9217}, Init: true},
			Address:             domain.Field[string]{Value: "Fredrik Bajers Vej 7", Init: true},
			ProximateLandmarks:  domain.Field[[]string]{Value: []string{"AAU", "Bus stop"}, Init: true},
			LocationDescription: domain.Field[string]{Value: "Main entrance", Init: true},
		},
		Description:  domain.Field[string]{Value: "Test incident", Init: true},
		CreatedAtUTC: domain.Field[time.Time]{Value: now, Init: true},
		UpdatedAtUTC: domain.Field[time.Time]{Value: now, Init: true},
		Version:      "1",
	}

	out := toDomain(toRow(in))

	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round trip changed the incident\n in: %+v\nout: %+v", in, out)
	}
}

// Fields that were never set must become NULL in the row and stay unset after reading back.
func TestIncidentRowRoundTripEmpty(t *testing.T) {
	in := &domain.Incident{ID: domain.IncidentID(uuid.New())}

	row := toRow(in)
	if row.Type != nil || row.Status != nil || row.Latitude != nil || row.Address != nil || row.ProximateLandmarks != nil {
		t.Fatalf("unset fields should be NULL, got %+v", row)
	}

	out := toDomain(row)
	if !reflect.DeepEqual(in, out) {
		t.Fatalf("round trip changed the incident\n in: %+v\nout: %+v", in, out)
	}
}
