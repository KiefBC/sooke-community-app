import SwiftUI

struct BusinessContactSection: View {
    @Environment(ThemeManager.self) private var themeManager
    let business: Business

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Contact")
                .font(.headline)
                .foregroundStyle(themeManager.colors.accent)

            if let phone = business.phone,
               let encoded = phone.addingPercentEncoding(withAllowedCharacters: .urlQueryAllowed),
               let url = URL(string: "tel:\(encoded)") {
                Link(destination: url) {
                    HStack(spacing: 12) {
                        Image(systemName: "phone.fill")
                            .foregroundStyle(themeManager.colors.accent)
                        Text(phone)
                            .foregroundStyle(.primary)
                    }
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                }
                .glassEffect(.regular.interactive(), in: .capsule)
            }

            if let email = business.email,
               let url = URL(string: "mailto:\(email)") {
                Link(destination: url) {
                    HStack(spacing: 12) {
                        Image(systemName: "envelope.fill")
                            .foregroundStyle(themeManager.colors.accent)
                        Text(email)
                            .foregroundStyle(.primary)
                    }
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                }
                .glassEffect(.regular.interactive(), in: .capsule)
            }

            if let website = business.website, let url = URL(string: website) {
                Link(destination: url) {
                    HStack(spacing: 12) {
                        Image(systemName: "globe")
                            .foregroundStyle(themeManager.colors.accent)
                        Text(website)
                            .foregroundStyle(.primary)
                            .lineLimit(1)
                    }
                    .padding(.horizontal, 12)
                    .padding(.vertical, 8)
                }
                .glassEffect(.regular.interactive(), in: .capsule)
            }
        }
        .padding(.horizontal)
    }
}
