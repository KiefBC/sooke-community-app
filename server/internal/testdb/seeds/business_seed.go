package seeds

import "database/sql"

// BusinessSeed inserts 5 businesses with full hours (35 rows), 3 menus, and
// 9 menu items. Calls CategorySeed and UserSeed first.
func BusinessSeed(tx *sql.Tx) {
	CategorySeed(tx)
	UserSeed(tx)

	// 5 businesses - one per category, mixed owner assignments
	Exec(tx, `
		INSERT INTO businesses (owner_id, category_id, name, slug, address, latitude, longitude) VALUES
			((SELECT id FROM users WHERE clerk_id = 'seed_general_user'),
			 (SELECT id FROM business_categories WHERE slug = 'restaurant'),
			 'Sooke Harbour House', 'sooke-harbour-house', '1528 Whiffen Spit Rd', 48.3538, -123.7256),

			((SELECT id FROM users WHERE clerk_id = 'seed_business_owner'),
			 (SELECT id FROM business_categories WHERE slug = 'cafe'),
			 'Moms Cafe', 'moms-cafe', '2036 Shields Rd', 48.3761, -123.7254),

			((SELECT id FROM users WHERE clerk_id = 'seed_super_admin'),
			 (SELECT id FROM business_categories WHERE slug = 'outdoor-recreation'),
			 'Sooke Landing Marina', 'sooke-landing-marina', '4549 Sooke Rd', 48.3720, -123.7100),

			(NULL,
			 (SELECT id FROM business_categories WHERE slug = 'community'),
			 'Sooke Community Hall', 'sooke-community-hall', '2037 Shields Rd', 48.3760, -123.7250),

			((SELECT id FROM users WHERE clerk_id = 'seed_owner_two'),
			 (SELECT id FROM business_categories WHERE slug = 'retail'),
			 'Wandering Moose Outfitters', 'wandering-moose-outfitters', '6691 Sooke Rd', 48.3750, -123.7200)
		ON CONFLICT DO NOTHING;
	`)

	// Hours - all 7 days for each business (35 rows total)

	// Sooke Harbour House: Wed-Sun 17:00-21:00, Mon-Tue closed
	Exec(tx, `
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed) VALUES
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 0, '17:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 1, '17:00', '21:00', true),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 2, '17:00', '21:00', true),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 3, '17:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 4, '17:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 5, '17:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 6, '17:00', '21:00', false)
		ON CONFLICT DO NOTHING;
	`)

	// Moms Cafe: Mon-Fri 07:00-15:00, Sat-Sun 08:00-15:00
	Exec(tx, `
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed) VALUES
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 0, '08:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 1, '07:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 2, '07:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 3, '07:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 4, '07:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 5, '07:00', '15:00', false),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 6, '08:00', '15:00', false)
		ON CONFLICT DO NOTHING;
	`)

	// Sooke Landing Marina: Mon-Sat 08:00-17:00, Sun 09:00-16:00
	Exec(tx, `
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed) VALUES
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 0, '09:00', '16:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 1, '08:00', '17:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 2, '08:00', '17:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 3, '08:00', '17:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 4, '08:00', '17:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 5, '08:00', '17:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'), 6, '08:00', '17:00', false)
		ON CONFLICT DO NOTHING;
	`)

	// Sooke Community Hall: Mon-Fri 09:00-21:00, Sat 09:00-17:00, Sun closed
	Exec(tx, `
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed) VALUES
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 0, '09:00', '21:00', true),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 1, '09:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 2, '09:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 3, '09:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 4, '09:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 5, '09:00', '21:00', false),
			((SELECT id FROM businesses WHERE slug = 'sooke-community-hall'), 6, '09:00', '17:00', false)
		ON CONFLICT DO NOTHING;
	`)

	// Wandering Moose Outfitters: Mon-Sat 10:00-18:00, Sun 11:00-16:00
	Exec(tx, `
		INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed) VALUES
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 0, '11:00', '16:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 1, '10:00', '18:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 2, '10:00', '18:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 3, '10:00', '18:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 4, '10:00', '18:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 5, '10:00', '18:00', false),
			((SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'), 6, '10:00', '18:00', false)
		ON CONFLICT DO NOTHING;
	`)

	// Menus - 3 total (Harbour House has 1, Moms has 2)
	Exec(tx, `
		INSERT INTO menus (business_id, name, description) VALUES
			((SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'), 'Dinner', 'Fresh Pacific Northwest cuisine'),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 'Breakfast', 'Classic home-style breakfast'),
			((SELECT id FROM businesses WHERE slug = 'moms-cafe'), 'Lunch', 'Soups, sandwiches, and salads')
		ON CONFLICT DO NOTHING;
	`)

	// Menu items - 3 per menu (9 total)
	Exec(tx, `
		INSERT INTO menu_items (menu_id, name, price) VALUES
			((SELECT id FROM menus WHERE name = 'Dinner' AND business_id = (SELECT id FROM businesses WHERE slug = 'sooke-harbour-house')),
			 'Pan-Seared Salmon', 32.00),
			((SELECT id FROM menus WHERE name = 'Dinner' AND business_id = (SELECT id FROM businesses WHERE slug = 'sooke-harbour-house')),
			 'Braised Short Rib', 38.00),
			((SELECT id FROM menus WHERE name = 'Dinner' AND business_id = (SELECT id FROM businesses WHERE slug = 'sooke-harbour-house')),
			 'Dungeness Crab Cake', 24.00),

			((SELECT id FROM menus WHERE name = 'Breakfast' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Moms Big Breakfast', 14.99),
			((SELECT id FROM menus WHERE name = 'Breakfast' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Blueberry Pancakes', 12.99),
			((SELECT id FROM menus WHERE name = 'Breakfast' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Eggs Benedict', 16.99),

			((SELECT id FROM menus WHERE name = 'Lunch' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Grilled Cheese & Tomato Soup', 13.99),
			((SELECT id FROM menus WHERE name = 'Lunch' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Turkey Club', 15.99),
			((SELECT id FROM menus WHERE name = 'Lunch' AND business_id = (SELECT id FROM businesses WHERE slug = 'moms-cafe')),
			 'Caesar Salad', 12.49);
	`)
}
