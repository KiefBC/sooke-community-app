import SwiftUI

struct ContentView: View {
    @Environment(ThemeManager.self) private var themeManager
    @State private var selectedTab: AppTab = .home

    var body: some View {
        TabView(selection: $selectedTab) {
            ForEach(AppTab.allCases, id: \.self) { tab in
                Group {
                    switch tab {
                    case .home: HomeView()
                    case .businesses: BusinessesPlaceholderView()
                    case .events: EventsPlaceholderView()
                    case .map: MapPlaceholderView()
                    }
                }
                .tabItem {
                    Label(tab.title, systemImage: tab.icon)
                }
                .tag(tab)
            }
        }
        .tint(themeManager.colors.primary)
    }
}

#Preview {
    ContentView()
        .environment(ThemeManager())
}
