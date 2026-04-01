import Foundation

// Category represents a business category returned by the API.
struct Category: Codable, Identifiable, Sendable, Equatable {
    let id: Int64
    let name: String
    let slug: String
}

// CategoryListResponse wraps the list of categories returned by GET /api/v1/categories.
struct CategoryListResponse: Codable, Sendable {
    let items: [Category]
}
