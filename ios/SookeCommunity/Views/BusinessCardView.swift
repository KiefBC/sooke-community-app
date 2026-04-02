//
//  BusinessCardView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import SwiftUI

struct BusinessCardView: View {
    @Environment(ThemeManager.self) private var themeManager
    let business: Business
    let details: BusinessDetails?

    init(business: Business, details: BusinessDetails? = nil) {
        self.business = business
        self.details = details
    }

    var body: some View {
        VStack(alignment: .leading, spacing: 0) {
            // TODO: replace with actual image when available
            ZStack {
                Rectangle()
                    .fill(themeManager.colors.accent.opacity(0.3))

                Image(systemName: "building.2")
                    .font(.system(size: 50))
                    .foregroundStyle(themeManager.colors.accent)
            }
            .frame(height: 180)
            .clipped()

            // Info
            VStack(alignment: .leading, spacing: 8) {
                Text(business.name)
                    .font(.headline)
                    .foregroundStyle(themeManager.colors.accent)
                    .lineLimit(2)

                // Reviews
                HStack(spacing: 4) {
                    ForEach(0..<5) { index in
                        Image(systemName: "star.fill")
                            .font(.caption)
                            .foregroundStyle(themeManager.colors.accent.opacity(0.5))
                    }
                    Text("(No reviews)")
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }

                // Hours
                hoursView

                // Category
                if !business.categoryName.isEmpty {
                    Text(business.categoryName)
                        .font(.caption)
                        .foregroundStyle(.secondary)
                }
            }
            .padding(12)
        }
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .glassEffect(.regular, in: .rect(cornerRadius: 12))
    }

    private var currentHoursStatus: HoursStatus? {
        details?.hoursStatus()
    }

    @ViewBuilder
    private var hoursView: some View {
        HStack(spacing: 4) {
            Image(systemName: "clock")
                .font(.caption)
                .foregroundStyle(hoursColor)

            if let status = currentHoursStatus {
                Text(status.displayText)
                    .font(.caption)
                    .foregroundStyle(hoursColor)
            } else {
                Text("Hours not available")
                    .font(.caption)
                    .foregroundStyle(.secondary)
            }
        }
    }

    private var hoursColor: Color {
        guard let status = currentHoursStatus else {
            return .secondary
        }
        return status.isOpen ? .green : .secondary
    }
}

#Preview {
    let sampleBusiness = Business(
        id: 1,
        name: "The Sooke Harbour House",
        slug: "sooke-harbour-house",
        description: "Fine dining with ocean views",
        categoryName: "Restaurant",
        categorySlug: "restaurant",
        address: "6971 West Coast Rd, Sooke, BC",
        latitude: 48.3754,
        longitude: -123.7322,
        phone: "(250) 642-3421",
        email: "info@sookeharbourhouse.com",
        website: "https://sookeharbourhouse.com"
    )

    BusinessCardView(business: sampleBusiness)
        .frame(width: 300)
        .padding()
        .environment(ThemeManager())
}
