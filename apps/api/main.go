package main

import (
	"database/sql"
	"flag"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/nawfal-btw/CE1-Group-3/apps/sql/internal"
	"net/http"
)

type incident struct {
	ID        int64  `json:"id"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

type app struct {
	queries *sqlinternal.Queries
}

func main() {
	addr := flag.String("addr", ":8080", "HTTP network address")
	// pgSocket := flag.String("pgsock", "postgres://nawfalaffald@/postgres?host=/home/nawfalaffald/Projects/TheGrift/db/pgsock&sslmode=disable", "Postgres Socket Path"

	err := http.ListenAndServe(*addr, app.routes())
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
