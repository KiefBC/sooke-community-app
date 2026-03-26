package repository

import (
	"context"
	"fmt"
)

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

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
