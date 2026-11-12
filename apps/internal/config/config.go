// Package config reads and validates startup configuration once. Other packages receive
// typed configuration or explicit dependencies rather than reading environment variables.
package config

import "flag"

type Config struct {
	Port       int
	Env        string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
}

func Load() Config {
	var cfg Config
	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Postgres user")
	flag.StringVar(&cfg.DBPassword, "db-pass", "", "Postgres password")
	flag.StringVar(&cfg.DBHost, "db-host", "localhost", "Postgres host")
	flag.StringVar(&cfg.DBPort, "db-port", "5432", "Postgres port")
	flag.StringVar(&cfg.DBName, "db-name", "death-star", "Postgres database name")
	flag.Parse()
	return cfg
}
func (c Config) DSN() string {
	return "host=" + c.DBHost +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" port=" + c.DBPort +
		" sslmode=disable"
}
