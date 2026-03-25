package database

import (
	"context"
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

	ctx := context.Background()

	// Pin advisory lock to a single connection so lock and unlock
	// execute on the same Postgres session.
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire dedicated connection: %w", err)
	}
	defer conn.Close()

	log.Println("Acquiring migration advisory lock...")
	if _, err := conn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", migrationLockID); err != nil {
		return fmt.Errorf("failed to acquire migration advisory lock: %w", err)
	}
	defer func() {
		var unlocked bool
		if err := conn.QueryRowContext(ctx, "SELECT pg_advisory_unlock($1)", migrationLockID).Scan(&unlocked); err != nil {
			log.Printf("Warning: failed to release migration advisory lock: %v", err)
		} else if !unlocked {
			log.Printf("Warning: pg_advisory_unlock returned false -- lock was not held")
		}
	}()

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
