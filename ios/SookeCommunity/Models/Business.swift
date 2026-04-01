import Foundation

// Business represents a business list item returned by the API.
struct Business: Codable, Identifiable, Sendable, Hashable {
    let id: Int64
    let name: String
    let slug: String
    let description: String?
    let categoryName: String
    let categorySlug: String
    let address: String
    let latitude: Double
    let longitude: Double
    let phone: String?
    let email: String?
    let website: String?

    enum CodingKeys: String, CodingKey {
        case id
        case name
        case slug
        case description
        case categoryName = "category_name"
        case categorySlug = "category_slug"
        case address
        case latitude
        case longitude
        case phone
        case email
        case website
    }
}

// BusinessHour represents the operating hours for a business on a specific day.
// day_of_week: 0 = Sunday, 1 = Monday, ..., 6 = Saturday
struct BusinessHour: Codable, Sendable {
    let dayOfWeek: Int
    let openTime: String
    let closeTime: String
    let isClosed: Bool

    enum CodingKeys: String, CodingKey {
        case dayOfWeek = "day_of_week"
        case openTime = "open_time"
        case closeTime = "close_time"
        case isClosed = "is_closed"
    }
}

// BusinessDetails represents a single business with its hours and menus.
struct BusinessDetails: Codable, Identifiable, Sendable {
    let id: Int64
    let name: String
    let slug: String
    let description: String?
    let categoryName: String
    let categorySlug: String
    let address: String
    let latitude: Double
    let longitude: Double
    let phone: String?
    let email: String?
    let website: String?
    let hours: [BusinessHour]
    let menus: [Menu]

    enum CodingKeys: String, CodingKey {
        case id
        case name
        case slug
        case description
        case categoryName = "category_name"
        case categorySlug = "category_slug"
        case address
        case latitude
        case longitude
        case phone
        case email
        case website
        case hours
        case menus
    }
}

