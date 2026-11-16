// The api command is the application's composition root. It will load configuration,
// construct adapters and services, and start the HTTP server without owning business rules.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	applayer "github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	transport "github.com/nawfal-btw/CE1-Group-3/apps/internal/transport/http"
	"go.opentelemetry.io/otel"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type config struct {
	port int
	env  string
	// databaseURL is the Postgres connection string, e.g.
	// postgres://user:password@localhost:5432/dbname?sslmode=disable
	databaseURL string
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
	// The Kubernetes deployment passes -pgsock; locally DATABASE_URL is used if the flag is omitted.
	flag.StringVar(&config.databaseURL, "pgsock", os.Getenv("DATABASE_URL"), "Postgres connection string")
	flag.Parse()

	// should use real attributes not `nil`
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	tracer := otel.Tracer("incident-api")
	meter := otel.Meter("incident-api")

	// 1. Database connection.
	db, err := openDB(config.databaseURL)
	if err != nil {
		logger.Error("cannot connect to database", "error", err)
		os.Exit(1)
	}
	logger.Info("database connection established")

	// 2. Wire the incident use case: repository -> service -> HTTP handler.
	incidentRepo := postgres.NewIncidentRepo(db, logger, tracer)

	// TODO(auth team): replace the placeholder authorizer and authenticator with real ones.
	incidentService, err := applayer.NewIncidentService(incidentRepo, applayer.PlaceholderAuthorizer{}, logger, tracer, meter)
	if err != nil {
		logger.Error("cannot create incident service", "error", err)
		os.Exit(1)
	}
	incidentHandler := transport.NewIncidentHandler(incidentService, transport.PlaceholderAuthenticator{}, logger, tracer)

	// 3. Router with all handlers.
	handler := transport.NewHandler(logger, tracer, otel.GetTextMapPropagator(), incidentHandler)

	app := &application{
		handler: handler,
		logger:  logger,
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

// openDB opens a GORM connection to Postgres and checks that the database actually answers,
// so a wrong address or password stops the program at startup instead of at the first request.
func openDB(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("no database address: set DATABASE_URL or pass -pgsock")
	}
	db, err := gorm.Open(gormpostgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
