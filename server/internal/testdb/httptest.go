package testdb

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Exec runs an HTTP request through the provided handler and returns the
// recorder. It must handle both call shapes used across this codebase:
//
//  1. A raw http.HandlerFunc - e.g. handler.GetEventHandler(tx) invoked
//     directly against a request with a hard-coded URL.
//  2. A chi router that needs to extract path params from the URL - e.g.
//     businesses/{slug}. The caller builds the router, registers the
//     handler on a route, and passes the router in.
//
// Both *chi.Mux and http.HandlerFunc satisfy http.Handler, so a single
// parameter type covers both cases without branching.
//
// TODO(human): implement Exec. Decide on the parameter shape (method,
// target URL, request body) and return *httptest.ResponseRecorder.
// Consider: should body be io.Reader (nil-friendly) or []byte? How do
// callers pass JSON bodies cleanly?
func Exec(t *testing.T, h http.Handler, method, target string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec

}
