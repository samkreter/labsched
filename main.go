package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/samkreter/labsched/internal/server"
)

func main() {
	logger := newLogger(os.Getenv)
	slog.SetDefault(logger)

	addr := os.Getenv("LABSCHED_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           server.New(logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info("listening", "addr", addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

// newLogger builds a logger from the environment.
//
//	LOG_FORMAT  "text" for local development, anything else gives JSON
//	LOG_LEVEL   debug, info, warn or error; unset or invalid gives info
//
// It takes getenv instead of calling os.Getenv so tests can pass a fake.
func newLogger(getenv func(string) string) *slog.Logger {
	var level slog.Level
	if err := level.UnmarshalText([]byte(getenv("LOG_LEVEL"))); err != nil {
		level = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: level}

	var handler slog.Handler
	if getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}
