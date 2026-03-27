import SwiftUI

struct BusinessesPlaceholderView: View {
    @Environment(ThemeManager.self) private var themeManager

    var body: some View {
        NavigationStack {
            ZStack {
                themeManager.colors.background.ignoresSafeArea()
                Text("Businesses coming in Milestone 3")
                    .foregroundStyle(themeManager.colors.muted)
            }
            .navigationTitle("Businesses")
            .toolbarColorScheme(.dark, for: .navigationBar)
        }
    }
}
