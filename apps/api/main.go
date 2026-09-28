package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

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

func (app *app) applyMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(app.db, &postgres.Config{})
	if err != nil {
		log.Println("couldnt connect to db: ", err)
		return err
	}
	migration, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		log.Println("DB instance and migration interop issue: ", err)
		return err
	}
	if err := migration.Up(); err != nil {
		log.Println("Migrations unsuccessful: ", err)
		return err
	}
	return nil
}
func setupDB(pgSocket *string) app {
	pg, err := sql.Open("pgx", *pgSocket)
	if err != nil {
		log.Fatalln("db connection failed: ", err, "\n Exiting...")
	}
	if err := pg.Ping(); err != nil {
		log.Fatalln("db ping failed: ", err, "\n Exiting...")
	}
	app := app{
		queries: db.New(pg),
		db:      pg,
	}
	if err := app.applyMigrations("file://sql/migrations"); err != nil {
		log.Println("No migrations applied: ", err)
	}
	return app

}
func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	pgSocket := flag.String("pgsock", "postgres://anakin:skywalker@localhost:5432/death-star?sslmode=disable", "Postgres connection string")
	flag.Parse()

	app := setupDB(pgSocket)

	err := http.ListenAndServe(*addr, app.routes())
	if err != nil {
		log.Fatalln("Failed to init webserver: ", err)
	}
}
