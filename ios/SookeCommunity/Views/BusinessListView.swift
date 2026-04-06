//
//  BusinessListView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-27.
//

import SwiftUI

struct BusinessListView: View {
    @Environment(ThemeManager.self) private var themeManager
    @Environment(\.apiClient) private var apiClient
    @State private var vm = BusinessListViewModel()

    var body: some View {
        NavigationStack {
            List {
                ForEach(vm.items) { item in
                    NavigationLink(value: item) {
                        BusinessCardView(business: item)
                    }
                    .listRowInsets(EdgeInsets(top: 8, leading: 16, bottom: 8, trailing: 16))
                    .listRowBackground(Color.clear)
                }
            }
            .searchable(text: $vm.searchText, prompt: "Search Businesses")
            .task(id: vm.searchText) {
                vm.apiClient = apiClient
                if !vm.searchText.isEmpty {
                    do {
                        try await Task.sleep(for: .milliseconds(300))
                    } catch { return }
                }
                await vm.fetchBusinesses()
                await vm.fetchCategories()
            }
            .onChange(of: vm.selectedCategory) {
                Task { await vm.fetchBusinesses() }
            }
            .scrollContentBackground(.hidden)
            .navigationDestination(for: Business.self) { business in
                BusinessDetailView(business: business)
            }
            .safeAreaInset(edge: .top) {
                ScrollView(.horizontal, showsIndicators: false) {
                    GlassEffectContainer(spacing: 8) {
                        HStack(spacing: 8) {
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
            .navigationTitle("Businesses")
            .navigationBarTitleDisplayMode(.inline)
            .background(themeManager.colors.background.ignoresSafeArea())
            .toolbarBackground(themeManager.colors.background, for: .navigationBar)
            .toolbarBackgroundVisibility(.visible, for: .navigationBar)
        }
    }
}

