package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/samkreter/labsched/internal/server"
)

func main() {
	opts := &slog.HandlerOptions{Level: slog.LevelDebug}
	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           server.New(logger),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("listening on :8080")
	log.Fatal(srv.ListenAndServe())
}
