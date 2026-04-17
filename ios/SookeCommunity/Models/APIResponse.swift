import Foundation

// PaginatedResponse wraps a list of items with pagination metadata.
// Matches the Go handler.PaginatedResponse[T] generic type.
struct PaginatedResponse<T: Codable & Sendable>: Codable, Sendable {
    let items: [T]
    let pagination: Pagination
}

// ListResponse wraps a list of items without pagination, used for small
// reference lookups like categories and event types.
// Matches the Go handler.ListResponse[T] generic type.
struct ListResponse<T: Codable & Sendable>: Codable, Sendable {
    let items: [T]
}

// Pagination contains page metadata included in every paginated API response.
struct Pagination: Codable, Sendable {
    let page: Int
    let perPage: Int
    let totalItems: Int
    let totalPages: Int

    enum CodingKeys: String, CodingKey {
        case page
        case perPage = "per_page"
        case totalItems = "total_items"
        case totalPages = "total_pages"
    }
}

// APIErrorResponse represents an error returned by the API.
struct APIErrorResponse: Codable, Sendable {
    let error: APIErrorDetail
}

// APIErrorDetail contains the code and human-readable message for an API error.
struct APIErrorDetail: Codable, Sendable {
    let code: String
    let message: String
}
