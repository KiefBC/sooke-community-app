package handler_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	testDB = testdb.Open()
	if testDB == nil {
		os.Exit(0)
	}
	os.Exit(m.Run())
}

const businessSeed = `
	INSERT INTO business_categories (name, slug) VALUES ('Restaurant', 'restaurant'), ('Cafe', 'cafe') ON CONFLICT DO NOTHING;
	INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
		VALUES ((SELECT id FROM business_categories WHERE slug = 'restaurant'), 'Sooke Harbour House', 'sooke-harbour-house', '1528 Whiffen Spit Rd', 48.3538, -123.7256);
	INSERT INTO businesses (category_id, name, slug, address, latitude, longitude)
		VALUES ((SELECT id FROM business_categories WHERE slug = 'cafe'), 'Moms Cafe', 'moms-cafe', '2036 Shields Rd', 48.3761, -123.7254);
	INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
		VALUES ((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 1, '09:00', '17:00', false);
	INSERT INTO menus (business_id, name) VALUES ((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 'Dinner');
	INSERT INTO menu_items (menu_id, name, price) VALUES ((SELECT id FROM menus WHERE name = 'Dinner'), 'Fish and Chips', '12.99');
`

func TestListBusinesses(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		wantStatus     int
		wantMinItems   int
		wantContentType string
	}{
		{
			name:            "200 with items and pagination",
			url:             "/api/v1/businesses",
			wantStatus:      http.StatusOK,
			wantMinItems:    2,
			wantContentType: "application/json",
		},
		{
			name:            "search returns filtered results",
			url:             "/api/v1/businesses?search=harbour",
			wantStatus:      http.StatusOK,
			wantMinItems:    1,
			wantContentType: "application/json",
		},
		{
			name:            "category filter works",
			url:             "/api/v1/businesses?category=cafe",
			wantStatus:      http.StatusOK,
			wantMinItems:    1,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(businessSeed); err != nil {
				t.Fatalf("seed failed: %v", err)
			}

			h := handler.ListBusinessesHandler(tx)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			h(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body handler.PaginatedResponse[repository.Business]
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(body.Items) < tt.wantMinItems {
				t.Errorf("got %d items, want at least %d", len(body.Items), tt.wantMinItems)
			}

			for _, b := range body.Items {
				if b.Name == "" || b.Slug == "" || b.Address == "" {
					t.Errorf("required fields missing: name=%q slug=%q address=%q", b.Name, b.Slug, b.Address)
				}
			}

			if body.Pagination.Page < 1 {
				t.Errorf("pagination.page = %d, want >= 1", body.Pagination.Page)
			}
			if body.Pagination.PerPage < 1 {
				t.Errorf("pagination.per_page = %d, want >= 1", body.Pagination.PerPage)
			}
			if body.Pagination.TotalItems < tt.wantMinItems {
				t.Errorf("pagination.total_items = %d, want >= %d", body.Pagination.TotalItems, tt.wantMinItems)
			}
			if body.Pagination.TotalPages < 1 {
				t.Errorf("pagination.total_pages = %d, want >= 1", body.Pagination.TotalPages)
			}
		})
	}
}

func TestGetBusiness(t *testing.T) {
	tests := []struct {
		name         string
		slug         string
		wantStatus   int
		wantName     string
		wantErrCode  string
		wantMinHours int
		wantMinMenus int
	}{
		{
			name:         "known slug returns 200 with full detail",
			slug:         "sooke-harbour-house",
			wantStatus:   http.StatusOK,
			wantName:     "Sooke Harbour House",
			wantMinHours: 1,
			wantMinMenus: 1,
		},
		{
			name:        "unknown slug returns 404 with error JSON",
			slug:        "nonexistent-slug",
			wantStatus:  http.StatusNotFound,
			wantErrCode: "not_found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(businessSeed); err != nil {
				t.Fatalf("seed failed: %v", err)
			}

			h := handler.GetBusinessHandler(tx)

			r := chi.NewRouter()
			r.Get("/businesses/{slug}", h)

			req := httptest.NewRequest(http.MethodGet, "/businesses/"+tt.slug, nil)
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tt.wantStatus, rec.Body.String())
			}

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			if tt.wantErrCode != "" {
				var body handler.ErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if body.Error.Code != tt.wantErrCode {
					t.Errorf("error code = %q, want %q", body.Error.Code, tt.wantErrCode)
				}
				if body.Error.Message == "" {
					t.Error("error message is empty")
				}
				return
			}

			var body repository.BusinessDetails
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if body.Name != tt.wantName {
				t.Errorf("name = %q, want %q", body.Name, tt.wantName)
			}
			if body.Slug == "" || body.Address == "" {
				t.Errorf("required fields missing: slug=%q address=%q", body.Slug, body.Address)
			}
			if len(body.Hours) < tt.wantMinHours {
				t.Errorf("hours count = %d, want >= %d", len(body.Hours), tt.wantMinHours)
			}
			if len(body.Menus) < tt.wantMinMenus {
				t.Errorf("menus count = %d, want >= %d", len(body.Menus), tt.wantMinMenus)
			}
			if len(body.Menus) > 0 && len(body.Menus[0].Items) < 1 {
				t.Errorf("menu items count = %d, want >= 1", len(body.Menus[0].Items))
			}
		})
	}
}

func TestGetCategories(t *testing.T) {
	const categorySeed = `
		INSERT INTO business_categories (name, slug) VALUES ('Restaurant', 'restaurant'), ('Cafe', 'cafe') ON CONFLICT DO NOTHING;
	`

	tests := []struct {
		name            string
		wantStatus      int
		wantMinItems    int
		wantContentType string
	}{
		{
			name:            "returns sorted list",
			wantStatus:      http.StatusOK,
			wantMinItems:    2,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(categorySeed); err != nil {
				t.Fatalf("seed failed: %v", err)
			}

			h := handler.ListCategoriesHandler(tx)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
			rec := httptest.NewRecorder()
			h(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body struct {
				Items []repository.Category `json:"items"`
			}
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(body.Items) < tt.wantMinItems {
				t.Fatalf("got %d categories, want at least %d", len(body.Items), tt.wantMinItems)
			}

			for i := 1; i < len(body.Items); i++ {
				if body.Items[i].Name < body.Items[i-1].Name {
					t.Errorf("categories not sorted: %q came after %q", body.Items[i].Name, body.Items[i-1].Name)
				}
			}
		})
	}
}
