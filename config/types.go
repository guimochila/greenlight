// Copyleft (c) 2024, guimochila. Happy Coding.
package config

import "time"

// Datasource configuration.
type Datasource struct {
	DSN          string
	Timeout      time.Duration
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

// Web server configuration.
type Server struct {
	Addr         string
	Port         int
	Readtimeout  time.Duration
	Writetimeout time.Duration
}
