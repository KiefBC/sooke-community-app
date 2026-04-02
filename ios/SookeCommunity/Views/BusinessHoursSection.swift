import SwiftUI

struct BusinessHoursSection: View {
    @Environment(ThemeManager.self) private var themeManager
    let details: BusinessDetails

    var body: some View {
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
}
