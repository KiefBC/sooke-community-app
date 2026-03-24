package database

import (
	"database/sql"
	"log"

	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB, migrationsPath string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	return nil
}
