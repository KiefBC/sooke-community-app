package handler_test

import (
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/testdb"
)

func TestMain(m *testing.M) {
	if testdb.Open() == nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}
