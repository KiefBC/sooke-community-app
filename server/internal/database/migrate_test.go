package database_test

import (
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/pressly/goose/v3"
)

// testDB holds the shared database connection for schema tests.
// Migrations are applied once in TestMain before any tests run.
var testDB *sql.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")

	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		os.Exit(0)
	}

	db, err := database.Connect(url)
	if err != nil {
		panic("failed to connect to test database: " + err.Error())
	}

	if err := goose.SetDialect("postgres"); err != nil {
		panic("failed to set goose dialect: " + err.Error())
	}

	migrationsPath := os.Getenv("TEST_MIGRATION_PATH")
	if migrationsPath == "" {
		migrationsPath = "../../../server/migrations"
	}

	// Apply migrations first so the goose tracking table exists,
	// then reset and reapply for a clean schema
	if err := goose.Up(db, migrationsPath); err != nil {
		panic("failed to apply migrations: " + err.Error())
	}
	if err := goose.Reset(db, migrationsPath); err != nil {
		panic("failed to reset migrations: " + err.Error())
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		panic("failed to reapply migrations: " + err.Error())
	}

	testDB = db

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func TestMigrationRoundTrip(t *testing.T) {
	migrationsPath := os.Getenv("TEST_MIGRATION_PATH")
	if migrationsPath == "" {
		migrationsPath = "../../../server/migrations"
	}

	if err := goose.Reset(testDB, migrationsPath); err != nil {
		t.Fatalf("goose.Reset failed: %v", err)
	}

	if err := goose.Up(testDB, migrationsPath); err != nil {
		t.Fatalf("goose.Up (second apply) failed: %v", err)
	}
}

func TestSchemaConstraints(t *testing.T) {
	// Common setup SQL fragments reused across test cases
	const (
		insertRole      = `INSERT INTO user_roles (name) VALUES ('test_role') ON CONFLICT DO NOTHING;`
		insertUser      = `INSERT INTO users (clerk_id, role_id, email, display_name) VALUES ('clerk_test_1', (SELECT id FROM user_roles WHERE name = 'test_role'), 'test@example.com', 'Test User');`
		insertCategory  = `INSERT INTO business_categories (name, slug) VALUES ('Restaurant', 'restaurant') ON CONFLICT DO NOTHING;`
		insertBusiness  = `INSERT INTO businesses (category_id, name, slug, address, latitude, longitude) VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Test Biz', 'test-biz', '123 Main St', 48.37, -123.73);`
		insertEventType = `INSERT INTO event_types (name, slug) VALUES ('Live Music', 'live-music') ON CONFLICT DO NOTHING;`
	)

	tests := []struct {
		name    string
		setup   string
		query   string
		wantErr bool
		check   string
		checkIs int
	}{
		{
			name:    "seed insert and read back",
			setup:   insertRole + insertCategory,
			query:   `INSERT INTO businesses (category_id, name, slug, address, latitude, longitude) VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Sooke Harbour House', 'sooke-harbour-house', '1528 Whiffen Spit Rd', 48.3538, -123.7256);`,
			wantErr: false,
		},
		{
			name:    "slug uniqueness on businesses",
			setup:   insertRole + insertCategory + `INSERT INTO businesses (category_id, name, slug, address, latitude, longitude) VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'First Biz', 'duplicate-slug', '1 Main St', 48.37, -123.73);`,
			query:   `INSERT INTO businesses (category_id, name, slug, address, latitude, longitude) VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Second Biz', 'duplicate-slug', '2 Main St', 48.37, -123.73);`,
			wantErr: true,
		},
		{
			name:    "email uniqueness on users",
			setup:   insertRole + `INSERT INTO users (clerk_id, role_id, email, display_name) VALUES ('clerk_dup_1', (SELECT id FROM user_roles WHERE name = 'test_role'), 'dup@example.com', 'User One');`,
			query:   `INSERT INTO users (clerk_id, role_id, email, display_name) VALUES ('clerk_dup_2', (SELECT id FROM user_roles WHERE name = 'test_role'), 'dup@example.com', 'User Two');`,
			wantErr: true,
		},
		{
			name:    "event location - both business_id and lat/lng set",
			setup:   insertRole + insertUser + insertCategory + insertBusiness + insertEventType,
			query:   `INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, latitude, longitude) VALUES ((SELECT id FROM event_types WHERE slug = 'live-music'), (SELECT id FROM users WHERE email = 'test@example.com'), (SELECT id FROM businesses WHERE slug = 'test-biz'), 'Bad Event', 'bad-event-both', NOW(), 48.37, -123.73);`,
			wantErr: true,
		},
		{
			name:    "event location - neither business_id nor lat/lng set (allowed for orphaned events)",
			setup:   insertRole + insertUser + insertCategory + insertBusiness + insertEventType,
			query:   `INSERT INTO events (event_type_id, submitted_by, name, slug, starts_at) VALUES ((SELECT id FROM event_types WHERE slug = 'live-music'), (SELECT id FROM users WHERE email = 'test@example.com'), 'Bad Event', 'bad-event-neither', NOW());`,
			wantErr: false,
		},
		{
			name:    "event location - business_id only succeeds",
			setup:   insertRole + insertUser + insertCategory + insertBusiness + insertEventType,
			query:   `INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at) VALUES ((SELECT id FROM event_types WHERE slug = 'live-music'), (SELECT id FROM users WHERE email = 'test@example.com'), (SELECT id FROM businesses WHERE slug = 'test-biz'), 'Good Event', 'good-event-biz', NOW());`,
			wantErr: false,
		},
		{
			name:    "event status CHECK rejects invalid value",
			setup:   insertRole + insertUser + insertCategory + insertBusiness + insertEventType,
			query:   `INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES ((SELECT id FROM event_types WHERE slug = 'live-music'), (SELECT id FROM users WHERE email = 'test@example.com'), (SELECT id FROM businesses WHERE slug = 'test-biz'), 'Nonsense Event', 'nonsense-event', NOW(), 'nonsense');`,
			wantErr: true,
		},
		{
			name:    "CASCADE business deletes business_hours",
			setup:   insertRole + insertCategory + insertBusiness + `INSERT INTO business_hours (business_id, day_of_week, open_time, close_time) VALUES ((SELECT id FROM businesses WHERE slug = 'test-biz'), 1, '09:00', '17:00');`,
			query:   `DELETE FROM businesses WHERE slug = 'test-biz';`,
			wantErr: false,
			check:   `SELECT COUNT(*) FROM business_hours WHERE business_id = (SELECT id FROM businesses WHERE slug = 'test-biz');`,
			checkIs: 0,
		},
		{
			name:    "CASCADE menu deletes menu_items",
			setup:   insertRole + insertCategory + insertBusiness + `INSERT INTO menus (business_id, name) VALUES ((SELECT id FROM businesses WHERE slug = 'test-biz'), 'Lunch');` + `INSERT INTO menu_items (menu_id, name, price) VALUES ((SELECT id FROM menus WHERE name = 'Lunch'), 'Fish and Chips', 12.99);`,
			query:   `DELETE FROM menus WHERE name = 'Lunch';`,
			wantErr: false,
			check:   `SELECT COUNT(*) FROM menu_items WHERE menu_id = (SELECT id FROM menus WHERE name = 'Lunch');`,
			checkIs: 0,
		},
		{
			name:    "SET NULL on business delete preserves events",
			setup:   insertRole + insertUser + insertCategory + insertBusiness + insertEventType + `INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at) VALUES ((SELECT id FROM event_types WHERE slug = 'live-music'), (SELECT id FROM users WHERE email = 'test@example.com'), (SELECT id FROM businesses WHERE slug = 'test-biz'), 'Orphaned Event', 'orphaned-event', NOW());`,
			query:   `DELETE FROM businesses WHERE slug = 'test-biz';`,
			wantErr: false,
			check:   `SELECT COUNT(*) FROM events WHERE slug = 'orphaned-event' AND business_id IS NULL;`,
			checkIs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if tt.setup != "" {
				if _, err := tx.Exec(tt.setup); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			_, err = tx.Exec(tt.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("query error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Run optional check query for CASCADE/SET NULL verification
			if tt.check != "" && err == nil {
				var count int
				if err := tx.QueryRow(tt.check).Scan(&count); err != nil {
					t.Fatalf("check query failed: %v", err)
				}
				if count != tt.checkIs {
					t.Errorf("check count = %d, want %d", count, tt.checkIs)
				}
			}
		})
	}
}
