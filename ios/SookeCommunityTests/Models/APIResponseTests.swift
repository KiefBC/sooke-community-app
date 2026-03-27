import Testing
import Foundation
@testable import SookeCommunity

@Suite("API Response Tests")
struct APIResponseTests {

    @Test func decodesPaginatedResponseOfBusiness() throws {
        let json = """
        {
            "items": [
                {
                    "id": 1,
                    "name": "Sooke Harbor House",
                    "slug": "sooke-harbor-house",
                    "description": "A waterfront restaurant",
                    "category_name": "Restaurants",
                    "category_slug": "restaurants",
                    "address": "1528 Whiffen Spit Rd",
                    "latitude": 48.3538,
                    "longitude": -123.7256,
                    "phone": "250-642-3421",
                    "email": "info@sookeharbourhouse.com",
                    "website": "https://sookeharbourhouse.com"
                }
            ],
            "pagination": {
                "page": 1,
                "per_page": 20,
                "total_items": 5,
                "total_pages": 1
            }
        }
        """
        let data = try #require(json.data(using: .utf8))
        let response = try JSONDecoder().decode(PaginatedResponse<Business>.self, from: data)

        #expect(response.items.count == 1)
        #expect(response.items[0].id == 1)
        #expect(response.items[0].name == "Sooke Harbor House")

        #expect(response.pagination.page == 1)
        #expect(response.pagination.perPage == 20)
        #expect(response.pagination.totalItems == 5)
        #expect(response.pagination.totalPages == 1)
    }

    @Test func decodesAPIErrorResponse() throws {
        let json = """
        {
            "error": {
                "code": "not_found",
                "message": "Business not found"
            }
        }
        """
        let data = try #require(json.data(using: .utf8))
        let errorResponse = try JSONDecoder().decode(APIErrorResponse.self, from: data)

        #expect(errorResponse.error.code == "not_found")
        #expect(errorResponse.error.message == "Business not found")
    }
}
