//
//  MapViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-04.
//

import CoreLocation
import Foundation

@MainActor
@Observable
final class MapViewModel {
    private let locationManager = CLLocationManager()
    var apiClient: APIClient?
    var businesses: [Business] = []
    var selectedCategory: Category? = nil
    var selectedBusiness: Business? = nil

    // TODO: Still not sure if I will need these state machines but keeping in case
    var isLoadingBusinesses: Bool = false
    var isLoadingCategories: Bool = false

    var timeZone: TimeZone = .current
    var error: Error? = nil
    var filteredBusinesses: [Business] {
        guard let category = selectedCategory else { return businesses }
        return businesses.filter { $0.categorySlug == category.slug }
    }
    var categories: [Category] = []

    func fetchBusinesses() async {
        guard let apiClient else { return }
        isLoadingBusinesses = true
        error = nil

        var queryItems: [URLQueryItem] = []
        queryItems.append(URLQueryItem(name: "tz", value: timeZone.identifier))

        do {
            let response: PaginatedResponse<Business> = try await apiClient.get(
                "/api/v1/businesses", queryItems: queryItems)
            businesses = response.items
        } catch {
            self.error = error
        }
        isLoadingBusinesses = false
    }

    func fetchCategories() async {
        guard let apiClient else { return }
        isLoadingCategories = true
        error = nil
        do {
            let response: ListResponse<Category> = try await apiClient.get("/api/v1/categories")
            categories = response.items
        } catch {
            self.error = error
        }
        isLoadingCategories = false
    }

    func selectCategory(_ category: Category?) {
        selectedCategory = category
    }

    func requestLocationPermission() {
        locationManager.requestWhenInUseAuthorization()
    }

    // TODO: implement map view model to fetch business locations and details for map annotations
}
