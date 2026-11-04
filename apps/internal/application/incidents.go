package application

import (
	"context"
	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
)

type LocationBuffer struct {
	Latitude            *float64
	Longitude           *float64
	Address             *string
	ProximateLandmarks  *[]string
	LocationDescription *string
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
	Patch(context.Context, *domain.Incident) error
	Get(context.Context) (*domain.Incident, error)
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
	if res := cmd.Description; res != nil {
		p.Description = domain.Field[string]{Value: *res, Init: true}
	}
	s.CreateLocation(cmd.Location, &p)

	i, err := domain.NewIncident(p)
	if err != nil {
		return nil, err
	}
	if err = s.incidentsRepo.Create(ctx, i); err != nil {
		return nil, err
	}
	return i, nil
}

func (s IncidentService) CreateLocation(loc LocationBuffer, p *domain.IncidentPatch) {
	if lat, long := loc.Latitude, loc.Longitude; lat != nil && long != nil {
		p.Location.GIS.Value.Latitude = *lat
		p.Location.GIS.Value.Longitude = *long
		p.Location.GIS.Init = true
	}
	if res := loc.Address; res != nil {
		p.Location.Address = domain.Field[string]{Value: *res, Init: true}
	}
	if res := loc.LocationDescription; res != nil {
		p.Location.LocationDescription = domain.Field[string]{Value: *res, Init: true}
	}
	if res := loc.ProximateLandmarks; res != nil {
		p.Location.ProximateLandmarks = domain.Field[[]string]{Value: *res, Init: true}
	}
}

// TODO: interact with ORM
func (s IncidentService) PatchIncident(ctx context.Context, cmd PatchIncidentTransport) (*domain.Incident, error) {
	// TODO: get that ORM both ways
	i, err := s.incidentsRepo.Get(ctx)
	if err != nil {
		return nil, err
	}
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
	if res := cmd.Description; res != nil {
		p.Description = domain.Field[string]{Value: *res, Init: true}
	}
	s.CreateLocation(cmd.Location, &p)

	i, err = domain.PatchIncident(p, i)
	if err != nil {
		return nil, err
	}
	return i, nil
}
