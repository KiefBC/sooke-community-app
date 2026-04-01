//
//  APIClientEnvironment.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import SwiftUI

enum APIConfig {
    static let baseURL = "http://localhost:8080"
}

private struct APIClientKey: EnvironmentKey {
    static let defaultValue = APIClient(baseURL: APIConfig.baseURL)
}

extension EnvironmentValues {
    var apiClient: APIClient {
        get { self[APIClientKey.self] }
        set { self[APIClientKey.self] = newValue }
    }
}
