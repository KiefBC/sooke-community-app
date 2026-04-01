//
//  BusinessListViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-27.
//

import Foundation

@MainActor
@Observable
final class BusinessListViewModel {
    private let apiClient: APIClient
    var items: [Business] = []
    var categories: [Category] = []
    var selectedCategory: Category? = nil
    var searchText: String = ""
    private(set) var isLoadingBusinesses: Bool = false
    private(set) var isLoadingCategories: Bool = false
    var error: Error? = nil

    var isLoading: Bool {
        isLoadingBusinesses || isLoadingCategories
    }

    init(apiClient: APIClient) {
        self.apiClient = apiClient
    }

    func fetchBusinesses() async {
        isLoadingBusinesses = true
        error = nil
        var queryItems: [URLQueryItem] = []
        if !searchText.isEmpty {
            queryItems.append(URLQueryItem(name: "search", value: searchText))
        }
        if let category = selectedCategory {
            queryItems.append(URLQueryItem(name: "category", value: category.slug))
        }
        do {
            let response: PaginatedResponse<Business> = try await apiClient.get("/api/v1/businesses", queryItems: queryItems)
            items = response.items
        } catch {
            self.error = error
        }
        isLoadingBusinesses = false
    }

    func fetchCategories() async {
        isLoadingCategories = true
        error = nil
        do {
            let response: CategoryListResponse = try await apiClient.get("/api/v1/categories")
            categories = response.items
        } catch {
            self.error = error
        }
        isLoadingCategories = false
    }

    func selectCategory(_ category: Category?) {
        selectedCategory = category
    }
}
