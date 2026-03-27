import SwiftUI

struct EventsPlaceholderView: View {
    @Environment(ThemeManager.self) private var themeManager

    var body: some View {
        NavigationStack {
            ZStack {
                themeManager.colors.background.ignoresSafeArea()
                Text("Events coming in Milestone 5")
                    .foregroundStyle(themeManager.colors.muted)
            }
            .navigationTitle("Events")
            .toolbarColorScheme(.dark, for: .navigationBar)
        }
    }
}
