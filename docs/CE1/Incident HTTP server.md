# config.go

This file holds the project configuration.

`Config` stores the server port, the environment, and the database user, password, host, port, and name.

`Load` registers the terminal flags, reads them, and returns a `Config`. The program starts by calling `config.Load()`.

`DSN` joins the database fields into one string. `gorm.Open` accepts that one string, so the fields stay separate in `Config` and are joined only in `DSN`.

```
host=localhost user=anakin password=skywalker dbname=death-star port=5432 sslmode=disable
```

# bin folder

The `bin` folder holds the executable files of the project, as in Let's Go Further. `go build -o bin/api ./cmd/api` writes the API program there. Each build replaces that file. now i have there the http server of the app 
Let's Go Further -- The bin directory will contain our compiled application binaries, ready for deployment to a production server.

# main HTTP server

## prepare main

`main` starts by calling `config.Load()`, which reads the configuration flags.

Then it makes a `*slog.Logger`. The handler is `slog.NewTextHandler` and the destination is `os.Stdout`, so each log line is plain text in the terminal.

`db` is a `*gorm.DB`. `gorm.Open` opens the connection to Postgres and returns this value. It is the handle the program uses to talk to the database. `main` does not run SQL with it. It passes `db` to `NewIncidentRepo`.

The incident handler is built with the service, the authenticator, and the logger. The service carries out the incident actions. The authenticator says who is calling. The logger is passed into `NewIncidentHandler`, so the handler can write logs.

## create the HTTP server

Then `main` creates the HTTP server with `srv := &http.Server{...}`.

`Addr` is the port stored on the `Config` struct, `cfg.Port`.

`Handler` is `handler.Routes()`. That call builds an `http.ServeMux` and returns it. The mux matches each URL to an incident method: `GET`, `POST`, and `PATCH`.

`ErrorLog` is our logger. `slog.NewLogLogger` connects it, and only `Error` level messages go there. Those are the HTTP server's own errors, not the incident request logs.

## start the HTTP server

An `Info` log is written first, with the message `starting server`, so it is clear that the server is starting. The line includes the address and the environment.

Then `srv.ListenAndServe()` starts the server and waits for requests.

If that call returns an error, the error is written as an `Error` log and the program closes with `os.Exit(1)`.

# logs for POST incident

The POST incident logs use the level that fits the event. Steps along the path are `Debug`. Expected user problems are `Warn`. A created incident is `Info`. An unexpected failure is `Error`. The response status matches the case: `201` when the incident is created, `400` for bad JSON or invalid data, `401` when authentication fails, `403` when the caller is not allowed, and `500` when the server fails.

Debug

- `PostIncident`
- `ValidateIncident`

Warn

- `authenticateIncident`
- `decodeIncidentInput`
- `PostIncident`
- `ValidateIncident`

Info

- `PostIncident`

Error

- `PostIncident`

