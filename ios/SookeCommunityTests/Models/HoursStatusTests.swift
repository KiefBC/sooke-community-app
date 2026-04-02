import Testing
import Foundation
@testable import SookeCommunity

@Suite("Business Hours Status Tests")
struct HoursStatusTests {

    // helper to build a date at a specific time today
    private func timeToday(_ timeString: String) -> Date {
        let formatter = DateFormatter()
        formatter.dateFormat = "HH:mm:ss"
        guard let parsed = formatter.date(from: timeString) else {
            fatalError("bad time string: \(timeString)")
        }
        return parsed
    }

    private func makeBusiness(hour: BusinessHour?) -> Business {
        Business(
            id: 1,
            name: "Test",
            slug: "test",
            description: nil,
            categoryName: "Cat",
            categorySlug: "cat",
            address: "1 Main St",
            latitude: 48.35,
            longitude: -123.72,
            phone: nil,
            email: nil,
            website: nil,
            todayHours: hour
        )
    }

    @Test func openWhenWithinHours() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: false)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("12:00:00"))

        guard case .open(let closesAt) = status else {
            Issue.record("expected .open, got \(status)")
            return
        }
        #expect(closesAt == "17:00:00")
    }

    @Test func closingSoonWithinOneHour() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: false)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("16:30:00"))

        guard case .closingSoon(let closesAt) = status else {
            Issue.record("expected .closingSoon, got \(status)")
            return
        }
        #expect(closesAt == "17:00:00")
    }

    @Test func opensSoonWithinOneHour() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: false)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("08:30:00"))

        guard case .opensSoon(let opensAt) = status else {
            Issue.record("expected .opensSoon, got \(status)")
            return
        }
        #expect(opensAt == "09:00:00")
    }

    @Test func closedWhenBeforeOpenAndNotSoon() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: false)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("06:00:00"))

        guard case .closed = status else {
            Issue.record("expected .closed, got \(status)")
            return
        }
    }

    @Test func closedWhenAfterClose() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: false)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("18:00:00"))

        guard case .closed = status else {
            Issue.record("expected .closed, got \(status)")
            return
        }
    }

    @Test func closedWhenIsClosedTrue() {
        let hour = BusinessHour(dayOfWeek: 0, openTime: "09:00:00", closeTime: "17:00:00", isClosed: true)
        let business = makeBusiness(hour: hour)
        let status = business.hoursStatus(for: timeToday("12:00:00"))

        guard case .closed = status else {
            Issue.record("expected .closed, got \(status)")
            return
        }
    }

    @Test func unknownWhenNoHours() {
        let business = makeBusiness(hour: nil)
        let status = business.hoursStatus()

        guard case .unknown = status else {
            Issue.record("expected .unknown, got \(status)")
            return
        }
    }

    @Test func closingSoonIsConsideredOpen() {
        let status = HoursStatus.closingSoon(closesAt: "17:00:00")
        #expect(status.isOpen == true)
    }

    @Test func opensSoonIsNotConsideredOpen() {
        let status = HoursStatus.opensSoon(opensAt: "09:00:00")
        #expect(status.isOpen == false)
    }
}
