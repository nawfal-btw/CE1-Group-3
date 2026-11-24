// This file translates resource HTTP requests and responses. Its scope is protocol handling;
// resource validation and persistence remain in the application and repository packages.
package http

import (
	"log/slog"

	"github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"go.opentelemetry.io/otel/trace"
)

type ResourceHandler struct {
}

func NewResourceHandler(resourceService *application.ResourceService, logger *slog.Logger, tracer trace.Tracer) *ResourceHandler {
	return nil
}
