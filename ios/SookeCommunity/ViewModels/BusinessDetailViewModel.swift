//
//  BusinessDetailViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import Foundation

@MainActor
@Observable
final class BusinessDetailViewModel {
    private let apiClient: APIClient
    private(set) var businessDetails: BusinessDetails?
    private(set) var isLoading = false
    private(set) var error: Error?

    init(apiClient: APIClient) {
        self.apiClient = apiClient
    }

    func fetchBusinessDetails(slug: String) async {
        isLoading = true
        error = nil

        do {
            businessDetails = try await apiClient.get("/api/v1/businesses/\(slug)")
        } catch {
            self.error = error
        }

        isLoading = false
    }
}
