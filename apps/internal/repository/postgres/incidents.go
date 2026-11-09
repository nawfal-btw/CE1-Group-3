package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
	"gorm.io/gorm"
)

var ErrIncidentNotFound = errors.New("incident not found")

type IncidentRepo struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewIncidentRepo(db *gorm.DB, logger *slog.Logger) *IncidentRepo {
	return &IncidentRepo{db: db, logger: logger.With("layer", "repository")}
}

func (r *IncidentRepo) PostIncident(ctx context.Context, incident *domain.Incident) error {
	r.logger.InfoContext(ctx, "storing incident", "incident_id", incident.ID)
	row := incidentToRow(incident)
	return r.db.WithContext(ctx).Create(&row).Error
}

func (r *IncidentRepo) PatchIncident(ctx context.Context, incident *domain.Incident) error {
	r.logger.InfoContext(ctx, "updating incident", "incident_id", incident.ID)
	row := incidentToRow(incident)
	result := r.db.WithContext(ctx).
		Model(&incidentRow{ID: row.ID}).
		Select("*").
		Omit("id").
		Updates(row)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrIncidentNotFound
	}
	return nil
}

func (r *IncidentRepo) GetIncident(ctx context.Context, id domain.IncidentID) (*domain.Incident, error) {
	r.logger.InfoContext(ctx, "loading incident", "incident_id", id)
	var row incidentRow
	err := r.db.WithContext(ctx).First(&row, "id = ?", uuid.UUID(id)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrIncidentNotFound
	}
	if err != nil {
		return nil, err
	}
	return rowToIncident(row), nil
}

// incidentRow is deliberately separate from domain.Incident so GORM remains a
// repository concern. A nil column value maps to an uninitialized domain field.
type incidentRow struct {
	ID                  uuid.UUID `gorm:"column:id;type:uuid;primaryKey"`
	Type                *int16    `gorm:"column:type"`
	Status              *int8     `gorm:"column:status"`
	Priority            *int8     `gorm:"column:priority"`
	Latitude            *float64  `gorm:"column:latitude"`
	Longitude           *float64  `gorm:"column:longitude"`
	Address             *string   `gorm:"column:address"`
	LocationDescription *string   `gorm:"column:location_description"`
	Description         *string   `gorm:"column:description"`
}

func (incidentRow) TableName() string {
	return "incidents"
}

func incidentToRow(incident *domain.Incident) incidentRow {
	return incidentRow{
		ID:                  uuid.UUID(incident.ID),
		Type:                domainOptionToColumnValue(incident.Type, func(value domain.IncidentType) int16 { return int16(value) }),
		Status:              domainOptionToColumnValue(incident.Status, func(value domain.IncidentStatus) int8 { return int8(value) }),
		Priority:            domainOptionToColumnValue(incident.Priority, func(value domain.IncidentPriority) int8 { return int8(value) }),
		Latitude:            domainOptionToColumnValue(incident.Location.Latitude, func(value float64) float64 { return value }),
		Longitude:           domainOptionToColumnValue(incident.Location.Longitude, func(value float64) float64 { return value }),
		Address:             domainOptionToColumnValue(incident.Location.Address, func(value string) string { return value }),
		LocationDescription: domainOptionToColumnValue(incident.Location.LocationDescription, func(value string) string { return value }),
		Description:         domainOptionToColumnValue(incident.Description, func(value string) string { return value }),
	}
}

func rowToIncident(row incidentRow) *domain.Incident {
	return &domain.Incident{
		ID:          domain.IncidentID(row.ID),
		Type:        columnValueToDomainOption(row.Type, func(value int16) domain.IncidentType { return domain.IncidentType(value) }),
		Status:      columnValueToDomainOption(row.Status, func(value int8) domain.IncidentStatus { return domain.IncidentStatus(value) }),
		Priority:    columnValueToDomainOption(row.Priority, func(value int8) domain.IncidentPriority { return domain.IncidentPriority(value) }),
		Description: columnValueToDomainOption(row.Description, func(value string) string { return value }),
		Location: domain.Location{
			Latitude:            columnValueToDomainOption(row.Latitude, func(value float64) float64 { return value }),
			Longitude:           columnValueToDomainOption(row.Longitude, func(value float64) float64 { return value }),
			Address:             columnValueToDomainOption(row.Address, func(value string) string { return value }),
			LocationDescription: columnValueToDomainOption(row.LocationDescription, func(value string) string { return value }),
		},
	}
}

// domainOptionToColumnValue maps an unset or zero-value domain option to SQL NULL.
func domainOptionToColumnValue[Domain, Column any](option domain.Option[Domain], convert func(Domain) Column) *Column {
	if !option.Set || option.Value == nil {
		return nil
	}
	value := convert(*option.Value)
	return &value
}

// columnValueToDomainOption maps SQL NULL to an unset domain option.
func columnValueToDomainOption[Column, Domain any](value *Column, convert func(Column) Domain) domain.Option[Domain] {
	if value == nil {
		return domain.Option[Domain]{}
	}
	converted := convert(*value)
	return domain.Option[Domain]{Value: &converted, Set: true}
}
