import Foundation

// EventType represents an event category returned by GET /api/v1/event-types.
struct EventType: Codable, Identifiable, Sendable, Equatable {
    let id: Int64
    let name: String
    let slug: String
}

struct EventTypeListResponse: Codable, Sendable {
    let items: [EventType]
}
