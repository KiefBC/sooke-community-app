import Foundation

// shared threshold for "soon" states
private let soonThreshold: TimeInterval = 3600 // 1 hour

private let timeParser: DateFormatter = {
    let f = DateFormatter()
    f.dateFormat = "HH:mm:ss"
    return f
}()

// Computes HoursStatus from a single BusinessHour and the current time
private func computeStatus(hour: BusinessHour, now: Date) -> HoursStatus {
    if hour.isClosed {
        return .closed
    }

    guard let openTime = timeParser.date(from: hour.openTime),
          let closeTime = timeParser.date(from: hour.closeTime),
          let currentTime = timeParser.date(from: timeParser.string(from: now)) else {
        return .unknown
    }

    if currentTime >= openTime && currentTime <= closeTime {
        let timeUntilClose = closeTime.timeIntervalSince(currentTime)
        if timeUntilClose <= soonThreshold {
            return .closingSoon(closesAt: hour.closeTime)
        }
        return .open(closesAt: hour.closeTime)
    } else if currentTime < openTime {
        let timeUntilOpen = openTime.timeIntervalSince(currentTime)
        if timeUntilOpen <= soonThreshold {
            return .opensSoon(opensAt: hour.openTime)
        }
        return .closed
    } else {
        return .closed
    }
}

extension BusinessDetails {
    func hoursStatus(for date: Date = Date()) -> HoursStatus {
        let calendar = Calendar.current
        let dayOfWeek = calendar.component(.weekday, from: date) - 1

        guard let todayHours = hours.first(where: { $0.dayOfWeek == dayOfWeek }) else {
            return .unknown
        }

        return computeStatus(hour: todayHours, now: date)
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

extension Business {
    func hoursStatus(for date: Date = Date()) -> HoursStatus {
        guard let hour = todayHours else {
            return .unknown
        }
        return computeStatus(hour: hour, now: date)
    }
}
