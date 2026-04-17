package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/kiefbc/sooke_app/server/internal/repository"
	"github.com/kiefbc/sooke_app/server/internal/testdb"
	"github.com/kiefbc/sooke_app/server/internal/testdb/seeds"
)

func TestListCategories(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
	}{
		{
			name:      "returns all categories",
			wantCount: 5,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := testdb.WithTx(t, seeds.CategorySeed)

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
