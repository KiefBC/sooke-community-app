import Testing
import Foundation
@testable import SookeCommunity

@Suite("Business Model Tests")
struct BusinessModelTests {

    // MARK: - Business decoding

    @Test func decodesBusinessWithAllFields() throws {
        let json = """
        {
            "id": 1,
            "name": "Sooke Harbor House",
            "slug": "sooke-harbor-house",
            "description": "A waterfront restaurant",
            "category_name": "Restaurants",
            "category_slug": "restaurants",
            "address": "1528 Whiffen Spit Rd",
            "latitude": 48.3538,
            "longitude": -123.7256,
            "phone": "250-642-3421",
            "email": "info@sookeharbourhouse.com",
            "website": "https://sookeharbourhouse.com"
        }
        """
        let data = try #require(json.data(using: .utf8))
        let business = try JSONDecoder().decode(Business.self, from: data)

        #expect(business.id == 1)
        #expect(business.name == "Sooke Harbor House")
        #expect(business.slug == "sooke-harbor-house")
        #expect(business.description == "A waterfront restaurant")
        #expect(business.categoryName == "Restaurants")
        #expect(business.categorySlug == "restaurants")
        #expect(business.address == "1528 Whiffen Spit Rd")
        #expect(business.latitude == 48.3538)
        #expect(business.longitude == -123.7256)
        #expect(business.phone == "250-642-3421")
        #expect(business.email == "info@sookeharbourhouse.com")
        #expect(business.website == "https://sookeharbourhouse.com")
    }

    @Test func decodesBusinessWithNullOptionals() throws {
        let json = """
        {
            "id": 2,
            "name": "Corner Store",
            "slug": "corner-store",
            "description": null,
            "category_name": "Retail",
            "category_slug": "retail",
            "address": "123 Main St",
            "latitude": 48.3700,
            "longitude": -123.7100,
            "phone": null,
            "email": null,
            "website": null
        }
        """
        let data = try #require(json.data(using: .utf8))
        let business = try JSONDecoder().decode(Business.self, from: data)

        #expect(business.id == 2)
        #expect(business.name == "Corner Store")
        #expect(business.description == nil)
        #expect(business.phone == nil)
        #expect(business.email == nil)
        #expect(business.website == nil)
    }

    // MARK: - Business today_hours decoding

    @Test func decodesBusinessWithTodayHours() throws {
        let json = """
        {
            "id": 5,
            "name": "Open Shop",
            "slug": "open-shop",
            "description": null,
            "category_name": "Retail",
            "category_slug": "retail",
            "address": "1 Main St",
            "latitude": 48.37,
            "longitude": -123.72,
            "phone": null,
            "email": null,
            "website": null,
            "today_hours": {
                "day_of_week": 3,
                "open_time": "09:00:00",
                "close_time": "17:00:00",
                "is_closed": false
            }
        }
        """
        let data = try #require(json.data(using: .utf8))
        let business = try JSONDecoder().decode(Business.self, from: data)

        let hours = try #require(business.todayHours)
        #expect(hours.dayOfWeek == 3)
        #expect(hours.openTime == "09:00:00")
        #expect(hours.closeTime == "17:00:00")
        #expect(hours.isClosed == false)
    }

    @Test func decodesBusinessWithNullTodayHours() throws {
        let json = """
        {
            "id": 6,
            "name": "No Hours Shop",
            "slug": "no-hours-shop",
            "description": null,
            "category_name": "Retail",
            "category_slug": "retail",
            "address": "2 Main St",
            "latitude": 48.37,
            "longitude": -123.72,
            "phone": null,
            "email": null,
            "website": null,
            "today_hours": null
        }
        """
        let data = try #require(json.data(using: .utf8))
        let business = try JSONDecoder().decode(Business.self, from: data)

        #expect(business.todayHours == nil)
    }

    @Test func decodesBusinessWithMissingTodayHours() throws {
        let json = """
        {
            "id": 7,
            "name": "Legacy Shop",
            "slug": "legacy-shop",
            "description": null,
            "category_name": "Retail",
            "category_slug": "retail",
            "address": "3 Main St",
            "latitude": 48.37,
            "longitude": -123.72,
            "phone": null,
            "email": null,
            "website": null
        }
        """
        let data = try #require(json.data(using: .utf8))
        let business = try JSONDecoder().decode(Business.self, from: data)

        #expect(business.todayHours == nil)
    }

    // MARK: - BusinessDetails decoding

    @Test func decodesBusinessDetailsWithHoursAndMenus() throws {
        let json = """
        {
            "id": 1,
            "name": "Sooke Harbor House",
            "slug": "sooke-harbor-house",
            "description": "A waterfront restaurant",
            "category_name": "Restaurants",
            "category_slug": "restaurants",
            "address": "1528 Whiffen Spit Rd",
            "latitude": 48.3538,
            "longitude": -123.7256,
            "phone": "250-642-3421",
            "email": "info@sookeharbourhouse.com",
            "website": "https://sookeharbourhouse.com",
            "hours": [
                {"day_of_week": 1, "open_time": "09:00:00", "close_time": "17:00:00", "is_closed": false}
            ],
            "menus": [
                {
                    "id": 1,
                    "name": "Dinner",
                    "description": "Evening menu",
                    "items": [
                        {"id": 1, "name": "Salmon", "description": "Fresh local salmon", "price": "28.99"}
                    ]
                }
            ]
        }
        """
        let data = try #require(json.data(using: .utf8))
        let details = try JSONDecoder().decode(BusinessDetails.self, from: data)

        #expect(details.id == 1)
        #expect(details.name == "Sooke Harbor House")

        // Hours
        #expect(details.hours.count == 1)
        let hour = try #require(details.hours.first)
        #expect(hour.dayOfWeek == 1)
        #expect(hour.openTime == "09:00:00")
        #expect(hour.closeTime == "17:00:00")
        #expect(hour.isClosed == false)

        // Menus
        #expect(details.menus.count == 1)
        let menu = try #require(details.menus.first)
        #expect(menu.id == 1)
        #expect(menu.name == "Dinner")
        #expect(menu.description == "Evening menu")

        // Menu items
        #expect(menu.items.count == 1)
        let item = try #require(menu.items.first)
        #expect(item.id == 1)
        #expect(item.name == "Salmon")
        #expect(item.description == "Fresh local salmon")
        #expect(item.price == "28.99")
    }

    @Test func decodesBusinessDetailsWithNullMenuDescription() throws {
        let json = """
        {
            "id": 3,
            "name": "Pizzeria",
            "slug": "pizzeria",
            "description": null,
            "category_name": "Restaurants",
            "category_slug": "restaurants",
            "address": "5 Ocean Blvd",
            "latitude": 48.3600,
            "longitude": -123.7200,
            "phone": null,
            "email": null,
            "website": null,
            "hours": [],
            "menus": [
                {
                    "id": 2,
                    "name": "Lunch",
                    "description": null,
                    "items": [
                        {"id": 5, "name": "Pizza", "description": null, "price": "14.99"}
                    ]
                }
            ]
        }
        """
        let data = try #require(json.data(using: .utf8))
        let details = try JSONDecoder().decode(BusinessDetails.self, from: data)

        let menu = try #require(details.menus.first)
        #expect(menu.description == nil)

        let item = try #require(menu.items.first)
        #expect(item.description == nil)
        #expect(item.price == "14.99")
    }
}
