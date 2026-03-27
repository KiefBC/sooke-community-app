import SwiftUI

enum AppTab: String, CaseIterable {
    case home
    case businesses
    case events
    case map

    var title: String {
        switch self {
        case .home: "Home"
        case .businesses: "Businesses"
        case .events: "Events"
        case .map: "Map"
        }
    }

    var icon: String {
        switch self {
        case .home: "house"
        case .businesses: "storefront"
        case .events: "calendar"
        case .map: "map"
        }
    }
}
