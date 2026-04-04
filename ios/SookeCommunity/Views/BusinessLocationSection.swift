import SwiftUI
import MapKit

struct BusinessLocationSection: View {
    @Environment(ThemeManager.self) private var themeManager
    let business: Business

    var body: some View {
        VStack(alignment: .leading, spacing: 12) {
                Text("Location")
                    .font(.headline)
                    .foregroundStyle(themeManager.colors.accent)

            HStack(spacing: 12) {
                Image(systemName: "mappin.and.ellipse")
                    .foregroundStyle(themeManager.colors.accent)
                Text(business.address)
                    .foregroundStyle(.primary)
                Spacer()
                Button(action: {
                    let mapItem = MKMapItem(
                        location: CLLocation(
                            latitude: business.latitude,
                            longitude: business.longitude
                        ),
                        address: nil
                    )
                    mapItem.name = business.name
                    
                    let launchOptions: [String: Any] = [
                        MKLaunchOptionsDirectionsModeKey: MKLaunchOptionsDirectionsModeDriving
                    ]
                    
                    mapItem.openInMaps(launchOptions: launchOptions)
                }) {
                    Text("Open in Maps")
                        .font(.caption)
                        .foregroundStyle(themeManager.colors.accent)
                    Image(systemName: "map")
                        .foregroundStyle(themeManager.colors.accent)
                }
                .padding()
                .glassEffect()
            }

            BusinessMapView(latitude: business.latitude, longitude: business.longitude, businessName: business.name)

// TODO: Good sample object? Maybe keep this
//            RoundedRectangle(cornerRadius: 12)
//                .fill(themeManager.colors.accent.opacity(0.2))
//                .frame(height: 200)
//                .overlay {
//                    VStack {
//                        Image(systemName: "map")
//                            .font(.system(size: 40))
//                            .foregroundStyle(themeManager.colors.accent)
//                        Text("Map coming soon")
//                            .font(.caption)
//                            .foregroundStyle(.secondary)
//                    }
//                }
        }
        .padding(.horizontal)
    }
}
