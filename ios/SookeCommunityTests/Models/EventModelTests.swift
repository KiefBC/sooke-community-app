//
//  EventModelTests.swift
//  SookeCommunityTests
//
//  Created by Kiefer Hay on 2026-04-20.
//

import Testing
import Foundation
@testable import SookeCommunity

@Suite("Event Model Tests")
struct EventModelTests {
    
    let seedData: String = """
    {
        "id": 1,
        "name": "Summer Music Festival",
        "slug": "summer-music-festival",
        "description": null,
        "event_type_name": "Festival",
        "event_type_slug": "festival",
        "latitude": 48.4284,
        "longitude": -123.3656,
        "start_time": "2024-07-15T10:00:00Z",
        "end_time": "2024-07-15T22:00:00Z",
        "status": "published",
        "business_name": "Downtown Events Co",
        "business_slug": "downtown-events-co"
    }
    """

    
    @Test func decodeEventWithAllFields() throws {
        
        let data = seedData.data(using: .utf8)!
        let event = try JSONDecoder().decode(Event.self, from: data)
        
        #expect(event.id == 1)
        #expect(event.name == "Summer Music Festival")
        #expect(event.slug == "summer-music-festival")
        #expect(event.description == nil)
        #expect(event.latitude == 48.4284)
        #expect(event.longitude == -123.3656)
        #expect(event.startTime == "2024-07-15T10:00:00Z")
        #expect(event.endTime == "2024-07-15T22:00:00Z")
        #expect(event.status == "published")
        #expect(event.businessName == "Downtown Events Co")
        #expect(event.businessSlug == "downtown-events-co")
    }
}
