package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   int
	}{
		{"get ok", http.MethodGet, http.StatusOK},
		{"post rejected", http.MethodPost, http.StatusMethodNotAllowed},
	}

	h := New(slog.New(slog.DiscardHandler))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), tt.method, "/healthz", nil)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tt.want {
				t.Errorf("%s /healthz = %d, want %d", tt.method, rec.Code, tt.want)
			}
		})
	}
}

func TestLogRequestsRecordsStatus(t *testing.T) {
	var buf bytes.Buffer
	h := New(slog.New(slog.NewJSONHandler(&buf, nil)))

	req := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/missing", nil)
	h.ServeHTTP(httptest.NewRecorder(), req)

	var line struct {
		Msg    string `json:"msg"`
		Path   string `json:"path"`
		Status int    `json:"status"`
	}
	if err := json.Unmarshal(buf.Bytes(), &line); err != nil {
		t.Fatalf("log line is not JSON: %v\n%s", err, buf.String())
	}
	if line.Msg != "request" || line.Path != "/missing" || line.Status != http.StatusNotFound {
		t.Errorf("got %+v, want msg=request path=/missing status=404", line)
	}
}
