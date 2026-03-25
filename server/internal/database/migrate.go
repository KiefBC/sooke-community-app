package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
)

// migrationLockID is an arbitrary constant used as the Postgres advisory lock key.
// It prevents concurrent instances from running migrations at the same time.
const migrationLockID int64 = 7_238_462_019

func Migrate(db *sql.DB, migrationsPath string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	// Acquire a session-level advisory lock so only one instance migrates at a time.
	// pg_advisory_lock blocks until the lock is available.
	log.Println("Acquiring migration advisory lock...")
	if _, err := db.Exec("SELECT pg_advisory_lock($1)", migrationLockID); err != nil {
		return fmt.Errorf("failed to acquire migration advisory lock: %w", err)
	}
	defer func() {
		if _, err := db.Exec("SELECT pg_advisory_unlock($1)", migrationLockID); err != nil {
			log.Printf("Warning: failed to release migration advisory lock: %v", err)
		}
	}()

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
