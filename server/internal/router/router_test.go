package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/router"
)

func TestRouter(t *testing.T) {
	r := router.New()

	tests := []struct {
		name           string
		method         string
		path           string
		wantStatusCode int
		wantBody       *handler.HealthResponse
	}{
		{
			name:           "GET /api/v1/health returns 200",
			method:         http.MethodGet,
			path:           "/api/v1/health",
			wantStatusCode: http.StatusOK,
			wantBody:       &handler.HealthResponse{Status: "ok"},
		},
		{
			name:           "GET /nonexistent returns 404",
			method:         http.MethodGet,
			path:           "/nonexistent",
			wantStatusCode: http.StatusNotFound,
			wantBody:       nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("expected status code %d, got %d", tt.wantStatusCode, rec.Code)
			}

			if tt.wantBody != nil {
				var resp handler.HealthResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

				if resp != *tt.wantBody {
					t.Errorf("expected body %+v, got %+v", *tt.wantBody, resp)
				}
			}
		})
	}
}
