//
//  BusinessMapView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-04.
//

import SwiftUI
import MapKit

struct BusinessMapView: View {
    let latitude: Double
    let longitude: Double
    let businessName: String
    
    var businessLocation: CLLocationCoordinate2D {
        CLLocationCoordinate2D(latitude: latitude, longitude: longitude)
    }
    
    var initialPosition: MapCameraPosition {
        .region(MKCoordinateRegion(
            center: businessLocation,
            span: MKCoordinateSpan(latitudeDelta: 0.01, longitudeDelta: 0.01)
        ))
    }
    
    var body: some View {
        Map(initialPosition: initialPosition) {
            Marker(businessName, coordinate: businessLocation)
                .tint(.red)
            
// TODO: Potential Idea For Cards
//            Annotation(businessName, coordinate: businessLocation) {
//                ZStack {
//                    RoundedRectangle(cornerRadius: 5)
//                        .fill(Color.red.opacity(0.8))
//                    Image(systemName: "figure.walk.diamond")
//                        .padding(5)
//                }
//            }

        }
        .mapControlVisibility(.hidden)
        .clipShape(RoundedRectangle(cornerRadius: 12))
        .frame(height: 200)
    }
}

#Preview {
    BusinessMapView(latitude: 48.3716, longitude: -123.7382, businessName: "MEOW")
}
