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
    var apiClient: APIClient?
    private(set) var businessDetails: BusinessDetails?
    private(set) var isLoading = false
    private(set) var error: Error?

    func fetchBusinessDetails(slug: String) async {
        guard let apiClient else { return }
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
