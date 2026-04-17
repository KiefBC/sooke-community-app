package handler

import (
	"context"
	"math"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kiefbc/sooke_app/server/internal/repository"
)

// GetBusinessHandler retrieves a single business by its slug. It returns a 404 if the business is not found.
func GetBusinessHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
		defer cancel()

		slug := chi.URLParam(r, "slug")

		business, err := repository.GetBusinessBySlug(ctx, db, slug)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to retrieve business")
			return
		}
		if business == nil {
			WriteError(w, http.StatusNotFound, "not_found", "Business not found")
			return
		}

		WriteJSON(w, http.StatusOK, business)
	}
}

// ListBusinessesHandler retrieves a list of businesses based on search and category filters, along with pagination. It returns a paginated response containing the list of businesses and pagination metadata.
func ListBusinessesHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), TIMEOUT)
		defer cancel()

		search := r.URL.Query().Get("search")
		category := r.URL.Query().Get("category")

		timeZone, err := TimeZoneHelper(r)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "invalid_parameter", err.Error())
			return
		}
		page, perPage, offset := PaginationHelper(r)

		businesses, total, err := repository.ListBusinesses(ctx, db, search, category, timeZone, perPage, offset)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list businesses")
			return
		}

		totalPages := int(math.Ceil(float64(total) / float64(perPage)))

		WriteJSON(w, http.StatusOK, PaginatedResponse[repository.Business]{
			Items: businesses,
			Pagination: Pagination{
				Page:       page,
				PerPage:    perPage,
				TotalItems: total,
				TotalPages: totalPages,
			},
		})
	}
}
