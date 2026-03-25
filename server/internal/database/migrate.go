package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, migrationsPath string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}
