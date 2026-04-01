//
//  HoursStatus.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import Foundation

enum HoursStatus {
    case open(closesAt: String)
    case opensSoon(opensAt: String)
    case closed
    case unknown

    var isOpen: Bool {
        if case .open = self {
            return true
        }
        return false
    }

    var displayText: String {
        switch self {
        case .open(let closesAt):
            return "Open - Closes at \(closesAt.formattedAsTime)"
        case .opensSoon(let opensAt):
            return "Opens at \(opensAt.formattedAsTime)"
        case .closed:
            return "Closed"
        case .unknown:
            return "Hours not available"
        }
    }
}
