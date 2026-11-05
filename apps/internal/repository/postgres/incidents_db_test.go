package postgres

// These tests run against a real Postgres (DATABASE_URL, or the local compose DB by default).
// Each test runs inside a transaction in a fresh temporary schema, which is rolled back at the
// end, so nothing is left behind and existing tables are never touched.
// If no database is reachable the tests are skipped.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const migrationUp = "../../../db/migrations/000001_create_incidents.up.sql"

func newTestRepo(t *testing.T) *IncidentRepo {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://anakin:skywalker@localhost:5432/death-star?sslmode=disable"
	}
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		t.Skipf("no database available, skipping: %v", err)
	}
	sqlDB, _ := db.DB()
	t.Cleanup(func() { sqlDB.Close() })

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("could not start transaction: %v", tx.Error)
	}
	t.Cleanup(func() { tx.Rollback() })

	schema := fmt.Sprintf("test_%s", uuid.NewString()[:8])
	migration, err := os.ReadFile(migrationUp)
	if err != nil {
		t.Fatalf("could not read migration: %v", err)
	}
	for _, stmt := range []string{
		"CREATE SCHEMA " + schema,
		"SET LOCAL search_path TO " + schema,
		string(migration),
	} {
		if err := tx.Exec(stmt).Error; err != nil {
			t.Fatalf("test setup failed on %q: %v", stmt, err)
		}
	}
	return NewIncidentRepo(tx)
}

func validIncident(t *testing.T) *domain.Incident {
	t.Helper()
	i, err := domain.NewIncident(domain.IncidentPatch{
		Type:        domain.Field[domain.IncidentType]{Value: 3, Init: true},
		Status:      domain.Field[domain.IncidentStatus]{Value: domain.PLANNING, Init: true},
		Priority:    domain.Field[domain.IncidentPriority]{Value: 2, Init: true},
		Description: domain.Field[string]{Value: "Fire in building", Init: true},
		Location: domain.LocationComponent{
			GIS:                domain.Field[domain.LocationGIS]{Value: domain.LocationGIS{Latitude: 57.0488, Longitude: 9.9217}, Init: true},
			ProximateLandmarks: domain.Field[[]string]{Value: []string{"AAU"}, Init: true},
		},
	})
	if err != nil {
		t.Fatalf("domain rejected test incident: %v", err)
	}
	return i
}

func TestRepoCreateAndGet(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	in := validIncident(t)

	if err := repo.Create(ctx, in); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	out, err := repo.Get(ctx, in.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if out.ID != in.ID || out.Description != in.Description || out.Status != in.Status ||
		out.Location.GIS != in.Location.GIS || out.Location.ProximateLandmarks.Value[0] != "AAU" {
		t.Fatalf("stored incident differs\n in: %+v\nout: %+v", in, out)
	}
	if out.Location.Address.Init {
		t.Fatalf("address was never set, but came back as set: %+v", out.Location.Address)
	}
}

func TestRepoGetNotFound(t *testing.T) {
	repo := newTestRepo(t)

	_, err := repo.Get(context.Background(), domain.IncidentID(uuid.New()))
	if !errors.Is(err, ErrIncidentNotFound) {
		t.Fatalf("expected ErrIncidentNotFound, got %v", err)
	}
}

func TestRepoPatch(t *testing.T) {
	repo := newTestRepo(t)
	ctx := context.Background()
	in := validIncident(t)
	if err := repo.Create(ctx, in); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Patch through the domain, exactly as the application layer will.
	patched, err := domain.PatchIncident(domain.IncidentPatch{
		Status:   domain.Field[domain.IncidentStatus]{Value: domain.IN_PROGRESS, Init: true},
		Location: domain.LocationComponent{Address: domain.Field[string]{Value: "Fredrik Bajers Vej 7", Init: true}},
	}, in)
	if err != nil {
		t.Fatalf("domain rejected patch: %v", err)
	}
	if err := repo.Patch(ctx, patched); err != nil {
		t.Fatalf("Patch failed: %v", err)
	}

	out, err := repo.Get(ctx, in.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if out.Status.Value != domain.IN_PROGRESS || out.Location.Address.Value != "Fredrik Bajers Vej 7" {
		t.Fatalf("patch not stored: %+v", out)
	}
	if out.Description != in.Description {
		t.Fatalf("untouched field changed: %+v", out.Description)
	}
}

func TestRepoPatchNotFound(t *testing.T) {
	repo := newTestRepo(t)

	err := repo.Patch(context.Background(), validIncident(t))
	if !errors.Is(err, ErrIncidentNotFound) {
		t.Fatalf("expected ErrIncidentNotFound, got %v", err)
	}
}
