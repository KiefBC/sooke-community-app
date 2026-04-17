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

func TestGetEvent(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
		wantSlug   string
		url        string
	}{
		{
			name:       "returns event by slug",
			wantStatus: http.StatusOK,
			wantSlug:   "friday-night-jazz",
			url:        "/api/v1/events/friday-night-jazz",
		},
		{
			name:       "returns 404 for non-existent slug",
			wantStatus: http.StatusNotFound,
			wantSlug:   "",
			url:        "/api/v1/events/nonexistent-event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.EventSeed)

			r := chi.NewRouter()
			r.Get("/api/v1/events/{slug}", handler.GetEventHandler(tx))

			rec := testdb.Exec(t, r, http.MethodGet, tt.url, nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if tt.wantSlug != "" {
				var event repository.Event
				testdb.DecodeJSON(t, rec, &event)
				if event.Slug != tt.wantSlug {
					t.Errorf("event slug = %q, want %q", event.Slug, tt.wantSlug)
				}
			}
		})
	}
}

func TestListEventTypes(t *testing.T) {
	tests := []struct {
		name            string
		wantStatus      int
		wantAmount      int
		wantContentType string
	}{
		{
			name:            "returns seeded event types",
			wantStatus:      http.StatusOK,
			wantAmount:      3,
			wantContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.EventSeed)

			rec := testdb.Exec(t, handler.ListEventTypesHandler(tx), http.MethodGet, "/api/v1/event-types", nil)
			testdb.AssertStatus(t, rec, tt.wantStatus)

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body handler.ListResponse[repository.EventType]
			testdb.DecodeJSON(t, rec, &body)
			if len(body.Items) < tt.wantAmount {
				t.Fatalf("got %d event types, want at least %d", len(body.Items), tt.wantAmount)
			}
		})
	}
}
