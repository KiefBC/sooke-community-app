package handler

import (
	"context"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kiefbc/sooke_app/server/internal/repository"
)

func GetBusinessHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
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

func ListBusinessesHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		search := r.URL.Query().Get("search")
		category := r.URL.Query().Get("category")

		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
		if perPage < 1 || perPage > 100 {
			perPage = 20
		}

		offset := (page - 1) * perPage

		businesses, total, err := repository.ListBusinesses(ctx, db, search, category, perPage, offset)
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
