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

    private var currentHoursStatus: HoursStatus {
        business.hoursStatus()
    }

    @ViewBuilder
    private var hoursView: some View {
        HStack(spacing: 4) {
            Image(systemName: "clock")
                .font(.caption)
                .foregroundStyle(hoursColor)

            Text(currentHoursStatus.displayText)
                .font(.caption)
                .foregroundStyle(hoursColor)
        }
    }

    private var hoursColor: Color {
        switch currentHoursStatus {
        case .open:
            return themeManager.colors.statusOpen
        case .closingSoon:
            return themeManager.colors.statusSoon
        case .opensSoon:
            return themeManager.colors.statusSoon
        case .closed:
            return themeManager.colors.statusClosed
        case .unknown:
            return .secondary
        }
    }
}

