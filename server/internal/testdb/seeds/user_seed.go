package seeds

import "database/sql"

// UserSeed inserts 3 user roles and 5 users. Leaf seed, no dependencies.
func UserSeed(tx *sql.Tx) {
	Exec(tx, `
		INSERT INTO user_roles (name) VALUES
			('general_user'),
			('business_owner'),
			('super_admin')
		ON CONFLICT DO NOTHING;
	`)

	Exec(tx, `
		INSERT INTO users (clerk_id, email, display_name, role_id) VALUES
			('seed_general_user',   'user@sooke.dev',      'Sooke Resident',  (SELECT id FROM user_roles WHERE name = 'general_user')),
			('seed_business_owner', 'owner@sooke.dev',     'Business Owner',  (SELECT id FROM user_roles WHERE name = 'business_owner')),
			('seed_super_admin',    'admin@sooke.dev',     'Super Admin',     (SELECT id FROM user_roles WHERE name = 'super_admin')),
			('seed_resident_two',   'resident2@sooke.dev', 'Second Resident', (SELECT id FROM user_roles WHERE name = 'general_user')),
			('seed_owner_two',      'owner2@sooke.dev',    'Second Owner',    (SELECT id FROM user_roles WHERE name = 'business_owner'))
		ON CONFLICT DO NOTHING;
	`)
}
