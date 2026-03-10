package database_test

import (
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/database"
)

func TestDBConnect(t *testing.T) {
	tests := []struct {
		name        string
		databaseURL string
		envKey      string
		wantErr     bool
	}{
		{
			name:        "empty URL returns error",
			databaseURL: "",
			wantErr:     true,
		},
		{
			name:        "invalid URL returns error",
			databaseURL: "postgres://invalid:invalid@localhost:9999/nonexistent",
			wantErr:     true,
		},
		{
			name:    "valid URL connects successfully",
			envKey:  "TEST_DATABASE_URL",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := tt.databaseURL
			if tt.envKey != "" {
				url = os.Getenv(tt.envKey)
				if url == "" {
					t.Skipf("environment variable %s not set, skipping test", tt.envKey)
				}
			}

			db, err := database.Connect(url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Connect() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if db != nil {
				db.Close()
			}
		})
	}
}
