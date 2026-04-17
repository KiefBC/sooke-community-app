package handler_test

import (
	"database/sql"
	"net/http"
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
)

func TestHealthHandler(t *testing.T) {
	var healthDB *sql.DB
	if url := os.Getenv("TEST_DATABASE_URL"); url != "" {
		var err error
		healthDB, err = database.Connect(url)
		if err != nil {
			t.Fatalf("failed to connect to test database: %v", err)
		}
		defer healthDB.Close()
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

	if healthDB != nil {
		tests = append(tests, struct {
			name           string
			db             *sql.DB
			wantStatusCode int
			wantStatus     string
			wantDBStatus   string
		}{
			name:           "real database returns 200 ok",
			db:             healthDB,
			wantStatusCode: http.StatusOK,
			wantStatus:     "ok",
			wantDBStatus:   "connected",
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := testdb.Exec(t, handler.HealthHandler(tt.db), http.MethodGet, "/api/v1/health", nil)
			testdb.AssertStatus(t, rec, tt.wantStatusCode)

			var resp handler.HealthResponse
			testdb.DecodeJSON(t, rec, &resp)

			if resp.Status != tt.wantStatus {
				t.Errorf("expected status %q, got %q", tt.wantStatus, resp.Status)
			}
			if resp.DBStatus != tt.wantDBStatus {
				t.Errorf("expected db_status %q, got %q", tt.wantDBStatus, resp.DBStatus)
			}
		})
	}
}
