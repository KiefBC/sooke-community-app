import SwiftUI

// MARK: - ThemeColors

struct ThemeColors: Sendable {
    let background: Color
    let foreground: Color
    let card: Color
    let primary: Color
    let accent: Color
    let secondary: Color
    let muted: Color
    let border: Color

    // Semantic status colors for hours display
    let statusOpen = Color(red: 0.30, green: 0.75, blue: 0.40)
    let statusClosed = Color(red: 0.80, green: 0.35, blue: 0.35)
    let statusSoon = Color(red: 0.85, green: 0.65, blue: 0.25)
}

// MARK: - SookeTheme

enum SookeTheme: String, CaseIterable, Sendable {
    case midnightStrait
    case rainforest
    case tidalPool
    case duskShore
    case cedarFog

    static let `default`: SookeTheme = .midnightStrait

    var displayName: String {
        switch self {
        case .midnightStrait: return "Midnight Strait"
        case .rainforest:     return "Rainforest"
        case .tidalPool:      return "Tidal Pool"
        case .duskShore:      return "Dusk Shore"
        case .cedarFog:       return "Cedar & Fog"
        }
    }

    // Approximate sRGB conversions from OKLCh originals
    var colors: ThemeColors {
        switch self {
        case .midnightStrait:
            return ThemeColors(
                background: Color(red: 0.10, green: 0.14, blue: 0.20),
                foreground: Color(red: 0.88, green: 0.90, blue: 0.92),
                card:       Color(red: 0.13, green: 0.18, blue: 0.25),
                primary:    Color(red: 0.30, green: 0.62, blue: 0.65),
                accent:     Color(red: 0.72, green: 0.60, blue: 0.30),
                secondary:  Color(red: 0.18, green: 0.23, blue: 0.30),
                muted:      Color(red: 0.15, green: 0.20, blue: 0.27),
                border:     Color(red: 0.20, green: 0.26, blue: 0.33)
            )
        case .rainforest:
            return ThemeColors(
                background: Color(red: 0.08, green: 0.15, blue: 0.10),
                foreground: Color(red: 0.88, green: 0.91, blue: 0.88),
                card:       Color(red: 0.11, green: 0.19, blue: 0.13),
                primary:    Color(red: 0.25, green: 0.58, blue: 0.38),
                accent:     Color(red: 0.70, green: 0.58, blue: 0.28),
                secondary:  Color(red: 0.15, green: 0.22, blue: 0.16),
                muted:      Color(red: 0.13, green: 0.20, blue: 0.14),
                border:     Color(red: 0.18, green: 0.26, blue: 0.19)
            )
        case .tidalPool:
            return ThemeColors(
                background: Color(red: 0.08, green: 0.14, blue: 0.16),
                foreground: Color(red: 0.88, green: 0.90, blue: 0.91),
                card:       Color(red: 0.11, green: 0.18, blue: 0.20),
                primary:    Color(red: 0.20, green: 0.58, blue: 0.60),
                accent:     Color(red: 0.72, green: 0.55, blue: 0.28),
                secondary:  Color(red: 0.14, green: 0.22, blue: 0.24),
                muted:      Color(red: 0.12, green: 0.19, blue: 0.21),
                border:     Color(red: 0.17, green: 0.25, blue: 0.27)
            )
        case .duskShore:
            return ThemeColors(
                background: Color(red: 0.09, green: 0.11, blue: 0.19),
                foreground: Color(red: 0.88, green: 0.88, blue: 0.92),
                card:       Color(red: 0.12, green: 0.14, blue: 0.24),
                primary:    Color(red: 0.35, green: 0.45, blue: 0.72),
                accent:     Color(red: 0.75, green: 0.55, blue: 0.35),
                secondary:  Color(red: 0.15, green: 0.17, blue: 0.28),
                muted:      Color(red: 0.13, green: 0.15, blue: 0.24),
                border:     Color(red: 0.18, green: 0.20, blue: 0.32)
            )
        case .cedarFog:
            return ThemeColors(
                background: Color(red: 0.12, green: 0.13, blue: 0.11),
                foreground: Color(red: 0.89, green: 0.89, blue: 0.87),
                card:       Color(red: 0.16, green: 0.17, blue: 0.14),
                primary:    Color(red: 0.40, green: 0.52, blue: 0.35),
                accent:     Color(red: 0.68, green: 0.55, blue: 0.32),
                secondary:  Color(red: 0.18, green: 0.20, blue: 0.17),
                muted:      Color(red: 0.16, green: 0.17, blue: 0.15),
                border:     Color(red: 0.22, green: 0.24, blue: 0.20)
            )
        }
    }
}

// MARK: - ThemeManager

@MainActor @Observable final class ThemeManager {
    private static let userDefaultsKey = "selectedTheme"
    private let defaults: UserDefaults

    var current: SookeTheme {
        didSet {
            defaults.set(current.rawValue, forKey: Self.userDefaultsKey)
        }
    }

    var colors: ThemeColors {
        current.colors
    }

    init(defaults: UserDefaults = .standard) {
        self.defaults = defaults
        if let raw = defaults.string(forKey: Self.userDefaultsKey),
           let saved = SookeTheme(rawValue: raw) {
            self.current = saved
        } else {
            self.current = SookeTheme.default
        }
    }
}
