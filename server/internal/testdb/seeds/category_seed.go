package seeds

import "database/sql"

// CategorySeed inserts 5 business categories.
func CategorySeed(tx *sql.Tx) {
	Exec(tx, `
		INSERT INTO business_categories (name, slug) VALUES
			('Restaurant', 'restaurant'),
			('Cafe', 'cafe'),
			('Retail', 'retail'),
			('Outdoor Recreation', 'outdoor-recreation'),
			('Community', 'community')
		ON CONFLICT DO NOTHING;
	`)
}
