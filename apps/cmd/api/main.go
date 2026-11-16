// The api command is the application's composition root. It will load configuration,
// construct adapters and services, and start the HTTP server without owning business rules.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	transport "github.com/nawfal-btw/CE1-Group-3/apps/internal/transport/http"
	"go.opentelemetry.io/otel"
)

type config struct {
	port int
	env  string
}

type application struct {
	handler *transport.Handler
	logger  *slog.Logger
	config  config
}

func main() {
	var config config
	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.env, "environment", "development", "Environment[dev, staging, prod]")
	flag.Parse()

	// should use real attributes not `nil`
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tracer := otel.Tracer("incident-api")
	handler := transport.NewHandler(logger, tracer, otel.GetTextMapPropagator())

	app := &application{
		handler: handler,
		config:  config,
	}
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.port),
		Handler:      app.handler.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("server starting...", "addr", server.Addr, "env", config.env)
	err := server.ListenAndServe()
	logger.Error("server stopped", "error", err)
	os.Exit(1)

}
