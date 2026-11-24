// This file implements resource management use cases. It coordinates resource domain rules and
// persistence while keeping transport and SQL details at the application's outer boundaries.
package application

import (
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

type ResourceService struct {
}

func NewResourceService(resourceRepository *postgres.ResourceRepository, logger *slog.Logger, tracer trace.Tracer) *ResourceService {
	return nil
}
