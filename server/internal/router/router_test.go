package router_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/router"
)

func TestRouter(t *testing.T) {
	var testDB *sql.DB
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		var err error
		testDB, err = database.Connect(url)
		if err != nil {
			t.Fatalf("failed to connect to test database: %v", err)
		}
		defer testDB.Close()
	}

	tests := []struct {
		name           string
		db             *sql.DB
		method         string
		path           string
		wantStatusCode int
		wantBody       *handler.HealthResponse
	}{
		{
			name:           "GET /api/v1/health with no DB returns 503",
			db:             nil,
			method:         http.MethodGet,
			path:           "/api/v1/health",
			wantStatusCode: http.StatusServiceUnavailable,
			wantBody:       &handler.HealthResponse{Status: "degraded", DBStatus: "disconnected"},
		},
		{
			name:           "GET /nonexistent returns 404",
			db:             nil,
			method:         http.MethodGet,
			path:           "/nonexistent",
			wantStatusCode: http.StatusNotFound,
			wantBody:       nil,
		},
	}

	if testDB != nil {
		tests = append(tests, struct {
			name           string
			db             *sql.DB
			method         string
			path           string
			wantStatusCode int
			wantBody       *handler.HealthResponse
		}{
			name:           "GET /api/v1/health with DB returns 200",
			db:             testDB,
			method:         http.MethodGet,
			path:           "/api/v1/health",
			wantStatusCode: http.StatusOK,
			wantBody:       &handler.HealthResponse{Status: "ok", DBStatus: "connected"},
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := router.New(tt.db)
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
