package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// HealthResponse represents the JSON response for the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(HealthResponse{Status: "ok"}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("health handler: failed to write response: %v", err)
	}
}
