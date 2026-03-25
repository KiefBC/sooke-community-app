package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kiefbc/sooke_app/server/internal/database"
	"github.com/kiefbc/sooke_app/server/internal/slug"
)

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load("./.env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	db, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := seed(db); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully")
}

func seed(db *sql.DB) error {
	var seedUserRoles = []struct {
		name string
	}{
		{
			"general_user",
		},
		{
			"business_owner",
		},
		{
			"super_admin",
		},
	}

	for _, role := range seedUserRoles {
		if _, err := db.Exec("INSERT INTO user_roles (name) VALUES ($1) ON CONFLICT (name) DO NOTHING", role.name); err != nil {
			return fmt.Errorf("failed to seed role %q: %w", role.name, err)
		}
	}

	var seedUsers = []struct {
		ClerkID     string
		Email       string
		DisplayName string
		Role        string
	}{
		{"seed_general_user", "user@sooke.dev", "Sooke Resident", "general_user"},
		{"seed_business_owner", "owner@sooke.dev", "Business Owner", "business_owner"},
		{"seed_super_admin", "admin@sooke.dev", "Super Admin", "super_admin"},
	}

	for _, u := range seedUsers {
		if _, err := db.Exec(
			`INSERT INTO users (clerk_id, role_id, email, display_name) VALUES ($1, (SELECT id FROM user_roles WHERE name = $2), $3, $4)
			 ON CONFLICT (clerk_id) DO UPDATE SET role_id = (SELECT id FROM user_roles WHERE name = $2), email = $3, display_name = $4`,
			u.ClerkID, u.Role, u.Email, u.DisplayName,
		); err != nil {
			return fmt.Errorf("failed to seed user %q: %w", u.Email, err)
		}
	}

	var seedBusinessCat = []struct {
		name string
	}{
		{
			"Restaurant",
		},
		{
			"Cafe",
		},
		{
			"Retail",
		},
		{
			"Outdoor Recreation",
		},
		{
			"Community",
		},
	}

	for _, cat := range seedBusinessCat {
		if _, err := db.Exec("INSERT INTO business_categories (name, slug) VALUES ($1, $2) ON CONFLICT (name) DO NOTHING", cat.name, slug.GenerateSlug(cat.name)); err != nil {
			return fmt.Errorf("failed to seed business category %q: %w", cat.name, err)
		}
	}

	var seedEventType = []struct {
		name string
	}{
		{
			"Live Music",
		},
		{
			"Market",
		},
		{
			"Community Meeting",
		},
	}

	for _, et := range seedEventType {
		if _, err := db.Exec("INSERT INTO event_types (name, slug) VALUES ($1, $2) ON CONFLICT (slug) DO NOTHING", et.name, slug.GenerateSlug(et.name)); err != nil {
			return fmt.Errorf("failed to seed event type %q: %w", et.name, err)
		}
	}

	var seedBusinesses = []struct {
		Name         string
		Description  string
		Address      string
		Lat          float64
		Lng          float64
		Category     string
		OwnerClerkID string
		Phone        string
		Email        string
		Website      string
	}{
		{"Sooke Harbour House", "Historic waterfront inn with Pacific Northwest fine dining", "1528 Whiffen Spit Rd", 48.356618349381755, -123.72733056442932, "Restaurant", "seed_general_user", "111-111-1111", "info@sookeharbourhouse.dev", "https://sookeharbourhouse.dev"},
		{"Mom's Cafe", "Family-friendly breakfast and lunch spot loved by locals", "2036 Shields Rd", 48.377112314153386, -123.7254915288472, "Cafe", "seed_business_owner", "222-222-2222", "hello@momscafe.dev", "https://momscafe.dev"},
		{"Sooke Landing Marina", "Full-service marina with boat rentals and moorage", "6585 Goodmere Rd", 48.37801614501773, -123.71648471022064, "Outdoor Recreation", "seed_super_admin", "333-333-3333", "dock@sookemarina.dev", "https://sookemarina.dev"},
		{"King Tide Fishing Charters", "Guided salmon and halibut fishing on the Strait of Juan de Fuca", "6969 Sea Lion Way", 48.35785881994283, -123.72657954591688, "Outdoor Recreation", "", "", "", ""},
		{"Sooke Community Hall", "Local gathering space for markets, meetings, and events", "2037 Shields Rd", 48.37759971334821, -123.72517166628057, "Community", "", "", "", ""},
		{"The Stick In The Mud's Roastoreum", "Small-batch coffee roaster and cozy neighbourhood cafe", "6711 Eustace Rd", 48.37789738050146, -123.7245055441908, "Cafe", "", "", "", ""},
	}

	for _, biz := range seedBusinesses {
		var ownerClerkID interface{}
		if biz.OwnerClerkID != "" {
			ownerClerkID = biz.OwnerClerkID
		}
		if _, err := db.Exec(
			`INSERT INTO businesses (owner_id, category_id, name, slug, description, phone, email, website, address, latitude, longitude)
			 VALUES ((SELECT id FROM users WHERE clerk_id = $1), (SELECT id FROM business_categories WHERE name = $2), $3, $4, $5, $6, $7, $8, $9, $10, $11)
			 ON CONFLICT (slug) DO UPDATE SET owner_id = COALESCE((SELECT id FROM users WHERE clerk_id = $1), businesses.owner_id),
			 category_id = (SELECT id FROM business_categories WHERE name = $2),
			 name = $3, description = $5, phone = $6, email = $7, website = $8, address = $9, latitude = $10, longitude = $11`,
			ownerClerkID, biz.Category, biz.Name, slug.GenerateSlug(biz.Name), biz.Description, biz.Phone, biz.Email, biz.Website, biz.Address, biz.Lat, biz.Lng,
		); err != nil {
			return fmt.Errorf("failed to seed business %q: %w", biz.Name, err)
		}
	}

	var seedBusinessHours = []struct {
		Business  string
		DayOfWeek int
		OpenTime  string
		CloseTime string
		IsClosed  bool
	}{
		// Sooke Harbour House -- Wed-Sun 17:00-21:00, closed Mon-Tue
		{"Sooke Harbour House", 0, "17:00", "21:00", false}, // Sun
		{"Sooke Harbour House", 1, "00:00", "00:00", true},  // Mon
		{"Sooke Harbour House", 2, "00:00", "00:00", true},  // Tue
		{"Sooke Harbour House", 3, "17:00", "21:00", false}, // Wed
		{"Sooke Harbour House", 4, "17:00", "21:00", false}, // Thu
		{"Sooke Harbour House", 5, "17:00", "21:00", false}, // Fri
		{"Sooke Harbour House", 6, "17:00", "21:00", false}, // Sat

		// Mom's Cafe -- Mon-Fri 07:00-15:00, Sat-Sun 08:00-15:00
		{"Mom's Cafe", 0, "08:00", "15:00", false}, // Sun
		{"Mom's Cafe", 1, "07:00", "15:00", false}, // Mon
		{"Mom's Cafe", 2, "07:00", "15:00", false}, // Tue
		{"Mom's Cafe", 3, "07:00", "15:00", false}, // Wed
		{"Mom's Cafe", 4, "07:00", "15:00", false}, // Thu
		{"Mom's Cafe", 5, "07:00", "15:00", false}, // Fri
		{"Mom's Cafe", 6, "08:00", "15:00", false}, // Sat

		// Sooke Landing Marina -- Mon-Sat 08:00-17:00, Sun 09:00-16:00
		{"Sooke Landing Marina", 0, "09:00", "16:00", false}, // Sun
		{"Sooke Landing Marina", 1, "08:00", "17:00", false}, // Mon
		{"Sooke Landing Marina", 2, "08:00", "17:00", false}, // Tue
		{"Sooke Landing Marina", 3, "08:00", "17:00", false}, // Wed
		{"Sooke Landing Marina", 4, "08:00", "17:00", false}, // Thu
		{"Sooke Landing Marina", 5, "08:00", "17:00", false}, // Fri
		{"Sooke Landing Marina", 6, "08:00", "17:00", false}, // Sat

		// King Tide Fishing Charters -- Mon-Sun 05:00-18:00
		{"King Tide Fishing Charters", 0, "05:00", "18:00", false}, // Sun
		{"King Tide Fishing Charters", 1, "05:00", "18:00", false}, // Mon
		{"King Tide Fishing Charters", 2, "05:00", "18:00", false}, // Tue
		{"King Tide Fishing Charters", 3, "05:00", "18:00", false}, // Wed
		{"King Tide Fishing Charters", 4, "05:00", "18:00", false}, // Thu
		{"King Tide Fishing Charters", 5, "05:00", "18:00", false}, // Fri
		{"King Tide Fishing Charters", 6, "05:00", "18:00", false}, // Sat

		// Sooke Community Hall -- Mon-Fri 09:00-21:00, Sat 09:00-17:00, closed Sun
		{"Sooke Community Hall", 0, "00:00", "00:00", true},  // Sun
		{"Sooke Community Hall", 1, "09:00", "21:00", false}, // Mon
		{"Sooke Community Hall", 2, "09:00", "21:00", false}, // Tue
		{"Sooke Community Hall", 3, "09:00", "21:00", false}, // Wed
		{"Sooke Community Hall", 4, "09:00", "21:00", false}, // Thu
		{"Sooke Community Hall", 5, "09:00", "21:00", false}, // Fri
		{"Sooke Community Hall", 6, "09:00", "17:00", false}, // Sat

		// The Stick In The Mud's Roastoreum -- Mon-Fri 06:30-17:00, Sat 07:00-17:00, closed Sun
		{"The Stick In The Mud's Roastoreum", 0, "00:00", "00:00", true},  // Sun
		{"The Stick In The Mud's Roastoreum", 1, "06:30", "17:00", false}, // Mon
		{"The Stick In The Mud's Roastoreum", 2, "06:30", "17:00", false}, // Tue
		{"The Stick In The Mud's Roastoreum", 3, "06:30", "17:00", false}, // Wed
		{"The Stick In The Mud's Roastoreum", 4, "06:30", "17:00", false}, // Thu
		{"The Stick In The Mud's Roastoreum", 5, "06:30", "17:00", false}, // Fri
		{"The Stick In The Mud's Roastoreum", 6, "07:00", "17:00", false}, // Sat
	}

	for _, bh := range seedBusinessHours {
		if _, err := db.Exec(
			`INSERT INTO business_hours (business_id, day_of_week, open_time, close_time, is_closed)
			 VALUES ((SELECT id FROM businesses WHERE name = $1), $2, $3, $4, $5)
			 ON CONFLICT (business_id, day_of_week) DO UPDATE SET open_time = $3, close_time = $4, is_closed = $5`,
			bh.Business, bh.DayOfWeek, bh.OpenTime, bh.CloseTime, bh.IsClosed,
		); err != nil {
			return fmt.Errorf("failed to seed business hours for %q day %d: %w", bh.Business, bh.DayOfWeek, err)
		}
	}

	var seedMenus = []struct {
		Business    string
		Name        string
		Description string
		Items       []struct {
			Name        string
			Description string
			Price       float64
		}
	}{
		{
			"Sooke Harbour House", "Dinner", "Fresh Pacific Northwest cuisine",
			[]struct {
				Name        string
				Description string
				Price       float64
			}{
				{"Pan-Seared Salmon", "Wild BC salmon with seasonal vegetables", 32.00},
				{"Braised Short Rib", "Slow-cooked with root vegetables and red wine jus", 38.00},
				{"Dungeness Crab Cake", "Local crab with remoulade and micro greens", 24.00},
			},
		},
		{
			"Mom's Cafe", "Breakfast", "Classic home-style breakfast",
			[]struct {
				Name        string
				Description string
				Price       float64
			}{
				{"Mom's Big Breakfast", "Two eggs, bacon, toast, and hash browns", 14.99},
				{"Blueberry Pancakes", "Stack of three with maple syrup", 12.99},
				{"Eggs Benedict", "Poached eggs with hollandaise on English muffin", 16.99},
			},
		},
		{
			"The Stick In The Mud's Roastoreum", "Drinks", "Freshly roasted coffee and specialty drinks",
			[]struct {
				Name        string
				Description string
				Price       float64
			}{
				{"House Drip Coffee", "Freshly roasted single origin", 4.50},
				{"Vanilla Latte", "Double shot with house-made vanilla syrup", 6.00},
				{"Matcha Latte", "Ceremonial grade matcha with oat milk", 6.50},
			},
		},
	}

	for _, menu := range seedMenus {
		var menuID int64
		err := db.QueryRow(
			`INSERT INTO menus (business_id, name, description)
			 VALUES ((SELECT id FROM businesses WHERE name = $1), $2, $3)
			 ON CONFLICT (business_id, name) DO NOTHING
			 RETURNING id`,
			menu.Business, menu.Name, menu.Description,
		).Scan(&menuID)
		if err != nil {
			if err == sql.ErrNoRows {
				continue // already seeded
			}
			return fmt.Errorf("failed to seed menu %q for %q: %w", menu.Name, menu.Business, err)
		}

		for _, item := range menu.Items {
			if _, err := db.Exec(
				"INSERT INTO menu_items (menu_id, name, description, price) VALUES ($1, $2, $3, $4)",
				menuID, item.Name, item.Description, item.Price,
			); err != nil {
				return fmt.Errorf("failed to seed menu item %q: %w", item.Name, err)
			}
		}
	}

	var seedEvents = []struct {
		Name        string
		Business    string
		StartTime   string
		EndTime     string
		Description string
		EventType   string
		Status      string
	}{
		{"Friday Night Jazz", "Sooke Harbour House", "2026-04-04T19:00:00-07:00", "2026-04-04T22:00:00-07:00", "Live jazz performance featuring local musicians at the waterfront dining room", "Live Music", "approved"},
		{"Sooke Saturday Market", "Sooke Community Hall", "2026-04-05T09:00:00-07:00", "2026-04-05T14:00:00-07:00", "Weekly community market with local produce, crafts, and baked goods", "Market", "approved"},
	}

	for _, event := range seedEvents {
		var endTime interface{}
		if event.EndTime != "" {
			endTime = event.EndTime
		}
		if _, err := db.Exec(
			`INSERT INTO events (event_type_id, submitted_by, business_id, name, slug, description, starts_at, ends_at, status)
			 VALUES (
			   (SELECT id FROM event_types WHERE name = $1),
			   (SELECT id FROM users WHERE clerk_id = 'seed_super_admin'),
			   (SELECT id FROM businesses WHERE name = $2),
			   $3, $4, $5, $6, $7, $8
			 )
			 ON CONFLICT (slug) DO UPDATE SET description = $5, starts_at = $6, ends_at = $7, status = $8`,
			event.EventType, event.Business, event.Name, slug.GenerateSlug(event.Name),
			event.Description, event.StartTime, endTime, event.Status,
		); err != nil {
			return fmt.Errorf("failed to seed event %q: %w", event.Name, err)
		}
	}

	return nil
}
