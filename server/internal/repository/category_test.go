package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/kiefbc/sooke_app/server/internal/repository"
)

func TestListCategories(t *testing.T) {
	const setup = `
		INSERT INTO business_categories (id, name, slug) VALUES
		(1, 'Restaurant', 'restaurant'),
		(2, 'Cafe', 'cafe'),
		(3, 'Retail', 'retail');
	`
	tests := []struct {
		name      string
		wantCount int
	}{
		{
			name:      "returns all categories",
			wantCount: 3,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx, err := testDB.Begin()
			if err != nil {
				t.Fatalf("failed to begin transaction: %v", err)
			}
			defer tx.Rollback()

			if _, err := tx.Exec(setup); err != nil {
				t.Fatalf("failed to set up test data: %v", err)
			}

			categories, err := repository.ListCategories(ctx, tx)
			if err != nil {
				t.Fatalf("ListCategories returned an error: %v", err)
			}

			if len(categories) != tt.wantCount {
				t.Errorf("expected %d categories, got %d", tt.wantCount, len(categories))	
			}
		})
	}
}
