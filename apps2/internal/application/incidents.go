// This file implements incident use cases, coordinating authorization, domain validation, and
// persistence. It does not know whether the caller was an HTTP handler or another adapter.
package application

import (
	domain "github.com/nawfal-btw/CE1-Group-3/apps2/internal/domain"
	repo "github.com/nawfal-btw/CE1-Group-3/apps2/internal/repository/postgres"
)

type IncidentService struct {
	Incident repo.IncidentORM
}

func (i *IncidentService) Create(input IncidentInput) (*domain.Incident, error) {

}
