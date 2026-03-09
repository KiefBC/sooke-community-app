package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kiefbc/sooke_app/server/internal/router"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := router.New()

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
