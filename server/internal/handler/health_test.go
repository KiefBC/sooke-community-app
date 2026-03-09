package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/handler"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatusCode int
		wantStatus     string
	}{
		{
			name:           "returns 200 and status ok",
			method:         http.MethodGet,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/api/v1/health", nil)
			rec := httptest.NewRecorder()

			handler.HealthHandler(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status code %d, got %d", tt.wantStatusCode, rec.Code)
			}

			var resp handler.HealthResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if resp.Status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, resp.Status)
			}
		})
	}
}
