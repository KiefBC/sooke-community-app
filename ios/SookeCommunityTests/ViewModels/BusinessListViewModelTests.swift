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
    func makeTestClient() -> APIClient {
        let config = URLSessionConfiguration.ephemeral
        config.protocolClasses = [MockURLProtocol.self]
        let session = URLSession(configuration: config)
        return APIClient(baseURL: "http://localhost:8080", session: session)
    }
    
    func makeCategoryJSON(
        categories: [(id: Int64, name: String, slug: String)] = [
            (1, "Food", "food"),
            (2, "Retail", "retail")
        ]
    ) -> Data {
        let items = categories.map { c in
            "{\"id\": \(c.id), \"name\": \"\(c.name)\", \"slug\": \"\(c.slug)\"}"
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
            (1, "Test Cafe", "test-cafe", "Food", "food")
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
    
    @Test func fetchesBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        await vm.fetchBusinesses()
        
        #expect(vm.items.count == 1)
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
        let vm = BusinessListViewModel(apiClient: makeTestClient())
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
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        await vm.fetchCategories()
        
        #expect(vm.categories.isEmpty)
        #expect(vm.items.isEmpty)
        #expect(vm.error != nil)
        #expect(vm.isLoading == false)
    }
    
    @Test func fetchesCategories() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeCategoryJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        await vm.fetchCategories()
        
        #expect(vm.categories.count == 2)
        #expect(vm.categories[0].name == "Food")
        #expect(vm.categories[1].name == "Retail")
        #expect(vm.isLoading == false)
    }
    
    @Test func selectCategoryUpdatesState() async throws {
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        vm.selectCategory(Category(id: 1, name: "Food", slug: "food"))
        #expect(vm.selectedCategory?.name == "Food")
        
        vm.selectCategory(nil)
        #expect(vm.selectedCategory == nil)
        
    }
    
    @Test func searchesBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        vm.searchText = "Cafe"
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.items.count == 1)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("search=Cafe") == true)
    }
    
    @Test func filtersByCategory() async   {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        vm.selectedCategory = Category(id: 1, name: "Food", slug: "food")
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.items.count == 1)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("category=food") == true)
    }
    
    @Test func selectCategoryFiltersBusinesses() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        let testCategory = Category(id: 1, name: "Food", slug: "food")
        vm.selectCategory(testCategory)
        await vm.fetchBusinesses()
        
        let url = MockURLProtocol.lastRequest?.url?.absoluteString
        
        #expect(vm.selectedCategory == testCategory)
        #expect(vm.items.count == 1)
        #expect(vm.items.first?.name == "Test Cafe")
        #expect(vm.isLoading == false)
        #expect(url?.contains("category=food") == true)
    }
    
    @Test func clearingCategoryFetchesAll() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
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
    }
    
    @Test func isLoadingDefaultsToFalse() async throws {
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        #expect(vm.isLoading == false)
        #expect(vm.isLoadingBusinesses == false)
        #expect(vm.isLoadingCategories == false)
    }

    @Test func clearsErrorBeforeFetch() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 500
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = BusinessListViewModel(apiClient: makeTestClient())
        await vm.fetchBusinesses()
        #expect(vm.error != nil)

        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makePaginatedBusinessJSON()
        await vm.fetchBusinesses()
        #expect(vm.error == nil)
    }
}
