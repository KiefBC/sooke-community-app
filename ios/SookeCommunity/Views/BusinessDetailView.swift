import SwiftUI

struct BusinessDetailView: View {
    @Environment(ThemeManager.self) private var themeManager
    let business: Business
    @State private var vm: BusinessDetailViewModel

    init(business: Business, apiClient: APIClient) {
        self.business = business
        self._vm = State(initialValue: BusinessDetailViewModel(apiClient: apiClient))
    }

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                // Card
                if let details = vm.businessDetails {
                    BusinessCardView(business: business, details: details)
                        .padding(.horizontal)
                } else {
                    BusinessCardView(business: business)
                        .padding(.horizontal)
                }

                // Description
                if let description = business.description, !description.isEmpty {
                    VStack(alignment: .leading, spacing: 8) {
                        Text("About")
                            .font(.headline)
                            .foregroundStyle(themeManager.colors.accent)

                        Text(description)
                            .font(.body)
                            .foregroundStyle(.primary)
                    }
                    .padding(.horizontal)
                }

                // Hours
                if let details = vm.businessDetails {
                    if !details.hours.isEmpty {
                        BusinessHoursSection(details: details)

                        #if DEBUG
                        hoursDebugView(details: details)
                        #endif
                    } else {
                        hoursEmptyState
                    }
                }

                // Menus
                if let details = vm.businessDetails, !details.menus.isEmpty {
                    BusinessMenusSection(menus: details.menus)
                }

                // Contact
                BusinessContactSection(business: business)

                // Location
                BusinessLocationSection(address: business.address)

                Spacer(minLength: 20)
            }
            .padding(.vertical)
        }
        .background(themeManager.colors.background.ignoresSafeArea())
        .navigationTitle(business.name)
        .navigationBarTitleDisplayMode(.large)
        .task {
            await vm.fetchBusinessDetails(slug: business.slug)
        }
        .overlay {
            if vm.isLoading && vm.businessDetails == nil {
                ProgressView()
                    .scaleEffect(1.5)
            }
        }
        .alert("Error Loading Details", isPresented: .constant(vm.error != nil)) {
            Button("OK") {
                // Dismiss
            }
        } message: {
            if let error = vm.error {
                Text(error.localizedDescription)
            }
        }
    }

    // TODO: Debug info - remove this later
    @ViewBuilder
    private func hoursDebugView(details: BusinessDetails) -> some View {
        VStack(alignment: .leading, spacing: 8) {
            Text("Debug Info")
                .font(.caption)
                .foregroundStyle(.secondary)
            Text(verbatim: "Hours count: \(details.hours.count)")
                .font(.caption)
                .foregroundStyle(.secondary)
            ForEach(details.hours.indices, id: \.self) { index in
                let hour = details.hours[index]
                Text(verbatim: "Day \(hour.dayOfWeek): \(hour.openTime) - \(hour.closeTime) (Closed: \(hour.isClosed))")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
        }
        .padding()
        .background(Color.yellow.opacity(0.2))
        .cornerRadius(8)
        .padding(.horizontal)
    }

    private var hoursEmptyState: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Hours")
                .font(.headline)
                .foregroundStyle(themeManager.colors.accent)

            HStack(spacing: 8) {
                Image(systemName: "clock.badge.questionmark")
                    .foregroundStyle(.secondary)
                Text("Hours not yet available for this business")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            }
        }
        .padding()
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(themeManager.colors.accent.opacity(0.1))
        )
        .padding(.horizontal)
    }
}

#Preview {
    let sampleBusiness = Business(
        id: 1,
        name: "The Sooke Harbour House",
        slug: "sooke-harbour-house",
        description: "Fine dining with ocean views and locally sourced ingredients. Experience the best of Vancouver Island cuisine in a beautiful waterfront setting.",
        categoryName: "Restaurant",
        categorySlug: "restaurant",
        address: "6971 West Coast Rd, Sooke, BC V9Z 0V1",
        latitude: 48.3754,
        longitude: -123.7322,
        phone: "(250) 642-3421",
        email: "info@sookeharbourhouse.com",
        website: "https://sookeharbourhouse.com",
        todayHours: nil
    )

    NavigationStack {
        BusinessDetailView(business: sampleBusiness, apiClient: APIClient(baseURL: APIConfig.baseURL))
    }
    .environment(ThemeManager())
}
