// This file registers versioned HTTP routes and connects them to handlers. It defines the API
// surface, but delegates request behavior to handlers and application services.
package http

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/observability"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type AuxHandler struct {
}

type Handler struct {
	incidentHandler   IncidentHandler
	assignmentHandler AssignmentHandler
	responseHandler   ResponseHandler
	logger            *slog.Logger
	tracer            trace.Tracer
	traceContext      propagation.TextMapPropagator
}

func NewHandler(logger *slog.Logger, tracer trace.Tracer, traceContext propagation.TextMapPropagator, incidentHandler IncidentHandler) *Handler {
	return &Handler{logger: logger, tracer: tracer, traceContext: traceContext, incidentHandler: incidentHandler}
}

func (h Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	// The literal path wins over the {id} wildcard below, so "options" is never parsed as an ID.
	mux.HandleFunc("GET /v1/incident/options", h.incidentHandler.GetIncidentOptions)
	mux.HandleFunc("GET /v1/incident/{id}", h.incidentHandler.GetIncident)
	mux.HandleFunc("POST /v1/incident", h.incidentHandler.PostIncident)
	mux.HandleFunc("PATCH /v1/incident/{id}", h.incidentHandler.PatchIncident)

	// mux.HandleFunc("GET /v1/checkhealth", h.auxHandler.checkhealth)
	return h.middleware(mux)
}

func (h Handler) middleware(next http.Handler) http.Handler {
	return h.tracingMiddleware(h.authMiddleware(next))
}

// authMiddleware is a pass-through for now: it forwards every request unchanged.
// Incident handlers already authenticate each request themselves (AuthenticateIncident).
// TODO(auth team): add request-wide authentication here and stop calling next on failure.
func (h Handler) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// Adds request headers to context - specifically tracing?
func (h Handler) tracingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := h.traceContext.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.NewString()
		}
		ctx = observability.AddMetadataToContext(ctx, observability.RequestMetadata{
			RequestID: requestID,
			Method:    r.Method,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
