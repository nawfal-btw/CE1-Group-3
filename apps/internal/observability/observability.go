// Package observability defines technology-neutral telemetry contracts.
package observability

import (
	"context"
	"log/slog"
)

type RequestMetadata struct {
	RequestID string
	TraceID   string
	UserID    string
	Method    string
	Endpoint  string
}

type contextKey struct{}

func AddMetadataToContext(ctx context.Context, metadata RequestMetadata) context.Context {
	return context.WithValue(ctx, contextKey{}, metadata)
}

func GetMetadataFromContext(ctx context.Context) RequestMetadata {
	metadata, _ := ctx.Value(contextKey{}).(RequestMetadata)
	return metadata
}

func AddPersistentAttributes(ctx context.Context, logger *slog.Logger, layer string) *slog.Logger {
	metadata := GetMetadataFromContext(ctx)
	return logger.With(
		"layer", layer,
		"request_id", metadata.RequestID,
		"trace_id", metadata.TraceID,
		"user_id", metadata.UserID,
		"http_method", metadata.Method,
		"http_endpoint", metadata.Endpoint,
	)
}
