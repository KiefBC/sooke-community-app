package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/handler"
	"github.com/kiefbc/sooke_app/server/internal/repository"
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
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.EventSeed(tx)

			h := handler.GetEventHandler(tx)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			rec := httptest.NewRecorder()
			h(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if tt.wantSlug != "" {
				var event repository.Event
				if err := json.NewDecoder(rec.Body).Decode(&event); err != nil {
					t.Fatalf("failed to decode response: %v", err)
				}

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
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.EventSeed(tx)

			h := handler.ListEventTypesHandler(tx)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/event-types", nil)
			rec := httptest.NewRecorder()
			h(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if ct := rec.Header().Get("Content-Type"); ct != tt.wantContentType {
				t.Errorf("Content-Type = %q, want %q", ct, tt.wantContentType)
			}

			var body handler.ListResponse[repository.EventType]
			if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if len(body.Items) < tt.wantAmount {
				t.Fatalf("got %d event types, want at least %d", len(body.Items), tt.wantAmount)
			}
		})
	}
}
