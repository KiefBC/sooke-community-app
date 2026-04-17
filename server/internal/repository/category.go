package repository

import (
	"context"
	"fmt"
)

// / Category represents a business category in the database.
type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // Name of the category (e.g., "Restaurant", "Retail", "Service")
	Slug string `json:"slug"` // URL-friendly identifier derived from the name (e.g., "restaurant", "retail", "service")
}

// / ListCategories retrieves all business categories from the database, ordered by name.
func ListCategories(ctx context.Context, q Querier) ([]Category, error) {
	rows, err := q.QueryContext(ctx, `
		SELECT id, name, slug
		FROM business_categories
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	categories := []Category{}
	for rows.Next() {
		var c Category
		err := rows.Scan(&c.ID, &c.Name, &c.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over category rows: %w", err)
	}

	return categories, nil
}
