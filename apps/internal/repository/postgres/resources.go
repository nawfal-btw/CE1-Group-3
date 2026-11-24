// This file implements resource persistence with PostgreSQL and sqlc. It owns query execution
// and mapping, while resource availability rules remain in the domain and application layers.
package postgres

import (
	"database/sql"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type ResourceRepository struct {
}

func NewResourceRepository(db *sql.DB, logger *slog.Logger, tracer trace.Tracer) *ResourceRepository {
	return nil
}
