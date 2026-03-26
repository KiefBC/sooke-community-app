package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/pressly/goose/v3"
)

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

	// Drop all tables (including goose tracking) for a guaranteed clean state.
	if _, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		panic("failed to reset schema: " + err.Error())
	}
	if err := goose.Up(db, migrationsPath); err != nil {
		panic("failed to apply migrations: " + err.Error())
	}

	testDB = db

	code := m.Run()

	db.Close()
	os.Exit(code)
}

func TestGetBusinessBySlug(t *testing.T) {
	const setup = `
		INSERT INTO business_categories (name, slug) VALUES ('Restaurant', 'restaurant') ON CONFLICT DO NOTHING;
		INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
						VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Sooke Harbour House', 'sooke-harbour-house', '1528 Whiffen Spit Rd', 48.3538, -123.7256);
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
						VALUES
								((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 3, '09:00', '17:00', false),
								((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 1, '09:00', '17:00', false),
								((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 2, '09:00', '17:00', false);
		INSERT INTO menus (business_id, name) VALUES ((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 'Dinner');
		INSERT INTO menu_items (menu_id, name, price) VALUES ((SELECT id FROM menus WHERE name = 'Dinner'), 'Fish and Chips', '12.99');
	`

	tests := []struct {
		name      string
		slug      string
		wantNil   bool
		wantName  string
		wantHours int
		wantMenus int
		wantItems int
	}{
		{
			name:      "existing business returns details",
			slug:      "sooke-harbour-house",
			wantNil:   false,
			wantName:  "Sooke Harbour House",
			wantHours: 3,
			wantMenus: 1,
			wantItems: 1,
		},
		{
			name:      "nonexistent business returns nil",
			slug:      "nonexistent-slug",
			wantNil:   true,
			wantName:  "",
			wantHours: 0,
			wantMenus: 0,
			wantItems: 0,
		},
		{
			name:      "menus and items are included",
			slug:      "sooke-harbour-house",
			wantNil:   false,
			wantName:  "Sooke Harbour House",
			wantHours: 3,
			wantMenus: 1,
			wantItems: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(setup); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			business, err := repository.GetBusinessBySlug(context.Background(), tx, tt.slug)
			if err != nil {
				t.Fatalf("GetBusinessBySlug returned error: %v", err)
			}

			// if we expect nil, assert that business is nil and return early to avoid dereferencing nil in further assertions
			if tt.wantNil {
				if business != nil {
					t.Errorf("got business %+v, want nil", business)
				}
				return
			}

			if business == nil {
				t.Fatalf("got nil, want business with name %q", tt.wantName)
			}

			if business.Name != tt.wantName {
				t.Errorf("got name %q, want %q", business.Name, tt.wantName)
			}
			if len(business.Hours) != tt.wantHours {
				t.Errorf("got %d hours, want %d", len(business.Hours), tt.wantHours)
			}

			for i := 1; i < len(business.Hours); i++ {
				if business.Hours[i].DayOfWeek < business.Hours[i-1].DayOfWeek {
					t.Errorf("hours not ordered: day %d came after day %d", business.Hours[i].DayOfWeek, business.Hours[i-1].DayOfWeek)
				}
			}

			if len(business.Menus) != tt.wantMenus {
				t.Errorf("got %d menus, want %d", len(business.Menus), tt.wantMenus)
			}

			if tt.wantMenus > 0 && len(business.Menus[0].Items) != tt.wantItems {
				t.Errorf("got %d items in first menu, want %d", len(business.Menus[0].Items), tt.wantItems)
			}
		})
	}
}

func TestListBusinesses(t *testing.T) {
	const setup = `
	INSERT INTO business_categories (name, slug) VALUES ('Restaurant', 'restaurant') ON CONFLICT DO NOTHING;
	INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
		VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Sooke Harbour House', 'sooke-harbour-house', '1528 Whiffen Spit Rd', 48.3538, -123.7256);
	INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
		VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Moms Cafe', 'moms-cafe', '2036 Shields Rd', 48.3761, -123.7254);
`

	tests := []struct {
		name      string
		search    string
		category  string
		limit     int
		offset    int
		wantCount int
		wantTotal int
	}{
		{
			name:      "no filters returns all businesses",
			search:    "",
			category:  "",
			limit:     20,
			offset:    0,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "search by name returns matching business",
			search:    "Harbour",
			category:  "",
			limit:     20,
			offset:    0,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "search case-insensitive",
			search:    "moms",
			category:  "",
			limit:     20,
			offset:    0,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "filter by category",
			search:    "",
			category:  "restaurant",
			limit:     20,
			offset:    0,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "search and filter together",
			search:    "Cafe",
			category:  "restaurant",
			limit:     20,
			offset:    0,
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "non-matching search returns no businesses",
			search:    "Nonexistent",
			category:  "",
			limit:     20,
			offset:    0,
			wantCount: 0,
			wantTotal: 0,
		},
		{
			name:      "pagination page 1 with 2 per page",
			search:    "",
			category:  "",
			limit:     2,
			offset:    0,
			wantCount: 2,
			wantTotal: 2,
		},
		{
			name:      "pagination page 2 with 1 per page",
			search:    "",
			category:  "",
			limit:     1,
			offset:    1,
			wantCount: 1,
			wantTotal: 2,
		},
		{
			name:      "pagination with invalid offset returns no businesses",
			search:    "",
			category:  "",
			limit:     20,
			offset:    100,
			wantCount: 0,
			wantTotal: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(setup); err != nil {
				t.Fatalf("setup failed: %v", err)
			}

			businesses, total, err := repository.ListBusinesses(context.Background(), tx, tt.search, tt.category, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("ListBusinesses returned error: %v", err)
			}

			if len(businesses) != tt.wantCount {
				t.Errorf("got %d businesses, want %d", len(businesses), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("got total %d, want %d", total, tt.wantTotal)
			}
		})
	}
}
