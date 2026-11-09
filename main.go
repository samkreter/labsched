package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samkreter/labsched/internal/server"
)

func main() {
	logger := newLogger(os.Getenv)
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("run exit with error", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	addr := os.Getenv("LABSCHED_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:              addr,
		Handler:           server.New(logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	svrErrCh := make(chan error, 1)
	go func() {
		logger.Info("listening", "addr", srv.Addr)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			svrErrCh <- err
		}
		close(svrErrCh)
	}()

	select {
	case err := <-svrErrCh:
		return fmt.Errorf("server: %w", err)
	case <-ctx.Done():
	}

	stop()
	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}

	return <-svrErrCh
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
