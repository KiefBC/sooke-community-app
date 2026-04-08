package seeds

import "database/sql"

// EventSeed inserts 3 event types and 5 events. Calls BusinessSeed first
// (which pulls in categories and users).
func EventSeed(tx *sql.Tx) {
	BusinessSeed(tx)

	Exec(tx, `
		INSERT INTO event_types (name, slug) VALUES
			('Live Music', 'live-music'),
			('Market', 'market'),
			('Community Meeting', 'community-meeting')
		ON CONFLICT DO NOTHING;
	`)

	// 5 events across different businesses, types, submitters, and statuses.
	// Events with a business_id must NOT have lat/long (schema constraint).

	// 1. Approved event at a business
	Exec(tx, `
		INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES
			((SELECT id FROM event_types WHERE slug = 'live-music'),
			 (SELECT id FROM users WHERE clerk_id = 'seed_super_admin'),
			 (SELECT id FROM businesses WHERE slug = 'sooke-harbour-house'),
			 'Friday Night Jazz', 'friday-night-jazz',
			 NOW() + INTERVAL '7 days', 'approved');
	`)

	// 2. Approved market event at community hall
	Exec(tx, `
		INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES
			((SELECT id FROM event_types WHERE slug = 'market'),
			 (SELECT id FROM users WHERE clerk_id = 'seed_super_admin'),
			 (SELECT id FROM businesses WHERE slug = 'sooke-community-hall'),
			 'Sooke Saturday Market', 'sooke-saturday-market',
			 NOW() + INTERVAL '8 days', 'approved');
	`)

	// 3. Approved community meeting at marina
	Exec(tx, `
		INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES
			((SELECT id FROM event_types WHERE slug = 'community-meeting'),
			 (SELECT id FROM users WHERE clerk_id = 'seed_business_owner'),
			 (SELECT id FROM businesses WHERE slug = 'sooke-landing-marina'),
			 'Marina Open House', 'marina-open-house',
			 NOW() + INTERVAL '14 days', 'approved');
	`)

	// 4. Pending review event
	Exec(tx, `
		INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES
			((SELECT id FROM event_types WHERE slug = 'live-music'),
			 (SELECT id FROM users WHERE clerk_id = 'seed_owner_two'),
			 (SELECT id FROM businesses WHERE slug = 'moms-cafe'),
			 'Cafe Acoustic Night', 'cafe-acoustic-night',
			 NOW() + INTERVAL '10 days', 'pending_review');
	`)

	// 5. Draft event
	Exec(tx, `
		INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, starts_at, status) VALUES
			((SELECT id FROM event_types WHERE slug = 'market'),
			 (SELECT id FROM users WHERE clerk_id = 'seed_general_user'),
			 (SELECT id FROM businesses WHERE slug = 'wandering-moose-outfitters'),
			 'Spring Craft Fair', 'spring-craft-fair',
			 NOW() + INTERVAL '21 days', 'draft');
	`)
}
