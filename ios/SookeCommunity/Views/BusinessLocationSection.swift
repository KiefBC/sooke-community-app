import SwiftUI

struct BusinessLocationSection: View {
    @Environment(ThemeManager.self) private var themeManager
    let address: String

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Location")
                .font(.headline)
                .foregroundStyle(themeManager.colors.accent)

            HStack(spacing: 12) {
                Image(systemName: "mappin.and.ellipse")
                    .foregroundStyle(themeManager.colors.accent)
                Text(address)
                    .foregroundStyle(.primary)
            }

            // placeholder until MapLibre integration
            RoundedRectangle(cornerRadius: 12)
                .fill(themeManager.colors.accent.opacity(0.2))
                .frame(height: 200)
                .overlay {
                    VStack {
                        Image(systemName: "map")
                            .font(.system(size: 40))
                            .foregroundStyle(themeManager.colors.accent)
                        Text("Map coming soon")
                            .font(.caption)
                            .foregroundStyle(.secondary)
                    }
                }
        }
        .padding(.horizontal)
    }
}
