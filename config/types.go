// Copyleft (c) 2024, guimochila. Happy Coding.
package config

import "time"

// Datasource configuration.
type Datasource struct {
	dsn string
}

// Web server configuration.
type Server struct {
	Addr         string
	Port         int
	Readtimeout  time.Duration
	Writetimeout time.Duration
}
