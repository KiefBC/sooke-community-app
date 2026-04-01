import SwiftUI

@main
struct SookeCommunityApp: App {
    @State private var themeManager = ThemeManager()
    private let apiClient = APIClient(baseURL: APIConfig.baseURL)

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(themeManager)
                .environment(\.apiClient, apiClient)
                .preferredColorScheme(.dark)
        }
    }
}
