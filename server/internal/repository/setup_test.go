package repository_test

import (
	"database/sql"
	"os"
	"testing"

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
