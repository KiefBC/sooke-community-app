//
//  Event.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-20.
//

import Foundation

struct Event: Codable, Identifiable, Sendable, Hashable {
    let id: Int64
    let name: String
    let slug: String
    /// Can be null if the event doesn't have a description, so this is optional
    let description: String?
    let eventTypeName: String
    let eventTypeSlug: String
    /// Can be null if there is a business associated with the event, so these are optional
    let latitude: Double?
    /// Can be null if there is a business associated with the event, so these are optional
    let longitude: Double?
    /// RFC3339 format via Go time.Time
    let startTime: String
    /// RFC3339 format via Go time.Time
    let endTime: String
    let status: String
    /// Can be null if the event isn't associated with a business, so these are optional
    let businessName: String?
    /// Can be null if the event isn't associated with a business, so these are optional
    let businessSlug: String?
    
    enum CodingKeys: String, CodingKey {
        case id, name, slug, description, latitude, longitude, status
        case eventTypeName = "event_type_name"
        case eventTypeSlug = "event_type_slug"
        case startTime = "start_time"
        case endTime = "end_time"
        case businessName = "business_name"
        case businessSlug = "business_slug"
    }
}
