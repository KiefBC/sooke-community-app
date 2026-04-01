//
//  BusinessListView.swift
//  SookeCommunity
//
//  Created by Kiefer Hay on 2026-03-27.
//

import SwiftUI

struct BusinessListView: View {
    @Environment(ThemeManager.self) private var themeManager
    let apiClient: APIClient
    @State private var vm: BusinessListViewModel

    init(apiClient: APIClient) {
        self.apiClient = apiClient
        self._vm = State(initialValue: BusinessListViewModel(apiClient: apiClient))
    }

    var body: some View {
        @Bindable var vm = vm
        NavigationStack {
            VStack {
                ScrollView(.horizontal, showsIndicators: false) {
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
                                    .background(
                                        RoundedRectangle(cornerRadius: 16)
                                            .fill(isSelected ? themeManager.colors.accent : Color.clear)
                                    )
                                    .overlay(
                                        RoundedRectangle(cornerRadius: 16)
                                            .stroke(themeManager.colors.accent, lineWidth: 1)
                                    )
                            }
                            .buttonStyle(.plain)
                        }
                    }
                    .padding(.horizontal)
                }
                .padding(.vertical, 5)
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
                    if !vm.searchText.isEmpty {
                        do {
                            try await Task.sleep(for: .milliseconds(300))
                        } catch { return }
                    }
                    await vm.fetchBusinesses()
                }
                .scrollContentBackground(.hidden)
                .navigationDestination(for: Business.self) { business in
                    BusinessDetailView(business: business, apiClient: apiClient)
                }
            }
            .task {
                await vm.fetchCategories()
            }
            .onChange(of: vm.selectedCategory) {
                Task { await vm.fetchBusinesses() }
            }
            .background(themeManager.colors.background.ignoresSafeArea())
            .toolbarBackground(themeManager.colors.background, for: .navigationBar)
            .toolbarBackgroundVisibility(.visible, for: .navigationBar)
        }
        .background(themeManager.colors.background.ignoresSafeArea())
    }
}

