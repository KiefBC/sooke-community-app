package router_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/router"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDB = testdb.Open()
	if testDB == nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestRouter(t *testing.T) {
	r := router.New(testDB)

	tests := []struct {
		name           string
		method         string
		path           string
		wantStatusCode int
	}{
		{
			name:           "GET /api/v1/health returns 200",
			method:         http.MethodGet,
			path:           "/api/v1/health",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "GET /nonexistent returns 404",
			method:         http.MethodGet,
			path:           "/nonexistent",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "GET /api/v1/businesses returns 200",
			method:         http.MethodGet,
			path:           "/api/v1/businesses",
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "GET /api/v1/businesses-typo returns 404",
			method:         http.MethodGet,
			path:           "/api/v1/businesses-typo",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "GET /api/v1/businesses/{slug} with unknown slug returns 404",
			method:         http.MethodGet,
			path:           "/api/v1/businesses/nonexistent-slug",
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "GET /api/v1/categories returns 200",
			method:         http.MethodGet,
			path:           "/api/v1/categories",
			wantStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatusCode)
			}
		})
	}
}

func TestHealthResponses(t *testing.T) {
	tests := []struct {
		name       string
		db         *sql.DB
		wantStatus int
		wantBody   handler.HealthResponse
	}{
		{
			name:       "no database returns 503 degraded",
			db:         nil,
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   handler.HealthResponse{Status: "degraded", DBStatus: "disconnected"},
		},
		{
			name:       "real database returns 200 ok",
			db:         testDB,
			wantStatus: http.StatusOK,
			wantBody:   handler.HealthResponse{Status: "ok", DBStatus: "connected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.New(tt.db)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			var resp handler.HealthResponse
			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}
			if resp != tt.wantBody {
				t.Errorf("body = %+v, want %+v", resp, tt.wantBody)
			}
		})
	}
}
