//
//  EventListViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-20.
//

import Foundation

@MainActor
@Observable
final class EventListViewModel {
    var apiClient: APIClient?
    var items: [EventDetails] = []
    var eventTypes: [EventType] = []
    var selectedEventType: EventType? = nil
    var searchText: String = ""
    private(set) var isLoadingEvent: Bool = false
    private(set) var isLoadingEventType: Bool = false
    var error: Error? = nil
    
    var isLoading: Bool {
        isLoadingEvent || isLoadingEventType
    }
    
    func fetchEvents() async {
        guard let apiClient else { return }
        isLoadingEvent = true
        error = nil
        var queryItems: [URLQueryItem] = []
        
        if !searchText.isEmpty {
            queryItems.append(URLQueryItem(name: "search", value: searchText))
        }
        if let eventType = selectedEventType {
            queryItems.append(URLQueryItem(name: "type", value: eventType.slug))
        }
        
        do {
            let response: PaginatedResponse<EventDetails> = try await apiClient.get("/api/v1/events", queryItems: queryItems)
            items = response.items
        } catch {
            self.error = error
        }
        isLoadingEvent = false
    }
    
    func fetchEventTypes() async {
        guard let apiClient else { return }
        isLoadingEventType = true
        error = nil
        
        do {
            let response: ListResponse<EventType> = try await apiClient.get("/api/v1/event-types")
            eventTypes = response.items
        } catch {
            self.error = error
        }
        isLoadingEventType = false
    }
    
    func selectEventType(_ eventType: EventType?) {
        selectedEventType = eventType
    }
}
