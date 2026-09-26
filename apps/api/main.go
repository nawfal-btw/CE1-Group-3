package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	db "github.com/nawfal-btw/CE1-Group-3/apps/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type incidentCanon struct {
	ID        int64  `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
type incidentPost struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type app struct {
	queries *db.Queries
	db      *sql.DB
}

func (app *app) getIncidentByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		os.Exit(1)
	}
	incidentDB, err := app.queries.GetIncidentByID(r.Context(), int64(id))
	if err != nil {
		os.Exit(1)
	}
	incident_json := incidentCanon{
		ID:        incidentDB.ID,
		Latitude:  incidentDB.Latitude,
		Longitude: incidentDB.Longitude,
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(incident_json); err != nil {
		os.Exit(1)
	}

}
func (app *app) postIncident(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var incident db.AddIncidentParams

	if err := decoder.Decode(&incident); err != nil {
		log.Println("post failed")
		return
	}
	app.queries.AddIncident(r.Context(), incident)
}
func (app *app) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /incident/{id}", app.getIncidentByID)
	// mux.HandleFunc("GET /incident", getAllIncidents)
	mux.HandleFunc("POST /incident", app.postIncident)
	return mux
}
func (app *app) applyMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(app.db, &postgres.Config{})
	if err != nil {
		log.Println("couldnt connect to db")
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		log.Println("couldnt migrate schema")
		return err
	}
	if err := migration.Up(); err != nil {
		log.Println("couldnt migrate no2")
		return err
	}
	return nil
}
func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	pgSocket := flag.String("pgsock", "postgres://anakin:skywalker@localhost:5432/death-star?sslmode=disable", "Postgres connection string")
	flag.Parse()

	pg, err := sql.Open("pgx", *pgSocket)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	if err := pg.Ping(); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	defer pg.Close()

	app := app{
		queries: db.New(pg),
		db:      pg,
	}
	if err := app.applyMigrations("file://sql/migrations"); err != nil {
		log.Println("couldnt apply migrations: %v", err)
	}

	err = http.ListenAndServe(*addr, app.routes())
	if err != nil {
		log.Println("aaa")
		os.Exit(1)
	}
}
