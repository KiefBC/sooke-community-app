//
//  BusinessDetailViewModelTests.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import Foundation
import Testing
@testable import SookeCommunity

@Suite("Business Detail ViewModel Tests")
@MainActor
struct BusinessDetailViewModelTests {
    func makeTestClient() -> APIClient {
        let config = URLSessionConfiguration.ephemeral
        config.protocolClasses = [MockURLProtocol.self]
        let session = URLSession(configuration: config)
        return APIClient(baseURL: "http://localhost:8080", session: session)
    }

    func makeBusinessDetailsJSON(
        id: Int64 = 1,
        name: String = "Test Cafe",
        slug: String = "test-cafe"
    ) -> Data {
        """
        {"id": \(id), "name": "\(name)", "slug": "\(slug)", "description": "A test cafe", "category_name": "Food", "category_slug": "food", "address": "123 Main St", "latitude": 48.37, "longitude": -123.72, "phone": null, "email": null, "website": null, "hours": [{"day_of_week": 1, "open_time": "09:00:00", "close_time": "17:00:00", "is_closed": false}], "menus": []}
        """.data(using: .utf8)!
    }

    func makeErrorJSON(code: String = "not_found", message: String = "Business not found") -> Data {
        """
        {"error": {"code": "\(code)", "message": "\(message)"}}
        """.data(using: .utf8)!
    }

    @Test func fetchesBusinessDetails() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeBusinessDetailsJSON()
        let vm = BusinessDetailViewModel(apiClient: makeTestClient())
        await vm.fetchBusinessDetails(slug: "test-cafe")

        #expect(vm.businessDetails != nil)
        #expect(vm.businessDetails?.name == "Test Cafe")
        #expect(vm.businessDetails?.hours.count == 1)
        #expect(vm.isLoading == false)
        #expect(vm.error == nil)
    }

    @Test func handlesError() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 404
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = BusinessDetailViewModel(apiClient: makeTestClient())
        await vm.fetchBusinessDetails(slug: "nonexistent")

        #expect(vm.businessDetails == nil)
        #expect(vm.error != nil)
        #expect(vm.isLoading == false)
    }

    @Test func clearsErrorBeforeFetch() async throws {
        MockURLProtocol.reset()
        MockURLProtocol.mockStatusCode = 404
        MockURLProtocol.mockResponseData = makeErrorJSON()
        let vm = BusinessDetailViewModel(apiClient: makeTestClient())
        await vm.fetchBusinessDetails(slug: "nonexistent")
        #expect(vm.error != nil)

        MockURLProtocol.reset()
        MockURLProtocol.mockResponseData = makeBusinessDetailsJSON()
        await vm.fetchBusinessDetails(slug: "test-cafe")
        #expect(vm.error == nil)
    }

    @Test func isLoadingDefaultsToFalse() async throws {
        let vm = BusinessDetailViewModel(apiClient: makeTestClient())
        #expect(vm.isLoading == false)
        #expect(vm.businessDetails == nil)
        #expect(vm.error == nil)
    }
}
