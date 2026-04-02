import Foundation

enum HoursStatus {
    case open(closesAt: String)
    case closingSoon(closesAt: String)
    case opensSoon(opensAt: String)
    case closed
    case unknown

    var isOpen: Bool {
        switch self {
        case .open, .closingSoon:
            return true
        default:
            return false
        }
    }

    var displayText: String {
        switch self {
        case .open(let closesAt):
            return "Open - Closes at \(closesAt.formattedAsTime)"
        case .closingSoon(let closesAt):
            return "Closing soon - \(closesAt.formattedAsTime)"
        case .opensSoon(let opensAt):
            return "Opens soon - \(opensAt.formattedAsTime)"
        case .closed:
            return "Closed"
        case .unknown:
            return "Hours not available"
        }
    }
}
