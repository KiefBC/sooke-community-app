import Testing
import Foundation
@testable import SookeCommunity

@Suite("Theme Tests")
@MainActor
struct ThemeTests {

    @Test func allThemesExist() {
        #expect(SookeTheme.allCases.count == 5)
    }

    @Test func defaultThemeIsMidnightStrait() {
        #expect(SookeTheme.default == .midnightStrait)
    }

    @Test func eachThemeHasUniqueName() {
        let names = SookeTheme.allCases.map { $0.displayName }
        let uniqueNames = Set(names)
        #expect(names.count == uniqueNames.count)
    }

    @Test func eachThemeHasColors() {
        for theme in SookeTheme.allCases {
            let colors = theme.colors
            // Verify these properties are accessible (non-nil via struct existence)
            _ = colors.background
            _ = colors.foreground
            _ = colors.card
            _ = colors.primary
            _ = colors.accent
        }
        #expect(SookeTheme.allCases.count == 5)
    }

    @Test func themeManagerDefaultsToMidnightStrait() {
        // Use a fresh UserDefaults suite so tests are isolated
        let defaults = UserDefaults(suiteName: "ThemeTests_default")!
        defaults.removeObject(forKey: "selectedTheme")
        let manager = ThemeManager(defaults: defaults)
        #expect(manager.current == .midnightStrait)
    }

    @Test func themeManagerSwitchesTheme() {
        let defaults = UserDefaults(suiteName: "ThemeTests_switch")!
        defaults.removeObject(forKey: "selectedTheme")
        let manager = ThemeManager(defaults: defaults)
        manager.current = .rainforest
        #expect(manager.current == .rainforest)
    }
}
