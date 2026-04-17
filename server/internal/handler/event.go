package handler

import (
	"context"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiefbc/sooke_app/server/internal/repository"
)

// ListEventsHandler retrieves a paginated list of events, optionally filtered by search term and category. It returns a paginated response with the total number of events and total pages.
func ListEventsHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
		defer cancel()

		search := r.URL.Query().Get("search")
		category := r.URL.Query().Get("category")
		page, perPage, offset := PaginationHelper(r)

		events, total, err := repository.ListEvents(ctx, db, search, category, perPage, offset)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list events")
			return
		}

		totalPages := int(math.Ceil(float64(total) / float64(perPage)))

		WriteJSON(w, http.StatusOK, PaginatedResponse[repository.Event]{
			Items: events,
			Pagination: Pagination{
				Page:       page,
				PerPage:    perPage,
				TotalItems: total,
				TotalPages: totalPages,
			},
		})
	}
}

// GetEventHandler retrieves a single event by its slug. It returns a 404 if the event is not found.
func GetEventHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
		defer cancel()

		slug := chi.URLParam(r, "slug")

		event, err := repository.GetEventBySlug(ctx, db, slug)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve event")
			return
		}
		if event == nil {
			WriteError(w, http.StatusNotFound, "not_found", "Event not found")
			return
		}

		WriteJSON(w, http.StatusOK, event)
	}
}

func ListEventTypesHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
		defer cancel()

		eventTypes, _, err := repository.ListEventTypes(ctx, db)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list event types")
			return
		}

		WriteJSON(w, http.StatusOK, ListResponse[repository.EventType]{Items: eventTypes})
	}
}
