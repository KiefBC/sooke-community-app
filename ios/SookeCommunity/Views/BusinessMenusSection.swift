import SwiftUI

struct BusinessMenusSection: View {
    @Environment(ThemeManager.self) private var themeManager
    let menus: [Menu]

    var body: some View {
        VStack(alignment: .leading, spacing: 16) {
            ForEach(menus) { menu in
                VStack(alignment: .leading, spacing: 12) {
                    Text(menu.name)
                        .font(.headline)
                        .foregroundStyle(themeManager.colors.accent)

                    if let description = menu.description, !description.isEmpty {
                        Text(description)
                            .font(.subheadline)
                            .foregroundStyle(.secondary)
                    }

                    VStack(alignment: .leading, spacing: 8) {
                        ForEach(menu.items) { item in
                            HStack(alignment: .top) {
                                VStack(alignment: .leading, spacing: 2) {
                                    Text(item.name)
                                        .font(.subheadline)
                                        .foregroundStyle(.primary)

                                    if let description = item.description, !description.isEmpty {
                                        Text(description)
                                            .font(.caption)
                                            .foregroundStyle(.secondary)
                                    }
                                }

                                Spacer()

                                Text("$\(item.price)")
                                    .font(.subheadline)
                                    .foregroundStyle(.primary)
                            }
                        }
                    }
                }
                .padding()
                .background(
                    RoundedRectangle(cornerRadius: 12)
                        .fill(themeManager.colors.accent.opacity(0.1))
                )
            }
        }
        .padding(.horizontal)
    }
}
