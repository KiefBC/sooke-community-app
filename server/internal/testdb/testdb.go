package testdb

import (
	"context"
	"database/sql"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/pressly/goose/v3"
)

// testdbLockID is an arbitrary constant used as the Postgres advisory lock key.
// It serializes schema reset across parallel test-package processes.
const testdbLockID int64 = 8_349_571_020

var (
	once sync.Once
	db   *sql.DB
)

// Open returns a shared *sql.DB connected to the test database with
// migrations already applied. It returns nil when TEST_DATABASE_URL is
// not set, allowing callers to skip gracefully.
//
// A Postgres advisory lock serializes the schema reset across parallel
// test-package processes. The first process to acquire the lock drops and
// recreates the schema; subsequent processes find migrations already
// applied and proceed immediately.
func Open() *sql.DB {
	_ = godotenv.Load("../../../.env")

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		return nil
	}

	once.Do(func() {
		conn, err := database.Connect(url)
		if err != nil {
			panic("testdb: failed to connect: " + err.Error())
		}

		if err := goose.SetDialect("postgres"); err != nil {
			panic("testdb: failed to set goose dialect: " + err.Error())
		}

		migrationsPath := os.Getenv("TEST_MIGRATION_PATH")
		if migrationsPath == "" {
			migrationsPath = "../../../server/migrations"
		}

		// Acquire advisory lock on a dedicated connection so that lock
		// and unlock run on the same Postgres session.
		ctx := context.Background()
		lockConn, err := conn.Conn(ctx)
		if err != nil {
			panic("testdb: failed to get lock connection: " + err.Error())
		}
		defer lockConn.Close()

		if _, err := lockConn.ExecContext(ctx, "SELECT pg_advisory_lock($1)", testdbLockID); err != nil {
			panic("testdb: failed to acquire advisory lock: " + err.Error())
		}
		defer func() {
			_, _ = lockConn.ExecContext(ctx, "SELECT pg_advisory_unlock($1)", testdbLockID)
		}()

		// Check whether migrations have already been applied by another
		// test-package process. If goose_db_version exists and is at the
		// latest version, skip
		var exists bool
		err = lockConn.QueryRowContext(ctx,
			"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'goose_db_version')").Scan(&exists)
		if err != nil {
			panic("testdb: failed to check goose table: " + err.Error())
		}

		if !exists {
			if _, err := conn.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
				panic("testdb: failed to reset schema: " + err.Error())
			}
		}

		if err := goose.Up(conn, migrationsPath); err != nil {
			panic("testdb: failed to apply migrations: " + err.Error())
		}

		db = conn
	})

	return db
}
