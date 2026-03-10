package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// HealthResponse represents the JSON response for the health check endpoint.
type HealthResponse struct {
	Status   string `json:"status"`
	DBStatus string `json:"db_status"`
}

// HealthHandler returns an HTTP handler function that checks the health of the application and its database connection.
func HealthHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		dbStatus := "connected"
		if db == nil {
			dbStatus = "disconnected"
		} else if err := db.PingContext(ctx); err != nil {
			dbStatus = "disconnected"
		}

		status := "ok"
		httpStatus := http.StatusOK
		if dbStatus != "connected" {
			status = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(HealthResponse{
			Status:   status,
			DBStatus: dbStatus,
		}); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		if _, err := w.Write(buf.Bytes()); err != nil {
			log.Printf("health handler: failed to write response: %v", err)
		}
	}
}
