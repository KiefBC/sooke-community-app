//
//  BusinessDetailView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

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
                // Business Card with details
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
                
                // Hours Section
                if let details = vm.businessDetails {
                    if !details.hours.isEmpty {
                        hoursSection(details: details)
                        
                        // TODO: Debug info - remove this later
                        #if DEBUG
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
                        #endif
                    } else {
                        // Empty state for businesses without hours
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
                
                // Menus
                if let details = vm.businessDetails, !details.menus.isEmpty {
                    menusSection(menus: details.menus)
                }

                // Contact Information
                contactSection
                
                // Location
                locationSection
                
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
    
    @ViewBuilder
    private func hoursSection(details: BusinessDetails) -> some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Hours")
                .font(.headline)
                .foregroundStyle(themeManager.colors.accent)
            
            if details.hours.isEmpty {
                Text("No hours available")
                    .font(.subheadline)
                    .foregroundStyle(.secondary)
            } else {
                VStack(alignment: .leading, spacing: 8) {
                    ForEach(0..<7, id: \.self) { dayIndex in
                        HStack {
                            Text(Calendar.current.weekdaySymbols[dayIndex])
                                .font(.subheadline)
                                .foregroundStyle(.primary)
                                .frame(width: 100, alignment: .leading)
                            
                            if let dayHours = details.hours.first(where: { $0.dayOfWeek == dayIndex }) {
                                if dayHours.isClosed {
                                    Text("Closed")
                                        .font(.subheadline)
                                        .foregroundStyle(.secondary)
                                } else {
                                    Text("\(dayHours.openTime.formattedAsTime) - \(dayHours.closeTime.formattedAsTime)")
                                        .font(.subheadline)
                                        .foregroundStyle(.primary)
                                }
                            } else {
                                Text("Not available")
                                    .font(.subheadline)
                                    .foregroundStyle(.secondary)
                            }
                        }
                    }
                }
            }
        }
        .padding(.horizontal)
        .padding(.vertical, 12)
        .background(
            RoundedRectangle(cornerRadius: 12)
                .fill(themeManager.colors.accent.opacity(0.1))
        )
        .padding(.horizontal)
    }
    
    @ViewBuilder
    private func menusSection(menus: [Menu]) -> some View {
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

    @ViewBuilder
    private var contactSection: some View {
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
                }
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
                }
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
                }
            }
        }
        .padding(.horizontal)
    }
    
    @ViewBuilder
    private var locationSection: some View {
        VStack(alignment: .leading, spacing: 12) {
            Text("Location")
                .font(.headline)
                .foregroundStyle(themeManager.colors.accent)
            
            HStack(spacing: 12) {
                Image(systemName: "mappin.and.ellipse")
                    .foregroundStyle(themeManager.colors.accent)
                Text(business.address)
                    .foregroundStyle(.primary)
            }
            
            // Placeholder for map - you can add MapKit here later
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
        website: "https://sookeharbourhouse.com"
    )
    
    NavigationStack {
        BusinessDetailView(business: sampleBusiness, apiClient: APIClient(baseURL: APIConfig.baseURL))
    }
    .environment(ThemeManager())
}
