package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/router"
)

func main() {
	// godotenv.Load will not return an error if the file is missing, so we can ignore it
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("./.env")

	port := os.Getenv("PORT")
	if port == "" {
		// I forget to set the PORT variable all the time, so let's default to 8080 if it's not set
		port = "8080"
	}

	migrationPath := os.Getenv("MIGRATION_PATH")
	if migrationPath == "" {
		// Hardcoding this path is not ideal, but it makes it easier to run the server without having to set the MIGRATION_PATH variable every time
		migrationPath = "./migrations"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection established...")

	if err := database.Migrate(db, migrationPath); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	log.Printf("Database migrations completed")

	r := router.New(db)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	log.Printf("Starting server on port %s...", port)
	log.Printf("Server: http://127.0.0.1:%s/", port)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed to start: %v", err)
	}
}
