package handler_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/handler"
)

func TestHealthHandler(t *testing.T) {
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
		wantStatusCode int
		wantStatus     string
		wantDBStatus   string
	}{
		{
			name:           "no database returns 503 degraded",
			db:             nil,
			wantStatusCode: http.StatusServiceUnavailable,
			wantStatus:     "degraded",
			wantDBStatus:   "disconnected",
		},
	}

	if testDB != nil {
		tests = append(tests, struct {
			name           string
			db             *sql.DB
			wantStatusCode int
			wantStatus     string
			wantDBStatus   string
		}{
			name:           "real database returns 200 ok",
			db:             testDB,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
			wantDBStatus:   "connected",
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler.HealthHandler(tt.db)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
			rec := httptest.NewRecorder()

			h(rec, req)

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

			if resp.DBStatus != tt.wantDBStatus {
				t.Errorf("expected db_status %q, got %q", tt.wantDBStatus, resp.DBStatus)
			}
		})
	}
}
