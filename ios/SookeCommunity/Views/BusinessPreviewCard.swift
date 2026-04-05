//
//  BusinessPreviewCard.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-04.
//

import SwiftUI

struct BusinessPreviewCard: View {
    @Environment(\.dismiss) private var dismiss
    @Environment(ThemeManager.self) private var themeManager
    let business: Business
    let onViewDetails: () -> Void
    var currentHoursStatus: HoursStatus { business.hoursStatus() }
    
    var body: some View {
        Text(business.name)
        Text(business.description ?? "No description available")
        Text(currentHoursStatus.displayText)
            .foregroundStyle(hoursColor)
        Button("View Details") {
            dismiss()
            onViewDetails()
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

//#Preview {
//    BusinessPreviewCard()
//}
