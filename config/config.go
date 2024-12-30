// Copyleft (c) 2024, guimochila. Happy Coding.
package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

// Application version.
const Version = "1.0.0"

// Default web server port.
const DefaultPort = 4000

// Default timeouts for webserver operations.
const (
	Readtimeout       = 5 * time.Second
	Writetimeout      = 10 * time.Second
	Datasourcetimeout = 5 * time.Second
)

// Default datasource connection limits.
const (
	MaxOpenConns    = 1
	MaxIdleConns    = 25
	MaxIdleTime     = 15 * time.Minute
	MaxQueryTimeout = 3 * time.Second
)

// Global application configuration.
type Config struct {
	Server     Server
	Env        string
	Datasource Datasource
	Version    string
}

func New(config *Config) {
	// Parse configuration from command line.
	flag.IntVar(&config.Server.Port, "port", DefaultPort, "API server port")
	flag.StringVar(&config.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&config.Datasource.DSN, "datasource", os.Getenv("DB_DSN"), "Postgres DSN")
	flag.IntVar(&config.Datasource.MaxOpenConns, "db-max-open-conns", MaxOpenConns, "Postgres max open connection")
	flag.IntVar(&config.Datasource.MaxIdleConns, "db-max-idle-coons", MaxIdleConns, "Postgres max idle connection")
	flag.DurationVar(&config.Datasource.MaxIdleTime, "db-max-idle-time", MaxIdleTime, "Postgres max connection idle time")
	flag.Parse()

	// Add default configuration for server.
	config.Version = Version
	config.Server.Readtimeout = Readtimeout
	config.Server.Writetimeout = Writetimeout
	config.Server.Addr = fmt.Sprintf(":%d", config.Server.Port)

	// Datasource settings.
	config.Datasource.Timeout = Datasourcetimeout
}
