package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// PaginatedResponse represents the structure of a paginated response, containing a list of items and pagination information. It is a generic type that can be used for any type of items.
type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

// Pagination represents the structure of pagination information, including the current page, items per page, total items, and total pages.
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

// ErrorResponse represents the structure of an error response, containing an error detail with a code and a message.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the structure of an error response, containing a code and a message.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteJSON writes a JSON response to the http.ResponseWriter with the specified status code and data. If encoding the data fails, it writes an internal server error response.
func WriteJSON(w http.ResponseWriter, statusCode int, data any) {

	// This was a subtle bug, if this is defined after w.WriteHeader and the encoding fails, it gets sent to WriteError and the second call to WriteHeader in WriteError will cause an error because the headers have already been sent. By defining the buffer before writing the header, we ensure that any encoding errors are handled properly before sending the response.
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to encode response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}

// WriteError writes an error response to the http.ResponseWriter with the specified status code, error code, and message. It constructs an ErrorResponse and encodes it as JSON. If encoding fails, it logs the error.
func WriteError(w http.ResponseWriter, statusCode int, code string, message string) {
	errorResponse := ErrorResponse{
		Error: ErrorDetail{
			code,
			message,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(errorResponse); err != nil {
		fmt.Printf("Failed to encode error response: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
