//
//  BusinessListViewModelTests.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-27.
//

import Foundation
import Testing
@testable import SookeCommunity

@Suite("Business List ViewModel Tests")
@MainActor
struct BusinessListViewModelTests {
    @Test func fetchesBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()

        #expect(vm.items.count == 3)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
    }

    @Test func handlesError() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON(
            code: "server_error",
            message: "Internal Server Error BUG"
        )
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()

        #expect(vm.items.isEmpty)
        #expect(vm.error != nil)
        #expect(vm.isLoading == false)
    }

    @Test func categoryFetchHandlesError() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON(
            code: "server_error",
            message: "Internal Server Error BUG"
        )
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchCategories()
        
        #expect(vm.categories.isEmpty)
        #expect(vm.items.isEmpty)
        #expect(vm.error != nil)
        #expect(vm.isLoading == false)
    }
    
    @Test func fetchesCategories() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeCategoryJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchCategories()

        #expect(vm.categories.count == 3)
        #expect(vm.categories[0].name == "Food")
        #expect(vm.categories[1].name == "Retail")
        #expect(vm.isLoading == false)
    }
    
    @Test func selectCategoryUpdatesState() async throws {
        let vm = BusinessListViewModel()
        vm.selectCategory(Category(id: 1, name: "Food", slug: "food"))
        #expect(vm.selectedCategory?.name == "Food")
        
        vm.selectCategory(nil)
        #expect(vm.selectedCategory == nil)
        
    }
    
    @Test func searchesBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        vm.searchText = "Cafe"
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.items.count == 3)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("search=Cafe") == true)
        #expect(url?.contains("tz=") == true)
    }
    
    @Test func filtersByCategory() async   {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        vm.selectedCategory = Category(id: 1, name: "Food", slug: "food")
        await vm.fetchBusinesses()

        let url = MockURLProtocol.lastRequest?.url?.absoluteString

        #expect(vm.items.count == 3)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("category=food") == true)
        #expect(url?.contains("tz=") == true)
    }

    @Test func selectCategoryFiltersBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        let testCategory = Category(id: 1, name: "Food", slug: "food")
        vm.selectCategory(testCategory)
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.selectedCategory == testCategory)
        #expect(vm.items.count == 3)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("category=food") == true)
        #expect(url?.contains("tz=") == true)
    }
    
    @Test func clearingCategoryFetchesAll() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        vm.selectedCategory = Category(id: 1, name: "Food", slug: "food")
        await vm.fetchBusinesses()

        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON(businesses: [
            (1, "Test Retail", "test-retail", "Food", "food"),
            (2, "Test Cafe", "test-cafe", "Retail", "retail")
        ])
        
        vm.selectCategory(nil)
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.selectedCategory == nil)
        #expect(vm.items.count == 2)
        #expect(vm.items.first?.name == "Test Retail")
        #expect(vm.isLoading == false)
        #expect(url?.contains("category=") == false)
        #expect(url?.contains("tz=") == true)
    }
    
    @Test func isLoadingDefaultsToFalse() async throws {
        let vm = BusinessListViewModel()
        #expect(vm.isLoading == false)
        #expect(vm.isLoadingBusinesses == false)
        #expect(vm.isLoadingCategories == false)
    }

    @Test func clearsErrorBeforeFetch() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = BusinessListViewModel()
        vm.apiClient = makeTestClient()
        await vm.fetchBusinesses()
        #expect(vm.error != nil)

        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        await vm.fetchBusinesses()
        #expect(vm.error == nil)
    }
}
