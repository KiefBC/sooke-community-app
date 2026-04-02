//
//  BusinessDetails+Hours.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-29.
//

import Foundation

extension BusinessDetails {
    func hoursStatus(for date: Date = Date()) -> HoursStatus {
        let calendar = Calendar.current
        let dayOfWeek = calendar.component(.weekday, from: date) - 1

        guard let todayHours = hours.first(where: { $0.dayOfWeek == dayOfWeek }) else {
            return .unknown
        }

        if todayHours.isClosed {
            return .closed
        }

        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm:ss"

        guard let openTime = formatter.date(from: todayHours.openTime),
              let closeTime = formatter.date(from: todayHours.closeTime) else {
            return .unknown
        }

        let currentTimeString = formatter.string(from: date)
        guard let currentTime = formatter.date(from: currentTimeString) else {
            return .unknown
        }

        if currentTime >= openTime && currentTime <= closeTime {
            return .open(closesAt: todayHours.closeTime)
        } else if currentTime < openTime {
            return .opensSoon(opensAt: todayHours.openTime)
        } else {
            return .closed
        }
    }

    func todayHoursString(for date: Date = Date()) -> String {
        let calendar = Calendar.current
        let dayOfWeek = calendar.component(.weekday, from: date) - 1

        guard let todayHours = hours.first(where: { $0.dayOfWeek == dayOfWeek }) else {
            return "Hours not available"
        }

        if todayHours.isClosed {
            return "Closed today"
        }

        return "\(todayHours.openTime.formattedAsTime) - \(todayHours.closeTime.formattedAsTime)"
    }
}
