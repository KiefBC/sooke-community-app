import Foundation

// Menu represents a menu associated with a business.
struct Menu: Codable, Identifiable, Sendable {
    let id: Int64
    let name: String
    let description: String?
    let items: [MenuItem]
}

// MenuItem represents an individual item on a menu.
// price is a String because the Go API returns Postgres NUMERIC as a string to avoid float rounding.
struct MenuItem: Codable, Identifiable, Sendable {
    let id: Int64
    let name: String
    let description: String?
    let price: String
}
