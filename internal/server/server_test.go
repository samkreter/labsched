package server

import (
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

	h := New()
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
