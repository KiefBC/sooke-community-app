package slug_test

import (
	"testing"

	"github.com/kiefbc/sooke_app/server/internal/slug"
)

func TestSlug(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		input    string
	}{
		{
			name:     "Basic name",
			expected: "john-doe",
			input:    "John Doe",
		},
		{
			name:     "Name with special characters",
			expected: "johns-bake-shop",
			input:    "John's Bake Shop!",
		},
		{
			name:     "Name with multiple spaces",
			expected: "the-great-outdoors",
			input:    "The    Great  Outdoors",
		},
		{
			name:     "Name spelled in numbers",
			expected: "mile-42-diner",
			input:    "-Mile 42 Diner-",
		},
		{
			name:     "Already slugified name",
			expected: "already-slugified",
			input:    "already-slugified",
		},
		{
			name:     "Empty name",
			expected: "",
			input:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := slug.GenerateSlug(tt.input)

			if slug != tt.expected {
				t.Errorf("expected slug '%s', got '%s'", tt.expected, slug)
			}
		})
	}
}
