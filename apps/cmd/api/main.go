// The api command is the application's composition root. It will load configuration,
// construct adapters and services, and start the HTTP server without owning business rules.
package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/config"
	postgresrepo "github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/transport/apihttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	incidentRepo := postgresrepo.NewIncidentRepo(db, logger)
	service := application.NewIncidentService(incidentRepo, incidentAuthorizer{}, logger)

	incidents := apihttp.NewIncidentHandler(service, incidentAuthenticator{}, logger)
	handler := apihttp.NewHandler(incidents)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      handler.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

type incidentAuthenticator struct{}

func (incidentAuthenticator) AuthenticateIncident(*http.Request) (application.Identity, error) {
	return application.Identity{}, nil
}

type incidentAuthorizer struct{}

func (incidentAuthorizer) AuthorizePostIncident(context.Context, application.Identity) error {
	return nil
}

func (incidentAuthorizer) AuthorizePatchIncident(context.Context, application.Identity, uuid.UUID) error {
	return nil
}

func (incidentAuthorizer) AuthorizeGetIncident(context.Context, application.Identity, uuid.UUID) error {
	return nil
}
