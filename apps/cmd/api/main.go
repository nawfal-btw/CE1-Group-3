// The api command is the application's composition root. It will load configuration,
// construct adapters and services, and start the HTTP server without owning business rules.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/jackc/pgx/v5/stdlib"
	application "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	repository "github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	transport "github.com/nawfal-btw/CE1-Group-3/apps/internal/transport/http"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type app struct {
	handler *transport.Handler
	logger  *slog.Logger
	config  config
}

func main() {
	var config config
	flag.IntVar(&config.port, "port", 4000, "API server port")
	flag.StringVar(&config.env, "environment", "development", "Environment[dev, staging, prod]")
	flag.StringVar(&config.db.dsn, "db-dsn", os.Getenv("DEATHSTAR_DB_DSN"), "PostgreSQL DSN")
	flag.Parse()

	// should use real attributes not `nil`
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tracer := otel.Tracer("incident-api")
	// Concrete persistence and authorization adapters are composed here once
	// their infrastructure configuration is available.
	db, err := setupDB(config)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	handler, err := setupHandler(db, logger, tracer)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &app{
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
	err = server.ListenAndServe()
	logger.Error("server stopped", "error", err)
	os.Exit(1)

}

func setupHandler(db *sql.DB, logger *slog.Logger, tracer trace.Tracer) (*transport.Handler, error) {
	incidentRepository := repository.NewIncidentRepository(db, logger, tracer)
	incidentService := application.NewIncidentService(incidentRepository, logger, tracer, nil)
	incidentHandler := transport.NewIncidentHandler(incidentService, logger, tracer)

	resourceRepository := repository.NewResourceRepository(db, logger, tracer)
	resourceService := application.NewResourceService(resourceRepository, logger, tracer)
	resourceHandler := transport.NewResourceHandler(resourceService, logger, tracer)

	// responseRepository := repository.NewResponseRepository(db, logger, tracer)
	// responseService := application.NewResponseService(responseRepository, logger, tracer)
	// responseHandler := transport.NewResponseHandler(responseService, logger, tracer)

	handler := transport.NewHandler(incidentHandler, resourceHandler, logger, tracer, otel.GetTextMapPropagator())

	return handler, nil

}
func setupDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil

}
