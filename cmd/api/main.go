// Copyleft (c) 2024, guimochila. Happy Coding.
package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/guimochila/greenlight/config"
)

type application struct {
	config config.Config
	logger *slog.Logger
}

func main() {
	var cfg config.Config

	config.New(&cfg)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  cfg.Server.Readtimeout,
		WriteTimeout: cfg.Server.Writetimeout,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
