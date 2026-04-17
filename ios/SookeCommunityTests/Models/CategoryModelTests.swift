import Testing
import Foundation
@testable import SookeCommunity

@Suite("Category Model Tests")
struct CategoryModelTests {

    @Test func decodesCategory() throws {
        let json = """
        {"id": 1, "name": "Restaurants", "slug": "restaurants"}
        """
        let data = try #require(json.data(using: .utf8))
        let category = try JSONDecoder().decode(Category.self, from: data)

        #expect(category.id == 1)
        #expect(category.name == "Restaurants")
        #expect(category.slug == "restaurants")
    }

    @Test func decodesCategoryListResponse() throws {
        let json = """
        {
            "items": [
                {"id": 1, "name": "Restaurants", "slug": "restaurants"},
                {"id": 2, "name": "Retail", "slug": "retail"}
            ]
        }
        """
        let data = try #require(json.data(using: .utf8))
        let response = try JSONDecoder().decode(ListResponse<SookeCommunity.Category>.self, from: data)

        #expect(response.items.count == 2)
        #expect(response.items[0].id == 1)
        #expect(response.items[0].name == "Restaurants")
        #expect(response.items[0].slug == "restaurants")
        #expect(response.items[1].id == 2)
        #expect(response.items[1].name == "Retail")
        #expect(response.items[1].slug == "retail")
    }
}
