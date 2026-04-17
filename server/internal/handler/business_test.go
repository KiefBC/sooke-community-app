package handler_test

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
	"github.com/kiefbc/sooke_app/server/internal/testdb/seeds"
)

func TestTimeZoneValidation(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		wantStatus  int
		wantErrCode string
	}{
		{
			name:        "valid time zone returns 200",
			url:         "/api/v1/businesses?tz=America%2FVancouver",
			wantStatus:  http.StatusOK,
			wantErrCode: "",
		},
		{
			name:        "invalid time zone returns 400 with error JSON",
			url:         "/api/v1/businesses?tz=Invalid%2FZone",
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "invalid_parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t)

			rec := testdb.Exec(t, handler.ListBusinessesHandler(tx), http.MethodGet, tt.url, nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if tt.wantErrCode != "" {
				var body handler.ErrorResponse
				testdb.DecodeJSON(t, rec, &body)
				if body.Error.Code != tt.wantErrCode {
					t.Errorf("error code = %q, want %q", body.Error.Code, tt.wantErrCode)
				}
				if body.Error.Message == "" {
					t.Error("error message is empty")
				}
			}
		})
	}
}

func TestListBusinesses(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		wantStatus      int
		wantMinItems    int
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
		{
			name:            "today_hours included for business with hours",
			url:             "/api/v1/businesses?search=harbour&tz=America%2FVancouver",
			wantStatus:      http.StatusOK,
			wantMinItems:    1,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.BusinessSeed)

			rec := testdb.Exec(t, handler.ListBusinessesHandler(tx), http.MethodGet, tt.url, nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body handler.PaginatedResponse[repository.Business]
			testdb.DecodeJSON(t, rec, &body)

			if len(body.Items) < tt.wantMinItems {
				t.Errorf("got %d items, want at least %d", len(body.Items), tt.wantMinItems)
			}

			for _, b := range body.Items {
				if b.Name == "" || b.Slug == "" || b.Address == "" {
					t.Errorf("required fields missing: name=%q slug=%q address=%q", b.Name, b.Slug, b.Address)
				}
			}

			if tt.name == "today_hours included for business with hours" {
				for _, b := range body.Items {
					if b.Slug == "sooke-harbour-house" && b.TodayHours == nil {
						t.Error("expected today_hours for sooke-harbour-house, got nil")
					}
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
			tx := testdb.WithTx(t, seeds.BusinessSeed)

			r := chi.NewRouter()
			r.Get("/businesses/{slug}", handler.GetBusinessHandler(tx))

			rec := testdb.Exec(t, r, http.MethodGet, "/businesses/"+tt.slug, nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", ct)
			}

			if tt.wantErrCode != "" {
				var body handler.ErrorResponse
				testdb.DecodeJSON(t, rec, &body)
				if body.Error.Code != tt.wantErrCode {
					t.Errorf("error code = %q, want %q", body.Error.Code, tt.wantErrCode)
				}
				if body.Error.Message == "" {
					t.Error("error message is empty")
				}
				return
			}

			var body repository.BusinessDetails
			testdb.DecodeJSON(t, rec, &body)

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
			tx := testdb.WithTx(t, seeds.CategorySeed)

			rec := testdb.Exec(t, handler.ListCategoriesHandler(tx), http.MethodGet, "/api/v1/categories", nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body handler.ListResponse[repository.Category]
			testdb.DecodeJSON(t, rec, &body)

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
