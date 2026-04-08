package repository_test

import (
	"context"
	"slices"
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb/seeds"
)

func TestListEvents(t *testing.T) {
	tests := []struct {
		name          string
		search        string
		eventTypeSlug string
		limit         int
		offset        int
		wantCount     int
		wantTotal     int
		checkFunc     func(*testing.T, []repository.Event)
	}{
		{
			name:          "returns only approved events",
			search:        "",
			eventTypeSlug: "",
			limit:         20,
			offset:        0,
			wantCount:     4,
			wantTotal:     4,
		},
		{
			name:          "search by name",
			search:        "Jazz",
			eventTypeSlug: "",
			limit:         20,
			offset:        0,
			wantCount:     1,
			wantTotal:     1,
			checkFunc: func(t *testing.T, events []repository.Event) {
				if events[0].Slug != "friday-night-jazz" {
					t.Errorf("expected friday-night-jazz, got %q", events[0].Slug)
				}
			},
		},
		{
			name:          "search by non-matching name",
			search:        "Nonexistent Event",
			eventTypeSlug: "",
			limit:         20,
			offset:        0,
			wantCount:     0,
			wantTotal:     0,
		},
		{
			name:          "filter by event type",
			search:        "",
			eventTypeSlug: "live-music",
			limit:         20,
			offset:        0,
			wantCount:     1,
			wantTotal:     1,
			checkFunc: func(t *testing.T, events []repository.Event) {
				if events[0].EventTypeSlug != "live-music" {
					t.Errorf("expected event type live-music, got %q", events[0].EventTypeSlug)
				}
			},
		},
		{
			name:          "pagination returns limited results but full total",
			search:        "",
			eventTypeSlug: "",
			limit:         1,
			offset:        0,
			wantCount:     1,
			wantTotal:     4,
		},
		{
			name:          "COALESCE resolves coordinates from both sources",
			search:        "",
			eventTypeSlug: "",
			limit:         20,
			offset:        0,
			wantCount:     4,
			wantTotal:     4,
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
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			seeds.EventSeed(tx)

			events, total, err := repository.ListEvents(context.Background(), tx, tt.search, tt.eventTypeSlug, tt.limit, tt.offset)
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
		name string
		slug string
		wantEventName string
		status string
	}{
		{
			name: "existing event",
			slug: "friday-night-jazz",
			wantEventName: "Friday Night Jazz",
			status: "approved",
		},
		{
			name: "nonexistent event",
			slug: "nonexistent-event",
			wantEventName: "",
			status: "",
		},
		{
			name: "pending review event",
			slug: "cafe-acoustic-night",
			wantEventName: "Cafe Acoustic Night",
			status: "pending_review",
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
		name string
		wantCount int
		wantSlugs []string
	}{
		{
			name: "returns all event types",
			wantCount: 3,
			wantSlugs: []string{"live-music", "market", "community-meeting"},
		},
		{
			name: "no event types",
			wantCount: 0,
			wantSlugs: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if tt.wantCount > 0 {
				seeds.EventSeed(tx)
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
