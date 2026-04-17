package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	TIMEOUT = 5 * time.Second
)

// TimeZoneHelper extracts the time zone from the `tz` query parameter. When
// the parameter is omitted, it returns the default "America/Vancouver". When
// the parameter is provided but invalid (not a known IANA zone), it returns
// an error so the handler can surface a 400 to the client instead of
// silently substituting a different zone.
func TimeZoneHelper(r *http.Request) (string, error) {
	timeZone := r.URL.Query().Get("tz")
	if timeZone == "" {
		return "America/Vancouver", nil
	}

	if _, err := time.LoadLocation(timeZone); err != nil {
		return "", fmt.Errorf("invalid time zone: %q", timeZone)
	}

	return timeZone, nil
}

// PaginationHelper extracts pagination parameters from the request and calculates the offset for database queries.
func PaginationHelper(r *http.Request) (page int, perPage int, offset int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	offset = (page - 1) * perPage

	return page, perPage, offset
}
