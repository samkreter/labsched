package server

import (
	"log/slog"
	"net/http"
	"time"
)

// logRequests writes one log line per request, after the handler finishes.
func logRequests(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		logger.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// statusRecorder wraps a ResponseWriter to remember the status code.
// Handlers that never call WriteHeader send 200, hence the default.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// Unwrap lets http.ResponseController reach the original writer,
// for features like flushing that the wrapper doesn't implement.
func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}
