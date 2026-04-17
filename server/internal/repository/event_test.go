package repository_test

import (
	"context"
	"database/sql"
	"slices"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
	"github.com/kiefbc/sooke_app/server/internal/testdb/seeds"
)

func TestListEvents(t *testing.T) {
	tests := []struct {
		name       string
		search     string
		eventTypes []string
		limit      int
		offset     int
		wantCount  int
		wantTotal  int
		checkFunc  func(*testing.T, []repository.Event)
	}{
		{
			name:      "returns only approved events",
			limit:     20,
			wantCount: 4,
			wantTotal: 4,
		},
		{
			name:      "search by name",
			search:    "Jazz",
			limit:     20,
			wantCount: 1,
			wantTotal: 1,
			checkFunc: func(t *testing.T, events []repository.Event) {
				if events[0].Slug != "friday-night-jazz" {
					t.Errorf("expected friday-night-jazz, got %q", events[0].Slug)
				}
			},
		},
		{
			name:   "search by non-matching name",
			search: "Nonexistent Event",
			limit:  20,
		},
		{
			name:       "filter by single event type",
			eventTypes: []string{"live-music"},
			limit:      20,
			wantCount:  1,
			wantTotal:  1,
			checkFunc: func(t *testing.T, events []repository.Event) {
				if events[0].EventTypeSlug != "live-music" {
					t.Errorf("expected event type live-music, got %q", events[0].EventTypeSlug)
				}
			},
		},
		{
			name:       "filter by multiple event types matches OR",
			eventTypes: []string{"live-music", "market"},
			limit:      20,
			wantCount:  2,
			wantTotal:  2,
			checkFunc: func(t *testing.T, events []repository.Event) {
				for _, e := range events {
					if e.EventTypeSlug != "live-music" && e.EventTypeSlug != "market" {
						t.Errorf("unexpected event type %q in multi-filter result", e.EventTypeSlug)
					}
				}
			},
		},
		{
			name:      "pagination returns limited results but full total",
			limit:     1,
			wantCount: 1,
			wantTotal: 4,
		},
		{
			name:      "COALESCE resolves coordinates from both sources",
			limit:     20,
			wantCount: 4,
			wantTotal: 4,
			checkFunc: func(t *testing.T, events []repository.Event) {
				for _, e := range events {
					if e.Latitude == nil || e.Longitude == nil {
						t.Errorf("event %q should have coordinates (via own or COALESCE)", e.Slug)
					}
					if e.Slug == "whiffin-spit-beach-cleanup" {
						if e.BusinessName != nil {
							t.Error("standalone event should not have a business name")
						}
					} else {
						if e.BusinessName == nil {
							t.Errorf("business-linked event %q should have a business name", e.Slug)
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.EventSeed)

			events, total, err := repository.ListEvents(context.Background(), tx, tt.search, tt.eventTypes, tt.limit, tt.offset)
			if err != nil {
				t.Fatalf("ListEvents returned an error: %v", err)
			}

			if len(events) != tt.wantCount {
				t.Errorf("ListEvents returned %d events, want %d", len(events), tt.wantCount)
			}
			if total != tt.wantTotal {
				t.Errorf("ListEvents total = %d, want %d", total, tt.wantTotal)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, events)
			}
		})
	}
}

func TestGetEventBySlug(t *testing.T) {
	tests := []struct {
		name          string
		slug          string
		wantEventName string
		status        string
	}{
		{
			name:          "existing event",
			slug:          "friday-night-jazz",
			wantEventName: "Friday Night Jazz",
			status:        "approved",
		},
		{
			name:          "nonexistent event",
			slug:          "nonexistent-event",
			wantEventName: "",
			status:        "",
		},
		{
			name:          "pending review event",
			slug:          "cafe-acoustic-night",
			wantEventName: "Cafe Acoustic Night",
			status:        "pending_review",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.EventSeed)

			event, err := repository.GetEventBySlug(context.Background(), tx, tt.slug)
			if err != nil {
				t.Fatalf("GetEventBySlug returned an error: %v", err)
			}

			if tt.wantEventName == "" {
				if event != nil {
					t.Errorf("expected nil for nonexistent slug, got event with name %q", event.Name)
				}
			} else {
				if event == nil {
					t.Fatalf("expected event with name %q, got nil", tt.wantEventName)
				}
				if event.Name != tt.wantEventName {
					t.Errorf("expected event name %q, got %q", tt.wantEventName, event.Name)
				}
			}

			if event != nil && event.Status != tt.status {
				t.Errorf("expected event status %q, got %q", tt.status, event.Status)
			}
		})
	}
}

func TestGetEventTypes(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
		wantSlugs []string
		seed      func(*sql.Tx)
	}{
		{
			name:      "returns all event types",
			wantCount: 3,
			wantSlugs: []string{"live-music", "market", "community-meeting"},
			seed:      seeds.EventSeed,
		},
		{
			name:      "no event types",
			wantCount: 0,
			wantSlugs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tx *sql.Tx
			if tt.seed != nil {
				tx = testdb.WithTx(t, tt.seed)
			} else {
				tx = testdb.WithTx(t)
			}

			eventTypes, total, err := repository.ListEventTypes(context.Background(), tx)
			if err != nil {
				t.Fatalf("ListEventTypes returned an error: %v", err)
			}

			if len(eventTypes) != tt.wantCount {
				t.Errorf("ListEventTypes returned %d event types, want %d", len(eventTypes), tt.wantCount)
			}
			if total != tt.wantCount {
				t.Errorf("ListEventTypes total = %d, want %d", total, tt.wantCount)
			}

			for i, et := range eventTypes {
				if slices.Contains(tt.wantSlugs, et.Slug) == false {
					t.Errorf("unexpected event type slug %q at index %d", et.Slug, i)
				}
			}
		})
	}
}
