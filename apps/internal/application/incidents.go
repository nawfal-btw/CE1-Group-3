// This file implements incident use cases, coordinating authorization, domain validation, and
// persistence. It does not know whether the caller was an HTTP handler or another adapter.
package application

import (
	"context"
	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
	// repo "github.com/nawfal-btw/CE1-Group-3/apps2/internal/repository/postgres"
	// websocket "github.com/nawfal-btw/CE1-Group-3/apps2/internal/transport/websocket"
)

type LocationBuffer struct {
	Latitude            float64
	Longitude           float64
	Address             string
	ProximateLandmarks  []string
	LocationDescription string
}
type PatchIncidentTransport struct {
	ID          *uuid.UUID
	Type        *int16
	Status      *int8
	Priority    *int8
	Location    *LocationBuffer
	Description *string
	// Owner        AgencyRef
}

type IncidentRepository interface {
	Create(context.Context, *domain.Incident) error
}
type IncidentService struct {
	incidentsRepo IncidentRepository
	// incidentsLogger IncidentLogger
}

func NewIncidentService(i IncidentRepository) IncidentService {
	return IncidentService{incidentsRepo: i}
}

func (s IncidentService) CreateIncident(ctx context.Context, cmd PatchIncidentTransport) (*domain.Incident, error) {
	p := domain.IncidentPatch{}

	if res := cmd.Type; res != nil {
		p.Type = domain.Field[domain.IncidentType]{Value: domain.IncidentType(*res), Init: true}
	}
	if res := cmd.Status; res != nil {
		p.Status = domain.Field[domain.IncidentStatus]{Value: domain.IncidentStatus(*res), Init: true}
	}
	if res := cmd.Priority; res != nil {
		p.Priority = domain.Field[domain.IncidentPriority]{Value: domain.IncidentPriority(*res), Init: true}
	}
	if res := cmd.Location; res != nil {
		loc := cmd.Location
		p.Location = domain.LocationComponent{
			GIS: domain.Field[domain.LocationGIS]{
				Value: domain.LocationGIS{
					Latitude:  loc.Latitude,
					Longitude: loc.Longitude,
				},
				Init: true,
			},
			Address: domain.Field[string]{
				Value: loc.Address,
				Init:  true,
			},
			ProximateLandmarks: domain.Field[[]string]{
				Value: loc.ProximateLandmarks,
				Init:  true,
			},
			LocationDescription: domain.Field[string]{
				Value: loc.LocationDescription,
				Init:  true,
			},
		}
	}
	if res := cmd.Location; res != nil {
		p.Location = *res
	}

	i, err := domain.NewIncident(p)
	if err != nil {
		return &domain.Incident{}, err
	}
	// TODO: interact with ORM
	if err = s.incidentsRepo.Create(ctx, i); err != nil {
		return &domain.Incident{}, err
	}
	return i, nil
}
func (s IncidentService) PatchIncident(ctx context.Context, cmd PatchIncidentTransport) (*domain.Incident, error) {
	// TODO: get that ORM both ways

	return nil, nil
}
