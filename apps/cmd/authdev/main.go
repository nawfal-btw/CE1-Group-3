//go:build authdev

// Development-only authentication harness. Never deploy as the application API.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/auth"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/repository/postgres"
)

const (
	address = "127.0.0.1:8444"
	origin  = "https://" + address
	issuer  = "https://localhost:8443/realms/ce1-dev"
)

func main() {
	if err := run(); err != nil {
		// run returns only sanitized messages, never underlying driver/OIDC errors.
		log.Fatal(err)
	}
}

func run() error {
	for _, name := range []string{"CE1_AUTH_TEST_DATABASE_URL", "AUTH_BFF_CLIENT_SECRET", "AUTHDEV_CA_FILE", "AUTHDEV_CERT_FILE", "AUTHDEV_KEY_FILE"} {
		if os.Getenv(name) == "" {
			return errors.New("authdev requires database URL, client secret, CA, certificate and key environment settings")
		}
	}
	certificate, err := tls.LoadX509KeyPair(os.Getenv("AUTHDEV_CERT_FILE"), os.Getenv("AUTHDEV_KEY_FILE"))
	if err != nil {
		return errors.New("authdev could not load its TLS certificate and key")
	}
	caPEM, err := os.ReadFile(os.Getenv("AUTHDEV_CA_FILE"))
	if err != nil {
		return errors.New("authdev could not read the development CA certificate")
	}
	roots := x509.NewCertPool()
	if !roots.AppendCertsFromPEM(caPEM) {
		return errors.New("authdev development CA certificate is invalid")
	}
	config, err := pgx.ParseConfig(os.Getenv("CE1_AUTH_TEST_DATABASE_URL"))
	if err != nil {
		return errors.New("authdev test database configuration is invalid")
	}
	// Check the resolved driver configuration before any connection, including
	// environment/service-file defaults and multi-host fallbacks.
	if config.Host != "127.0.0.1" || config.Port != 55432 || config.Database != "ce1_auth_test" || config.User != "ce1_auth_test" {
		return errors.New("authdev refuses any database target other than the isolated ce1_auth_test fixture")
	}
	for _, fallback := range config.Fallbacks {
		if fallback.Host != config.Host || fallback.Port != config.Port {
			return errors.New("authdev refuses alternate database hosts or ports")
		}
	}
	config.ConnectTimeout = 5 * time.Second
	db := sql.OpenDB(stdlib.GetConnector(*config))
	defer db.Close()
	db.SetMaxOpenConns(5)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var database, user string
	err = db.QueryRowContext(ctx, "SELECT current_database(), current_user").Scan(&database, &user)
	cancel()
	if err != nil || database != "ce1_auth_test" || user != "ce1_auth_test" {
		return errors.New("authdev could not verify the isolated database identity")
	}
	validator, err := auth.NewOIDC(context.Background(), auth.OIDCConfig{
		IssuerURL: issuer, ClientID: "ce1-bff-dev", APIAudience: "ce1-api",
		HumanRoles: []string{"administrator", "dispatcher", "operator"}, RootCAs: roots,
	})
	if err != nil {
		return errors.New("authdev HTTPS OIDC discovery failed; verify issuer and certificate trust")
	}
	login, err := auth.NewLoginHandler(validator, postgres.NewAuthStore(db), auth.LoginConfig{
		PublicOrigin: origin, ClientSecret: os.Getenv("AUTH_BFF_CLIENT_SECRET"),
	})
	if err != nil {
		return errors.New("authdev login configuration is invalid")
	}
	mux := http.NewServeMux()
	mux.Handle("/auth/", login.Routes()) // Preserve the complete /auth/... path.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Referrer-Policy", "same-origin")
		_, _ = io.WriteString(w, `<!doctype html><html lang="en"><meta charset="utf-8">
<title>CE1 local authentication test</title><h1>Local authentication test</h1>
<form method="post" action="/auth/login"><button>Log in with Keycloak</button></form>
<p><a href="/auth/me">Check current session</a></p>
<form method="post" action="/auth/logout"><button>Log out of this application session</button></form>
</html>`)
	})
	server := &http.Server{
		Addr:              address,
		TLSConfig:         &tls.Config{MinVersion: tls.VersionTLS12, Certificates: []tls.Certificate{certificate}},
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 30 * time.Second,
		// No access logging or raw request diagnostics, including callback queries.
		ErrorLog: log.New(io.Discard, "", 0),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			w.Header().Set("Referrer-Policy", "no-referrer")
			w.Header().Set("X-Frame-Options", "DENY")
			if r.TLS == nil || r.Host != address {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}
			mux.ServeHTTP(w, r)
		}),
	}
	log.Print("Development-only authentication harness: " + origin)
	if err := server.ListenAndServeTLS("", ""); err != nil {
		return errors.New("authdev HTTPS server stopped")
	}
	return nil
}
