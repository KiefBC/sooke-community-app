//
//  SookeMapView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-04-04.
//

import SwiftUI
import MapKit

struct SookeMapView: View {
    @Environment(ThemeManager.self) private var themeManager
    var sookeLocation: CLLocationCoordinate2D = CLLocationCoordinate2D(latitude: 48.3725, longitude: -123.7255)
    var initialPosition: MapCameraPosition {
        .region(MKCoordinateRegion(
            center: sookeLocation,
            span: MKCoordinateSpan(latitudeDelta: 0.01, longitudeDelta: 0.01)
        ))
    }
    @State private var vm: MapViewModel
    @State private var selectedMarker: Int64?
    @State var navigateToBusiness: Business?
    @State private var navigationPath = NavigationPath()
    let apiClient: APIClient

    init(apiClient: APIClient) {
        self.apiClient = apiClient
        self._vm = State(initialValue: MapViewModel(apiClient: apiClient))
    }
    
    var body: some View {
        @Bindable var vm = vm
        NavigationStack(path: $navigationPath) {
            Map(initialPosition: initialPosition, selection: $selectedMarker) {
                ForEach(vm.filteredBusinesses) { business in
                    let businessLocation = CLLocationCoordinate2D(latitude: business.latitude, longitude: business.longitude)
                    Marker(business.name, coordinate: businessLocation)
                        .tag(business.id)
                        .tint(.red)
                }
                UserAnnotation()
            }
            .preferredColorScheme(ColorScheme.dark)
            .mapControlVisibility(.hidden)
            .task {
                await vm.fetchBusinesses()
                await vm.fetchCategories()
                vm.requestLocationPermission()
            }
            .onChange(of: selectedMarker) {
                guard let id = selectedMarker else { return }
                vm.selectedBusiness = vm.filteredBusinesses.first { $0.id == id }
            }
            .onChange(of: vm.selectedBusiness) {
                withAnimation {
                    selectedMarker = vm.selectedBusiness?.id
                }
            }
            .safeAreaInset(edge: .top) {
                ScrollView(.horizontal, showsIndicators: false) {
                    GlassEffectContainer(spacing: 8) {
                        HStack(spacing: 8) {
                            Button {
                                vm.selectedCategory = nil
                            } label: {
                                Text("All")
                                    .font(.subheadline)
                                    .padding(.horizontal, 12)
                                    .padding(.vertical, 6)
                                    .foregroundColor(vm.selectedCategory == nil ? .white : themeManager.colors.accent)
                            }
                            .buttonStyle(.plain)
                            .glassEffect(
                                vm.selectedCategory == nil ? .regular.tint(themeManager.colors.accent).interactive() : .regular.interactive(),
                                in: .capsule
                            )
                            ForEach(vm.categories) { cat in
                                let isSelected = vm.selectedCategory == cat
                                Button {
                                    vm.selectCategory(isSelected ? nil : cat)
                                } label: {
                                    Text(cat.name)
                                        .font(.subheadline)
                                        .padding(.horizontal, 12)
                                        .padding(.vertical, 6)
                                        .foregroundColor(isSelected ? .white : themeManager.colors.accent)
                                }
                                .buttonStyle(.plain)
                                .glassEffect(
                                    isSelected ? .regular.tint(themeManager.colors.accent).interactive() : .regular.interactive(),
                                    in: .capsule
                                )
                            }
                        }
                        .padding(.horizontal)
                    }
                }
                .padding(.vertical, 5)
                .background(themeManager.colors.background)
            }
            .sheet(item: $vm.selectedBusiness) { business in
                VStack {
                    BusinessPreviewCard(business: business) {
                        navigationPath.append(business)
                    }
                }
                .presentationDetents([.height(120), .medium])
            }
            .navigationDestination(for: Business.self) { business in
                BusinessDetailView(business: business, apiClient: apiClient)
            }
        }
    }
}

#Preview {
    SookeMapView(apiClient: APIClient(baseURL: APIConfig.baseURL))
}
