package server

import (
	"log/slog"
	"net/http"
)

// New returns the HTTP handler for the service, with request logging.
func New(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz)
	return logRequests(logger, mux)
}

func handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
