import Testing
import Foundation
@testable import SookeCommunity

@Suite("APIClient Tests")
struct APIClientTests {

    @Test func fetchesBusinessList() async throws {
        MockURLProtocol.reset()

        let payload = """
        {
            "items": [
                {
                    "id": 1,
                    "name": "Test Cafe",
                    "slug": "test-cafe",
                    "description": "A cozy cafe",
                    "category_name": "Food",
                    "category_slug": "food",
                    "address": "123 Main St",
                    "latitude": 48.3686,
                    "longitude": -123.7210,
                    "phone": null,
                    "email": null,
                    "website": null
                }
            ],
            "pagination": {
                "page": 1,
                "per_page": 20,
                "total_items": 1,
                "total_pages": 1
            }
        }
        """.data(using: .utf8)!

        MockURLProtocol.mockResponseData = payload
        MockURLProtocol.mockStatusCode = 200

        let client = makeTestClient()
        let result: PaginatedResponse<Business> = try await client.get("/api/v1/businesses")

        #expect(result.items.count == 1)
        #expect(result.items[0].name == "Test Cafe")
        #expect(result.items[0].slug == "test-cafe")
        #expect(result.pagination.page == 1)
        #expect(result.pagination.totalItems == 1)
    }

    @Test func fetchesSingleBusiness() async throws {
        MockURLProtocol.reset()

        let payload = """
        {
            "id": 42,
            "name": "Sooke Harbor",
            "slug": "sooke-harbor",
            "description": "Waterfront dining",
            "category_name": "Restaurants",
            "category_slug": "restaurants",
            "address": "1 Wharf St",
            "latitude": 48.3686,
            "longitude": -123.7210,
            "phone": "250-555-1234",
            "email": "info@sookeharbor.ca",
            "website": "https://sookeharbor.ca",
            "hours": [
                {
                    "day_of_week": 1,
                    "open_time": "09:00",
                    "close_time": "17:00",
                    "is_closed": false
                }
            ],
            "menus": []
        }
        """.data(using: .utf8)!

        MockURLProtocol.mockResponseData = payload
        MockURLProtocol.mockStatusCode = 200

        let client = makeTestClient()
        let result: BusinessDetails = try await client.get("/api/v1/businesses/sooke-harbor")

        #expect(result.id == 42)
        #expect(result.name == "Sooke Harbor")
        #expect(result.slug == "sooke-harbor")
        #expect(result.phone == "250-555-1234")
        #expect(result.hours.count == 1)
        #expect(result.hours[0].dayOfWeek == 1)
        #expect(result.menus.isEmpty)
    }

    @Test func throwsOnHTTPError() async throws {
        MockURLProtocol.reset()

        let errorPayload = """
        {
            "error": {
                "code": "not_found",
                "message": "Business not found"
            }
        }
        """.data(using: .utf8)!

        MockURLProtocol.mockResponseData = errorPayload
        MockURLProtocol.mockStatusCode = 404

        let client = makeTestClient()

        do {
            let _: BusinessDetails = try await client.get("/api/v1/businesses/nonexistent")
            Issue.record("Expected an error to be thrown but none was")
        } catch let error as APIError {
            if case .httpError(let statusCode, let response) = error {
                #expect(statusCode == 404)
                #expect(response?.error.code == "not_found")
                #expect(response?.error.message == "Business not found")
            } else {
                Issue.record("Expected APIError.httpError but got \(error)")
            }
        }
    }

    @Test func throwsOnNetworkError() async throws {
        MockURLProtocol.reset()

        MockURLProtocol.mockError = URLError(.notConnectedToInternet)

        let client = makeTestClient()

        do {
            let _: PaginatedResponse<Business> = try await client.get("/api/v1/businesses")
            Issue.record("Expected an error to be thrown but none was")
        } catch let error as URLError {
            #expect(error.code == .notConnectedToInternet)
        }
    }
}
