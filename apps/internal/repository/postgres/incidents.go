package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	domain "github.com/nawfal-btw/CE1-Group-3/apps/internal/domain"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/observability"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// ErrIncidentNotFound is the domain's not-found error, re-exported for callers of this package.
var ErrIncidentNotFound = domain.ErrIncidentNotFound

type IncidentRepo struct {
	db     *gorm.DB
	logger *slog.Logger
	tracer trace.Tracer
}

func NewIncidentRepo(db *gorm.DB, logger *slog.Logger, tracer trace.Tracer) *IncidentRepo {
	return &IncidentRepo{db: db, logger: logger, tracer: tracer}
}

func (r *IncidentRepo) PostIncident(ctx context.Context, incident *domain.Incident) (err error) {
	ctx, span := r.tracer.Start(ctx, "repository.PostIncident")
	defer func() { endSpan(span, err) }()
	logger := observability.AddPersistentAttributes(ctx, r.logger, "repository")
	row := incidentToRow(incident)
	if err = r.db.WithContext(ctx).Create(&row).Error; err != nil {
		logger.ErrorContext(ctx, "incident insert failed", "incident_id", incident.ID, "error", err)
	}
	return err
}

func (r *IncidentRepo) PatchIncident(ctx context.Context, incident *domain.Incident) (err error) {
	ctx, span := r.tracer.Start(ctx, "repository.PatchIncident")
	defer func() { endSpan(span, err) }()
	logger := observability.AddPersistentAttributes(ctx, r.logger, "repository")
	row := incidentToRow(incident)
	result := r.db.WithContext(ctx).
		Model(&incidentRow{ID: row.ID}).
		Select("*").
		Omit("id").
		Updates(row)
	if result.Error != nil {
		logger.ErrorContext(ctx, "incident update failed", "incident_id", incident.ID, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.WarnContext(ctx, "incident update did not find incident", "incident_id", incident.ID)
		return ErrIncidentNotFound
	}
	return nil
}

func (r *IncidentRepo) GetIncident(ctx context.Context, id domain.IncidentID) (incident *domain.Incident, err error) {
	ctx, span := r.tracer.Start(ctx, "repository.GetIncident")
	defer func() { endSpan(span, err) }()
	logger := observability.AddPersistentAttributes(ctx, r.logger, "repository")
	var row incidentRow
	err = r.db.WithContext(ctx).First(&row, "id = ?", uuid.UUID(id)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.WarnContext(ctx, "incident not found", "incident_id", id)
		return nil, ErrIncidentNotFound
	}
	if err != nil {
		logger.ErrorContext(ctx, "incident lookup failed", "incident_id", id, "error", err)
		return nil, err
	}
	return rowToIncident(row), nil
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
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
