//
//  MapViewModelTests.swift
//  SookeCommunityTests
//
//  Created by Kiefer Hay on 2026-04-04.
//

import Foundation
import Testing
@testable import SookeCommunity

@Suite("Map ViewModel Tests")
@MainActor
struct MapViewModelTests {
    @Test func fetchesBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        // TODO: Seriously consider locations NEAR the user in the main View this is an idea for now
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString

        #expect(vm.businesses.count == 3)
        #expect(vm.businesses[0].name == "Test Cafe")
        #expect(vm.businesses[2].categoryName == "Food")
        #expect(url?.contains("tz=") == true)
    }

    @Test func fetchesCategories() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeCategoryJSON()
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchCategories()

        #expect(vm.categories.count == 3)
        #expect(vm.categories[0].name == "Food")
        #expect(vm.categories[1].slug == "retail")
    }

    @Test func selectCategoryFiltersLocally() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        vm.selectedCategory = Category(id: 1, name: "Food", slug: "food")
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.filteredBusinesses.count == 2)
        #expect(vm.businesses.count == 3)
        #expect(url?.contains("tz=") == true)
    }
    
    @Test func nilCategoryShowsAllLocations() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        vm.selectedCategory = nil
        await vm.fetchBusinesses()
        
        #expect(vm.filteredBusinesses.count == 3)
        #expect(vm.businesses.count == 3)
        #expect(vm.selectedCategory == nil)
        
        vm.selectedCategory = Category(id: 1, name: "Food", slug: "food")
        #expect(vm.filteredBusinesses.count == 2)
        
        vm.selectedCategory = nil
        #expect(vm.filteredBusinesses.count == 3)
    }
    
    @Test func selectedBusinessUpdatesState() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()

        let business = vm.businesses[1]
        vm.selectedBusiness = business
        
        #expect(vm.selectedBusiness?.id == business.id)
        #expect(vm.selectedBusiness?.name == business.name)
        
        vm.selectedBusiness = nil
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.selectedBusiness == nil)
        #expect(vm.businesses.count == 3)
        #expect(vm.filteredBusinesses.count == 3)
        #expect(url?.contains("tz=") == true)
    }
    
    @Test func handlesError() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = MapViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()

        #expect(vm.error != nil)
        #expect(vm.businesses.isEmpty)
        #expect(vm.categories.isEmpty)
    }
}
