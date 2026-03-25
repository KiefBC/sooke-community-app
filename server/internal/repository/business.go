package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Why are we using pointers for nullable fields?
// We use pointers for nullable fields in the Business struct to allow us to represent the absence of a value (i.e., null) in the database.

// Business represents a business entity in the database.
type Business struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  *string `json:"description"` // nullable
	CategoryName string  `json:"category_name"`
	CategorySlug string  `json:"category_slug"`
	Address      string  `json:"address"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	Phone        *string `json:"phone"`   // nullable
	Email        *string `json:"email"`   // nullable
	Website      *string `json:"website"` // nullable
}

// BusinessDetails represents a business along with its hours and menus.
type BusinessDetails struct {
	Business                // Embed the Business struct to include its fields
	Hours    []BusinessHour `json:"hours"`
	Menus    []Menu         `json:"menus"`
}

// BusinessHour represents the operating hours for a business on a specific day of the week.
type BusinessHour struct {
	DayOfWeek int    `json:"day_of_week"` // 0 = Sunday, 1 = Monday, ..., 6 = Saturday
	OpenTime  string `json:"open_time"`   // Format: "HH:MM:SS"
	CloseTime string `json:"close_time"`  // Format: "HH:MM:SS"
	IsClosed  bool   `json:"is_closed"`   // Indicates if the business is closed on this day
}

// Menu represents a menu associated with a business, containing multiple menu items.
type Menu struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"` // nullable
	Items       []MenuItem `json:"items"`
}

// MenuItem represents an individual item on a menu, including its name, description, and price.
type MenuItem struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"` // nullable
	// Why use a string for the Price field in MenuItem instead of a numeric type?
	// The database stores price as NUMERIC(10,2). The pgx driver returns this as a string (`"12.99"`). Converting to `float64` introduces rounding
	Price string `json:"price"`
}

func ListBusinesses(ctx context.Context, db *sql.DB, search, category_slug string, limit, offset int) ([]Business, int, error) {
	var countTotal int
	err := db.QueryRowContext(ctx,
		`SELECT COUNT(*)
		 FROM businesses b
		 JOIN business_categories bc ON b.category_id = bc.id
		 WHERE ($1 = '' OR b.name ILIKE '%' || $1 || '%')
			 AND ($2 = '' OR bc.slug = $2)`,
		search, category_slug,
	).Scan(&countTotal)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count businesses: %w", err)
	}

	// The $1 = '' trick prevents dynamic sql injections with string concatentation. The OR short-circuits to always true
	rows, err := db.QueryContext(ctx,
		`SELECT b.id, b.name, b.slug, b.description, bc.name AS category_name, bc.slug AS category_slug, b.address, b.latitude, b.longitude, b.phone, b.email, b.website
			 FROM businesses b
		JOIN business_categories bc ON b.category_id = bc.id
		WHERE ($1 = '' OR b.name ILIKE '%' || $1 || '%')
			AND ($2 = '' OR bc.slug = $2)
		ORDER BY b.name ASC
		LIMIT $3 OFFSET $4`,
		search, category_slug, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query businesses: %w", err)
	}
	defer rows.Close()

	businesses := []Business{} // Initialize an empty slice to hold the results and costs less memory than preallocating with a large capacity
	for rows.Next() {
		var b Business
		err := rows.Scan(&b.ID, &b.Name, &b.Slug, &b.Description, &b.CategoryName, &b.CategorySlug, &b.Address, &b.Latitude, &b.Longitude, &b.Phone, &b.Email, &b.Website)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan business: %w", err)
		}
		businesses = append(businesses, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over business rows: %w", err)
	}

	return businesses, countTotal, nil
}

func GetBusinessBySlug(ctx context.Context, db *sql.DB, slug string) (*BusinessDetails, error) {
	var bd BusinessDetails
	err := db.QueryRowContext(ctx,
		`SELECT b.id, b.name, b.slug, b.description, bc.name AS category_name, bc.slug AS category_slug, b.address, b.latitude, b.longitude, b.phone, b.email, b.website
		 FROM businesses b
		 JOIN business_categories bc ON b.category_id = bc.id
		 WHERE b.slug = $1`,
		slug,
	).Scan(&bd.ID, &bd.Name, &bd.Slug, &bd.Description, &bd.CategoryName, &bd.CategorySlug, &bd.Address, &bd.Latitude, &bd.Longitude, &bd.Phone, &bd.Email, &bd.Website)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No business found with the given slug
		}
		return nil, fmt.Errorf("failed to query business: %w", err)
	}

	hoursRows, err := db.QueryContext(ctx,
		`SELECT day_of_week, open_time, close_time, is_closed
		 FROM business_hours
		 WHERE business_id = $1
		 ORDER BY day_of_week ASC`,
		bd.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query business hours: %w", err)
	}
	defer hoursRows.Close()

	for hoursRows.Next() {
		var bh BusinessHour
		err := hoursRows.Scan(&bh.DayOfWeek, &bh.OpenTime, &bh.CloseTime, &bh.IsClosed)
		if err != nil {
			return nil, fmt.Errorf("failed to scan business hour: %w", err)
		}
		bd.Hours = append(bd.Hours, bh)
	}
	if err := hoursRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over business hours rows: %w", err)
	}

	menuRows, err := db.QueryContext(ctx,
		`SELECT id, name, description
		 FROM menus
		 WHERE business_id = $1`,
		bd.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query menus: %w", err)
	}
	defer menuRows.Close()

	for menuRows.Next() {
		var m Menu
		err := menuRows.Scan(&m.ID, &m.Name, &m.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to scan menu: %w", err)
		}

		itemRows, err := db.QueryContext(ctx,
			`SELECT id, name, description, price
			 FROM menu_items
			 WHERE menu_id = ANY($1)
			 ORDER BY name ASC`,
			m.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query menu items: %w", err)
		}

		for itemRows.Next() {
			var mi MenuItem
			err := itemRows.Scan(&mi.ID, &mi.Name, &mi.Description, &mi.Price)
			if err != nil {
				return nil, fmt.Errorf("failed to scan menu item: %w", err)
			}
			m.Items = append(m.Items, mi)
		}
		if err := itemRows.Err(); err != nil {
			return nil, fmt.Errorf("error iterating over menu items rows: %w", err)
		}
		bd.Menus = append(bd.Menus, m)
		itemRows.Close() // Close the item rows before the next iteration to prevent too many open connections
	}
	if err := menuRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over menus rows: %w", err)
	}

	return &bd, nil
}
