package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kiefbc/sooke_app/server/internal/handler"
)

// New creates a new chi router with the defined routes and middleware.
func New(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.HealthHandler(db))
		r.Get("/businesses", handler.ListBusinessesHandler(db))
		r.Get("/businesses/{slug}", handler.GetBusinessHandler(db))
		r.Get("/categories", handler.ListCategoriesHandler(db))
	})

	return r
}
