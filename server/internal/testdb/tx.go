package testdb

import (
	"database/sql"
	"testing"
)

// WithTx opens the test DB (skipping the test if unavailable), begins a
// transaction, runs each seed against it, and registers a rollback via
// t.Cleanup. The returned *sql.Tx is owned by the test - do not commit.
func WithTx(t *testing.T, seeds ...func(*sql.Tx)) *sql.Tx {
	t.Helper()

	db := Open()
	if db == nil {
		t.Skip("TEST_DATABASE_URL not set")
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("testdb: begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback() })

	for _, seed := range seeds {
		seed(tx)
	}

	return tx
}
