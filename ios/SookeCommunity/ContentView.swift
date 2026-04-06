import SwiftUI

struct ContentView: View {
    @Environment(ThemeManager.self) private var themeManager
    @Environment(\.apiClient) private var apiClient
    @State private var selectedTab: AppTab = .home

    var body: some View {
        TabView(selection: $selectedTab) {
            Tab(AppTab.home.title, systemImage: AppTab.home.icon, value: .home) {
                HomeView()
            }
            Tab(AppTab.businesses.title, systemImage: AppTab.businesses.icon, value: .businesses) {
                BusinessListView()
            }
            Tab(AppTab.events.title, systemImage: AppTab.events.icon, value: .events) {
                EventsPlaceholderView()
            }
            Tab(AppTab.map.title, systemImage: AppTab.map.icon, value: .map) {
                SookeMapView()
            }
        }
        .tint(themeManager.colors.primary)
        .tabBarMinimizeBehavior(.onScrollDown)
    }
}

#Preview {
    ContentView()
        .environment(ThemeManager())
}
