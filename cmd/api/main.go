// Copyleft (c) 2024, guimochila. Happy Coding.
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/guimochila/greenlight/config"
	_ "github.com/lib/pq"
)

type application struct {
	config config.Config
	logger *slog.Logger
}

func main() {
	var cfg config.Config

	config.New(&cfg)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := initDB(logger, cfg.Datasource)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection poll established")

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

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func initDB(logger *slog.Logger, ds config.Datasource) (*sql.DB, error) {
	db, err := sql.Open("postgres", ds.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(ds.MaxOpenConns)
	db.SetMaxIdleConns(ds.MaxIdleConns)
	db.SetConnMaxIdleTime(ds.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), ds.Timeout)

	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		logger.Error("failed to ping database", "error", err.Error())

		return nil, err
	}

	return db, nil
}
