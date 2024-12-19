// Copyleft (c) 2024, guimochila. Happy Coding.
package config

import (
	"flag"
	"fmt"
	"time"
)

// Application version.
const Version = "1.0.0"

// Default web server port.
const DefaultPort = 4000

// Default timeouts for webserver operations.
const (
	Readtimeout  = 5 * time.Second
	Writetimeout = 10 * time.Second
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
	flag.Parse()

	// Add default configuration for server.
	config.Version = Version
	config.Server.Readtimeout = Readtimeout
	config.Server.Writetimeout = Writetimeout
	config.Server.Addr = fmt.Sprintf(":%d", config.Server.Port)
}
