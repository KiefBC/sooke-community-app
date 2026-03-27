import SwiftUI

@main
struct SookeCommunityApp: App {
    @State private var themeManager = ThemeManager()

    var body: some Scene {
        WindowGroup {
            ContentView()
                .environment(themeManager)
                .preferredColorScheme(.dark)
        }
    }
}
