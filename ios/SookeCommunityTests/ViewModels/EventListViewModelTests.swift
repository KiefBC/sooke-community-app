//
//  EventListViewModelTests.swift
//  SookeCommunityTests
//
//  Created by Kiefer Hay on 2026-04-20.
//

import Foundation
import Testing
@testable import SookeCommunity

@Suite("Event List ViewModel Tests")
@MainActor
struct EventListViewModelTests {
    @Test func fetchesEvents() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedEventJSON()
        let vm = EventListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchEvents()
        
        #expect(vm.items.count == 3)
        #expect(vm.items.first?.name == "Farmers Market")
    }
    
    @Test func handlesError() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON(
            code: "server_error",
            message: "Internal Server Error BUG"
        )
        let vm = EventListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchEvents()
        
        #expect(vm.items.isEmpty)
        #expect(vm.error != nil)
    }
    
    @Test func fetchesEventTypes() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeEventTypeJSON()
        let vm = EventListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchEventTypes()
        
        #expect(vm.eventTypes.count == 3)
        #expect(vm.eventTypes[0].name == "Market")
        #expect(vm.eventTypes[1].slug == "music")
    }
    
    @Test func sendsEventTypeAsQueryParameter() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedEventJSON()
        let vm = EventListViewModel()
        vm.apiClient = makeTestClient()
        vm.selectedEventType = EventType(id: 1, name: "Market", slug: "market")
        await vm.fetchEvents()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.items.count == 3)
        #expect(url?.contains("type=market") == true)
    }
    
    @Test func clearsErrorBeforeFetch() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = EventListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchEventTypes()
        #expect(vm.error != nil)

        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        await vm.fetchEventTypes()
        #expect(vm.error == nil)
    }
}
