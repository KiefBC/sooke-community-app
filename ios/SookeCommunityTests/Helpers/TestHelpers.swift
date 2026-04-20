//
//  TestHelpers.swift
//  SookeCommunityTests
//
//  Created by Kiefer Hay on 2026-04-06.
//

import Foundation
@testable import SookeCommunity

func makeTestClient() -> APIClient {
    let config = URLSessionConfiguration.ephemeral
    config.protocolClasses = [MockURLProtocol.self]
    let session = URLSession(configuration: config)
    return APIClient(baseURL: "http://localhost:8080", session: session)
}

func makeCategoryJSON(
    categories: [(id: Int64, name: String, slug: String)] = [
        (1, "Food", "food"),
        (2, "Retail", "retail"),
        (3, "Food", "food")
    ]
) -> Data {
    let items = categories.map { c in
        "{\"id\": \(c.id), \"name\": \"\(c.name)\", \"slug\": \"\(c.slug)\"}"
    }.joined(separator: ",")
                                                                                                                          
    return """
      {"items": [\(items)]}
      """.data(using: .utf8)!
}

func makeEventTypeJSON(
    eventTypes: [(id: Int64, name: String, slug: String)] = [
        (1, "Market", "market"),
        (2, "Music", "music"),
        (3, "Art", "art")
    ]
) -> Data {
    let items = eventTypes.map { e in
        "{\"id\": \(e.id), \"name\": \"\(e.name)\", \"slug\": \"\(e.slug)\"}"
    }.joined(separator: ",")

    return """
      {"items": [\(items)]}
      """.data(using: .utf8)!
}

func makeErrorJSON(code: String = "server_error", message: String = "Internal Server Error") -> Data {
      """
      {"error": {"code": "\(code)", "message": "\(message)"}}
      """.data(using: .utf8)!
}

func makePaginatedBusinessJSON(
    businesses: [(id: Int64, name: String, slug: String, categoryName: String, categorySlug: String)] = [
        (1, "Test Cafe", "test-cafe", "Food", "food"),
        (2, "Test Store", "test-store", "Retail", "retail"),
        (3, "Test Deli", "test-deli", "Food", "food")
    ],
    page: Int = 1,
    perPage: Int = 20,
    totalItems: Int? = nil,
    totalPages: Int = 1
) -> Data {
    let items = businesses.map { b in
        "{\"id\": \(b.id), \"name\": \"\(b.name)\", \"slug\": \"\(b.slug)\", \"description\": null, \"category_name\": \"\(b.categoryName)\", \"category_slug\": \"\(b.categorySlug)\", \"address\": \"123 Main St\", \"latitude\": 48.37, \"longitude\": -123.72, \"phone\": null, \"email\": null, \"website\": null}"
    }.joined(separator: ",")

    return """
    {"items": [\(items)], "pagination": {"page": \(page), "per_page": \(perPage), "total_items": \(totalItems ?? businesses.count), "total_pages": \(totalPages)}}
    """.data(using: .utf8)!
}
func makePaginatedEventJSON(
    events: [(id: Int64, name: String, slug: String, eventTypeName: String, eventTypeSlug: String, startTime: String, endTime: String, status: String, businessName: String?, businessSlug: String?)] = [
        (1, "Farmers Market", "farmers-market", "Market", "market", "2024-07-15T10:00:00Z", "2024-07-15T14:00:00Z", "upcoming", nil, nil),
        (2, "Live Music Night", "live-music-night", "Music", "music", "2024-07-16T19:00:00Z", "2024-07-16T22:00:00Z", "upcoming", nil, nil),
        (3, "Art Walk", "art-walk", "Art", "art", "2024-07-17T12:00:00Z", "2024-07-17T16:00:00Z", "upcoming", "Test Gallery", "test-gallery")
    ],
    page: Int = 1,
    perPage: Int = 20,
    totalItems: Int? = nil,
    totalPages: Int = 1
) -> Data {
    let items = events.map { e in
        let businessNameJSON = e.businessName.map { "\"\($0)\"" } ?? "null"
        let businessSlugJSON = e.businessSlug.map { "\"\($0)\"" } ?? "null"
        return "{\"id\": \(e.id), \"name\": \"\(e.name)\", \"slug\": \"\(e.slug)\", \"description\": null, \"event_type_name\": \"\(e.eventTypeName)\", \"event_type_slug\": \"\(e.eventTypeSlug)\", \"latitude\": null, \"longitude\": null, \"start_time\": \"\(e.startTime)\", \"end_time\": \"\(e.endTime)\", \"status\": \"\(e.status)\", \"business_name\": \(businessNameJSON), \"business_slug\": \(businessSlugJSON)}"
    }.joined(separator: ",")

    return """
    {"items": [\(items)], "pagination": {"page": \(page), "per_page": \(perPage), "total_items": \(totalItems ?? events.count), "total_pages": \(totalPages)}}
    """.data(using: .utf8)!
}

