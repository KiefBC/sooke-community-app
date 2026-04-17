package testdb

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

// AssertStatus fails the test when rec.Code does not match want. The
// response body is included in the failure message so a 500 with a JSON
// error payload still tells you what went wrong.
func AssertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, want, rec.Body.String())
	}
}

// DecodeJSON decodes rec.Body into target or fails the test. target must
// be a pointer (same rules as json.Unmarshal).
func DecodeJSON(t *testing.T, rec *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(rec.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v, body = %s", err, rec.Body.String())
	}
}
