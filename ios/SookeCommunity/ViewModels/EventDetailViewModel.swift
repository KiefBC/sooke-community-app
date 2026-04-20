//
//  EventDetailViewModel.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-20.
//

import Foundation

@MainActor
@Observable
final class EventDetailViewModel {
    var apiClient: APIClient?
    private(set) var eventDetails: EventDetails?
    private(set) var isLoading = false
    private(set) var error: Error?

    func fetchEventDetails(slug: String) async {
        guard let apiClient else { return }
        isLoading = true
        error = nil

        do {
            eventDetails = try await apiClient.get("/api/v1/events/\(slug)")
        } catch {
            self.error = error
        }

        isLoading = false
    }
}



