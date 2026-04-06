import SwiftUI

struct BusinessDetailView: View {
    @Environment(ThemeManager.self) private var themeManager
    @Environment(\.apiClient) private var apiClient
    @State private var vm = BusinessDetailViewModel()
    let business: Business

    var body: some View {
        ScrollView {
            VStack(alignment: .leading, spacing: 20) {
                // TODO: Do we want this Card?
                BusinessCardView(business: business)
                    .padding(.horizontal)
                
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
                    } else {
                        hoursEmptyState
                    }
                }

                // Menus
                BusinessMenusSection(menus: vm.businessDetails?.menus ?? [])

                // Contact
                BusinessContactSection(business: business)

                // Location
                BusinessLocationSection(business: business)

                Spacer(minLength: 20)
            }
            .padding(.vertical)
        }
        .background(themeManager.colors.background.ignoresSafeArea())
        .navigationTitle(business.name)
        .navigationBarTitleDisplayMode(.large)
        .task {
            vm.apiClient = apiClient
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
