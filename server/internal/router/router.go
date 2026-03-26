package router

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/kiefbc/sooke_app/server/internal/handler"
)

// New creates a new chi router with the defined routes and middleware.
func New(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handler.HealthHandler(db))
		r.Get("/businesses", handler.ListBusinessesHandler(db))
		r.Get("/businesses/{slug}", handler.GetBusinessHandler(db))
		r.Get("/categories", handler.ListCategoriesHandler(db))
	})

	return r
}
