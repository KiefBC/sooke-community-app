package testdb

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Exec runs an HTTP request through the provided handler and returns the
// recorder. It must handle both call shapes used across this codebase.
// Both *chi.Mux and http.HandlerFunc satisfy http.Handler, so a single
// parameter type covers both cases without branching.
func Exec(t *testing.T, h http.Handler, method, target string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	return rec
}
