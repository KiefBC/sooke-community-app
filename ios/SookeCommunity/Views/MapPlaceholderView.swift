import SwiftUI

struct MapPlaceholderView: View {
    @Environment(ThemeManager.self) private var themeManager

    var body: some View {
        NavigationStack {
            ZStack {
                themeManager.colors.background.ignoresSafeArea()
                Text("Map coming in Milestone 4")
                    .foregroundStyle(themeManager.colors.muted)
            }
            .navigationTitle("Map")
            .toolbarColorScheme(.dark, for: .navigationBar)
        }
    }
}
