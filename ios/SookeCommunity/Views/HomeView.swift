import SwiftUI

struct HomeView: View {
    @Environment(ThemeManager.self) private var themeManager

    var body: some View {
        NavigationStack {
            ZStack {
                themeManager.colors.background.ignoresSafeArea()
                VStack(spacing: 16) {
                    Text("Sooke Community")
                        .font(.largeTitle.bold())
                        .foregroundStyle(themeManager.colors.foreground)
                    Text("Your local guide to Sooke, BC")
                        .font(.subheadline)
                        .foregroundStyle(themeManager.colors.muted)
                }
            }
            .navigationTitle("Home")
            .toolbarColorScheme(.dark, for: .navigationBar)
        }
    }
}
