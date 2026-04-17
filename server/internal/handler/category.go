package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/kiefbc/sooke_app/server/internal/repository"
)

func ListCategoriesHandler(db repository.Querier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		categories, err := repository.ListCategories(ctx, db)
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to list categories")
			return
		}

		WriteJSON(w, http.StatusOK, ListResponse[repository.Category]{Items: categories})
	}
}
