//
//  MapViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-04.
//

import Foundation
import CoreLocation


@MainActor
@Observable
final class MapViewModel {
    private let locationManager = CLLocationManager()
    private let apiClient: APIClient
    var businesses: [Business] = []
    var selectedCategory: Category? = nil
    var selectedBusiness: Business? = nil
    var isLoading: Bool = false
    var error: Error? = nil
    var filteredBusinesses: [Business] {
        guard let category = selectedCategory else { return businesses }
        return businesses.filter { $0.categorySlug == category.slug }
    }
    var categories: [Category] = []
    
    init(apiClient: APIClient) {
        self.apiClient = apiClient
    }
    
    func fetchBusinesses() async {
        isLoading = true
        error = nil
        do {
            let response: PaginatedResponse<Business> = try await apiClient.get("/api/v1/businesses")
            businesses = response.items
        } catch {
            self.error = error
        }
        isLoading = false
    }
    
    func fetchCategories() async {
        isLoading = true
        error = nil
        do {
            let response: CategoryListResponse = try await apiClient.get("/api/v1/categories")
            categories = response.items
        } catch {
            self.error = error
        }
        isLoading = false
    }
    
    func selectCategory(_ category: Category?) {
        selectedCategory = category
    }
    
    func requestLocationPermission() {
        locationManager.requestWhenInUseAuthorization()
    }
    
    // TODO: implement map view model to fetch business locations and details for map annotations
}
