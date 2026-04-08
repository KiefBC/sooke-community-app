package repository_test

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
	"github.com/kiefbc/sooke_app/server/internal/testdb/seeds"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDB = testdb.Open()
	if testDB == nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestGetBusinessBySlug(t *testing.T) {
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
			wantHours: 7,
			wantMenus: 1,
			wantItems: 3,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.BusinessSeed(tx)

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

func TestListBusinessesTodayHours(t *testing.T) {
	// Edge-case businesses layered on top of the master seed.
	// Uses EXTRACT(DOW FROM NOW()) so the inserted hours match "today" regardless of when tests run.
	const edgeCases = `
		INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
			VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Has Hours Today', 'has-hours-today', '1 Main St', 48.35, -123.72);
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
			VALUES ((SELECT id FROM businesses WHERE slug = 'has-hours-today'), EXTRACT(DOW FROM NOW())::int, '09:00', '17:00', false);
		INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
			VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'No Hours', 'no-hours', '2 Main St', 48.35, -123.72);
		INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
			VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Wrong Day Only', 'wrong-day-only', '3 Main St', 48.35, -123.72);
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
			VALUES ((SELECT id FROM businesses WHERE slug = 'wrong-day-only'), (EXTRACT(DOW FROM NOW())::int + 1) % 7, '10:00', '18:00', false);
		INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
			VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Closed Today', 'closed-today', '4 Main St', 48.35, -123.72);
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
			VALUES ((SELECT id FROM businesses WHERE slug = 'closed-today'), EXTRACT(DOW FROM NOW())::int, '09:00', '17:00', true);
	`

	tests := []struct {
		name         string
		search       string
		wantSlug     string
		wantHasHours bool
		wantIsClosed bool
		wantOpenTime string
	}{
		{
			name:         "business with hours today has TodayHours populated",
			search:       "Has Hours Today",
			wantSlug:     "has-hours-today",
			wantHasHours: true,
			wantIsClosed: false,
			wantOpenTime: "09:00:00",
		},
		{
			name:         "business without any hours has nil TodayHours",
			search:       "No Hours",
			wantSlug:     "no-hours",
			wantHasHours: false,
		},
		{
			name:         "business with hours for different day has nil TodayHours",
			search:       "Wrong Day Only",
			wantSlug:     "wrong-day-only",
			wantHasHours: false,
		},
		{
			name:         "business closed today has TodayHours with IsClosed true",
			search:       "Closed Today",
			wantSlug:     "closed-today",
			wantHasHours: true,
			wantIsClosed: true,
			wantOpenTime: "09:00:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.BusinessSeed(tx)
			if _, err := tx.Exec(edgeCases); err != nil {
				t.Fatalf("edge-case setup failed: %v", err)
			}

			businesses, _, err := repository.ListBusinesses(context.Background(), tx, tt.search, "", "America/Vancouver", 20, 0)
			if err != nil {
				t.Fatalf("ListBusinesses returned error: %v", err)
			}

			if len(businesses) != 1 {
				t.Fatalf("got %d businesses, want 1", len(businesses))
			}

			b := businesses[0]
			if b.Slug != tt.wantSlug {
				t.Errorf("got slug %q, want %q", b.Slug, tt.wantSlug)
			}

			if tt.wantHasHours {
				if b.TodayHours == nil {
					t.Fatal("got nil TodayHours, want non-nil")
				}
				if b.TodayHours.IsClosed != tt.wantIsClosed {
					t.Errorf("got IsClosed %v, want %v", b.TodayHours.IsClosed, tt.wantIsClosed)
				}
				if b.TodayHours.OpenTime != tt.wantOpenTime {
					t.Errorf("got OpenTime %q, want %q", b.TodayHours.OpenTime, tt.wantOpenTime)
				}
			} else {
				if b.TodayHours != nil {
					t.Errorf("got TodayHours %+v, want nil", b.TodayHours)
				}
			}
		})
	}
}

func TestListBusinesses(t *testing.T) {
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
			wantCount: 5,
			wantTotal: 5,
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
			wantCount: 1,
			wantTotal: 1,
		},
		{
			name:      "search and filter together",
			search:    "Cafe",
			category:  "cafe",
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
			name:      "pagination first page",
			search:    "",
			category:  "",
			limit:     3,
			offset:    0,
			wantCount: 3,
			wantTotal: 5,
		},
		{
			name:      "pagination second page",
			search:    "",
			category:  "",
			limit:     3,
			offset:    3,
			wantCount: 2,
			wantTotal: 5,
		},
		{
			name:      "pagination with invalid offset returns no businesses",
			search:    "",
			category:  "",
			limit:     20,
			offset:    100,
			wantCount: 0,
			wantTotal: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.BusinessSeed(tx)

			businesses, total, err := repository.ListBusinesses(context.Background(), tx, tt.search, tt.category, "America/Vancouver", tt.limit, tt.offset)
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
