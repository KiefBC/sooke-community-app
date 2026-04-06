package repository

import (
	"context"
	"database/sql"
	"fmt"
)

// Querier is the shared interface between *sql.DB and *sql.Tx.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Business represents a business entity in the database.
type Business struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Slug         string        `json:"slug"`
	Description  *string       `json:"description"` // nullable
	CategoryName string        `json:"category_name"`
	CategorySlug string        `json:"category_slug"`
	Address      string        `json:"address"`
	Latitude     float64       `json:"latitude"`
	Longitude    float64       `json:"longitude"`
	Phone        *string       `json:"phone"`       // nullable
	Email        *string       `json:"email"`       // nullable
	Website      *string       `json:"website"`     // nullable
	TodayHours   *BusinessHour `json:"today_hours"` // nullable - only today's hours for list view
}

// BusinessDetails represents a business along with its hours and menus.
type BusinessDetails struct {
	Business
	Hours []BusinessHour `json:"hours"`
	Menus []Menu         `json:"menus"`
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
	// The database stores price as NUMERIC(10,2). The pgx driver returns this as a string (`"12.99"`). Converting to `float64` introduces rounding
	Price string `json:"price"`
}

// ListBusinesses retrieves a list of businesses from the database based on the provided search and category filters, along with pagination parameters. It returns the list of businesses, the total count of matching businesses (ignoring pagination), and any error encountered during the operation.
func ListBusinesses(ctx context.Context, q Querier, search, categorySlug, tz string, limit, offset int) ([]Business, int, error) {
	var countTotal int
	err := q.QueryRowContext(ctx,
		`SELECT COUNT(*)
		 FROM businesses b
		 JOIN business_categories bc ON b.category_id = bc.id
		 WHERE ($1 = '' OR b.name ILIKE '%' || $1 || '%')
			 AND ($2 = '' OR bc.slug = $2)`,
		search, categorySlug,
	).Scan(&countTotal)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count businesses: %w", err)
	}

	// The $1 = '' trick prevents dynamic sql with string concatentation. The OR short-circuits to always true
	rows, err := q.QueryContext(ctx,
		`SELECT b.id, b.name, b.slug, b.description, bc.name AS category_name, bc.slug AS category_slug,
		        b.address, b.latitude, b.longitude, b.phone, b.email, b.website,
		        bh.day_of_week, bh.open_time, bh.close_time, bh.is_closed
		 FROM businesses b
		 JOIN business_categories bc ON b.category_id = bc.id
		 LEFT JOIN business_hours bh ON bh.business_id = b.id
		     AND bh.day_of_week = EXTRACT(DOW FROM NOW() AT TIME ZONE $3)
		 WHERE ($1 = '' OR b.name ILIKE '%' || $1 || '%')
		     AND ($2 = '' OR bc.slug = $2)
		 ORDER BY b.name ASC
		 LIMIT $4 OFFSET $5`,
		search, categorySlug, tz, limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query businesses: %w", err)
	}
	defer rows.Close()

	businesses := []Business{}
	for rows.Next() {
		var b Business
		var bhDayOfWeek sql.NullInt64
		var bhOpenTime, bhCloseTime sql.NullString
		var bhIsClosed sql.NullBool
		err := rows.Scan(
			&b.ID, &b.Name, &b.Slug, &b.Description, &b.CategoryName, &b.CategorySlug,
			&b.Address, &b.Latitude, &b.Longitude, &b.Phone, &b.Email, &b.Website,
			&bhDayOfWeek, &bhOpenTime, &bhCloseTime, &bhIsClosed,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan business: %w", err)
		}
		if bhDayOfWeek.Valid {
			b.TodayHours = &BusinessHour{
				DayOfWeek: int(bhDayOfWeek.Int64),
				OpenTime:  bhOpenTime.String,
				CloseTime: bhCloseTime.String,
				IsClosed:  bhIsClosed.Bool,
			}
		}
		businesses = append(businesses, b)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over business rows: %w", err)
	}

	return businesses, countTotal, nil
}

// GetBusinessBySlug retrieves a business from the database based on its slug, along with its operating hours and associated menus. It returns a BusinessDetails struct containing all the relevant information, or nil if no business is found with the given slug. Any error encountered during the operation is also returned.
func GetBusinessBySlug(ctx context.Context, q Querier, slug string) (*BusinessDetails, error) {
	var bd BusinessDetails
	err := q.QueryRowContext(ctx,
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

	bd.Hours = []BusinessHour{}
	bd.Menus = []Menu{}

	hoursRows, err := q.QueryContext(ctx,
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

	// Collect all menus first, then close menuRows before querying items.
	// A *sql.Tx uses a single connection, so you can't have two active
	// result sets open at the same time (pgx returns "bad connection").
	menuRows, err := q.QueryContext(ctx,
		`SELECT id, name, description
		 FROM menus
		 WHERE business_id = $1`,
		bd.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query menus: %w", err)
	}

	for menuRows.Next() {
		var m Menu
		err := menuRows.Scan(&m.ID, &m.Name, &m.Description)
		if err != nil {
			menuRows.Close()
			return nil, fmt.Errorf("failed to scan menu: %w", err)
		}
		bd.Menus = append(bd.Menus, m)
	}
	if err := menuRows.Err(); err != nil {
		menuRows.Close()
		return nil, fmt.Errorf("error iterating over menus rows: %w", err)
	}
	menuRows.Close()

	// Now query items for each menu with no other result set open
	for i := range bd.Menus {
		itemRows, err := q.QueryContext(ctx,
			`SELECT id, name, description, price
			 FROM menu_items
			 WHERE menu_id = $1
			 ORDER BY name ASC`,
			bd.Menus[i].ID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query menu items: %w", err)
		}

		for itemRows.Next() {
			var mi MenuItem
			err := itemRows.Scan(&mi.ID, &mi.Name, &mi.Description, &mi.Price)
			if err != nil {
				itemRows.Close()
				return nil, fmt.Errorf("failed to scan menu item: %w", err)
			}
			bd.Menus[i].Items = append(bd.Menus[i].Items, mi)
		}
		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, fmt.Errorf("error iterating over menu items rows: %w", err)
		}
		itemRows.Close()
	}

	return &bd, nil
}
