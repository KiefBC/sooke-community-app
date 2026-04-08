package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// EventType represents an event type entity in the database.
type EventType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // Name of the event type
	Slug string `json:"slug"` // URL-friendly identifier derived from the name
}

// Event represents an event entity in the database.
type Event struct {
	ID            int64      `json:"id"`
	Name          string     `json:"name"`            // Name of the event
	Slug          string     `json:"slug"`            // URL-friendly identifier derived from the name
	Description   *string    `json:"description"`     // Description of the event, nullable
	EventTypeName string     `json:"event_type_name"` // Name of the event type this event belongs to
	EventTypeSlug string     `json:"event_type_slug"` // Slug of the event type this event belongs to
	Latitude      *float64   `json:"latitude"`        // Geographic latitude of the event location, nullable (if the event does not have a specified location)
	Longitude     *float64   `json:"longitude"`       // Geographic longitude of the event location, nullable (if the event does not have a specified location)
	StartsAt      time.Time  `json:"starts_at"`       // Start time of the event
	EndsAt        *time.Time `json:"ends_at"`         // End time of the event, nullable (if the event has no specified end time)
	Status        string     `json:"status"`          // Status of the event (e.g., "upcoming", "ongoing", "past")
	BusinessName  *string    `json:"business_name"`   // Name of the business hosting the event, nullable (if the event is not associated with a specific business)
	BusinessSlug  *string    `json:"business_slug"`   // Slug of the business hosting the event, nullable (if the event is not associated with a specific business)
}

// ListEvents retrieves all events from the database, optionally filtered by search term and event type, with pagination.
func ListEvents(ctx context.Context, q Querier, search, eventTypeSlug string, limit, offset int) ([]Event, int, error) {
	var countTotal int
	err := q.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM events e
	JOIN event_types et ON e.event_type_id = et.id
	WHERE ($1 = '' OR e.name ILIKE '%' || $1 || '%')
	AND ($2 = '' OR et.slug = $2)
	AND e.status IN ('upcoming', 'ongoing', 'approved')
	AND e.starts_at >= NOW()
	`, search, eventTypeSlug).Scan(&countTotal)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}

	rows, err := q.QueryContext(ctx, `
	SELECT e.id, e.name, e.slug, e.description, et.name, et.slug, COALESCE(e.latitude, b.latitude) AS latitude, COALESCE(e.longitude, b.longitude) AS longitude, e.starts_at, e.ends_at, e.status, b.name, b.slug
	FROM events e
	JOIN event_types et ON e.event_type_id = et.id
	LEFT JOIN businesses b ON e.business_id = b.id
	WHERE ($1 = '' OR e.name ILIKE '%' || $1 || '%')
	AND ($2 = '' OR et.slug = $2)
	AND e.status IN ('upcoming', 'ongoing', 'approved')
	AND e.starts_at >= NOW()
	ORDER BY e.starts_at ASC
	LIMIT $3 OFFSET $4
	`, search, eventTypeSlug, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count events: %w", err)
	}
	defer rows.Close()

	events := []Event{}
	for rows.Next() {
		var e Event
		err := rows.Scan(
			&e.ID, &e.Name, &e.Slug, &e.Description,
			&e.EventTypeName, &e.EventTypeSlug,
			&e.Latitude, &e.Longitude,
			&e.StartsAt, &e.EndsAt, &e.Status,
			&e.BusinessName, &e.BusinessSlug,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
	}

	return events, countTotal, nil
}

// / GetEventBySlug retrieves a single event from the database based on its slug. It returns nil if no event is found with the given slug.
func GetEventBySlug(ctx context.Context, q Querier, slug string) (*Event, error) {
	var e Event
	err := q.QueryRowContext(ctx, `
	SELECT e.id, e.name, e.slug, e.description, et.name, et.slug, COALESCE(e.latitude, b.latitude) AS latitude, COALESCE(e.longitude, b.longitude) AS longitude, e.starts_at, e.ends_at, e.status, b.name, b.slug
	FROM events e
	JOIN event_types et ON e.event_type_id = et.id
	LEFT JOIN businesses b ON e.business_id = b.id
	WHERE e.slug = $1
	ORDER BY e.starts_at ASC
	`, slug).Scan(&e.ID, &e.Name, &e.Slug, &e.Description, &e.EventTypeName, &e.EventTypeSlug, &e.Latitude, &e.Longitude, &e.StartsAt, &e.EndsAt, &e.Status, &e.BusinessName, &e.BusinessSlug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found with the given slug
		}
		return nil, fmt.Errorf("failed to get event by slug: %w", err)
	}

	return &e, nil
}

// ListEventTypes retrieves all event types from the database, ordered by name.
func ListEventTypes(ctx context.Context, q Querier) ([]EventType, int, error) {
	var countTotal int
	err := q.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM event_types
	`).Scan(&countTotal)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count event types: %w", err)
	}

	rows, err := q.QueryContext(ctx, `
	SELECT id, name, slug
	FROM event_types
	ORDER BY name ASC
	`)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query event types: %w", err)
	}
	defer rows.Close()

	eventTypes := []EventType{}
	for rows.Next() {
		var et EventType
		err := rows.Scan(&et.ID, &et.Name, &et.Slug)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan event type: %w", err)
		}
		eventTypes = append(eventTypes, et)
	}

	return eventTypes, countTotal, nil
}
