---
title: "SwiftUI Data Flow Guide"
---


# SwiftUI Data Flow Guide

A practical reference for how data moves from the Postgres database to SwiftUI views in this project, and when to use each Swift property wrapper.


## The Full Data Pipeline

Data flows through five layers. Each layer has one job.

```
Postgres (stores data)
   |
   v
Go API (queries DB, returns JSON)
   |  GET /api/v1/businesses
   v
APIClient (sends HTTP requests, decodes JSON)
   |  async/await
   v
ViewModel (holds state, calls APIClient)
   |  @Observable
   v
SwiftUI View (renders UI from ViewModel)
```

Here is how a business listing reaches the screen:

1. Postgres holds the row in the `businesses` table.
2. The Go handler calls `repository.ListBusinesses()`, which runs SQL and returns Go structs.
3. The handler encodes the structs as JSON and writes the HTTP response.
4. `APIClient.get()` sends the request, receives the JSON, and decodes it into Swift `Codable` structs.
5. `BusinessListViewModel.fetchBusinesses()` calls the APIClient and stores the result in `items`.
6. `BusinessListView` reads `vm.items` and renders a `List`.

When the ViewModel's `items` array changes, SwiftUI automatically re-renders the view because the ViewModel is `@Observable`.


## Models: Matching Go JSON to Swift Structs

The Go API returns snake_case JSON. Swift structs use camelCase. `CodingKeys` bridges them.

Go struct:

```go
type Business struct {
    ID           int64   `json:"id"`
    CategoryName string  `json:"category_name"`
    Phone        *string `json:"phone"`
}
```

Swift struct:

```swift
struct Business: Codable, Identifiable, Sendable {
    let id: Int64
    let categoryName: String
    let phone: String?

    enum CodingKeys: String, CodingKey {
        case id
        case categoryName = "category_name"
        case phone
    }
}
```

Key rules:
- Go `*string` (pointer) maps to Swift `String?` (optional). A Go `nil` pointer becomes JSON `null`, which Swift decodes as `nil`.
- Go `string` (non-pointer) maps to Swift `String` (non-optional).
- Go `[]Type` must be initialized to an empty slice (`[]Type{}`), not left as `nil`. Go marshals `nil` slices as `null`, but Swift decodes `[Type]` as a non-optional array and rejects `null`. See `repository/business.go` where we set `bd.Hours = []BusinessHour{}` before querying.


## APIClient: Making Network Requests

`APIClient` is a generic HTTP client. It sends a request and decodes the response into any `Codable` type.

```swift
// Fetch a paginated list
let response: PaginatedResponse<Business> = try await apiClient.get(
    "/api/v1/businesses",
    queryItems: [URLQueryItem(name: "search", value: "cafe")]
)
let businesses = response.items

// Fetch a single resource
let details: BusinessDetails = try await apiClient.get(
    "/api/v1/businesses/\(slug)"
)
```

The caller specifies the return type (`PaginatedResponse<Business>` or `BusinessDetails`), and `APIClient.get()` decodes the JSON into that type. If decoding fails, it throws `APIError.decodingError`.


## ViewModels: Where State Lives

ViewModels manage the data a view needs. They call the APIClient, store results, and track loading/error state.

```swift
@MainActor
@Observable
final class BusinessListViewModel {
    private let apiClient: APIClient
    var items: [Business] = []
    private(set) var isLoadingBusinesses = false
    var error: Error? = nil

    init(apiClient: APIClient) {
        self.apiClient = apiClient
    }

    func fetchBusinesses() async {
        isLoadingBusinesses = true
        error = nil
        do {
            let response: PaginatedResponse<Business> = try await apiClient.get("/api/v1/businesses")
            items = response.items
        } catch {
            self.error = error
        }
        isLoadingBusinesses = false
    }
}
```

`@MainActor` ensures all property access and mutation happens on the main thread (required for UI updates). `@Observable` tells SwiftUI to watch for property changes and re-render views that read them. `final` prevents subclassing -- ViewModels are concrete, not abstract.


## Property Wrappers: When to Use Each One

This is the most confusing part of SwiftUI. Here is the decision tree.


### `let` -- Immutable data passed from a parent

Use `let` when a view receives data it only reads and never changes.

```swift
struct BusinessCardView: View {
    let business: Business      // parent passes this in, card just displays it
    let details: BusinessDetails?

    var body: some View {
        Text(business.name)
    }
}
```

The card does not own or modify the business. It just renders it.


### `var` (plain) -- Computed or non-reactive properties

Use a plain `var` for computed values or properties that do not need to trigger re-renders.

```swift
var isLoading: Bool {
    isLoadingBusinesses || isLoadingCategories
}
```


### `@State` -- View-owned mutable state

Use `@State` when a view owns a piece of state and is the single source of truth for it. When `@State` changes, the view re-renders.

```swift
struct ContentView: View {
    @State private var selectedTab: AppTab = .home  // this view owns which tab is selected

    var body: some View {
        TabView(selection: $selectedTab) { ... }
    }
}
```

`@State` is also how views own their ViewModels:

```swift
struct BusinessListView: View {
    @State private var vm: BusinessListViewModel

    init(apiClient: APIClient) {
        self._vm = State(initialValue: BusinessListViewModel(apiClient: apiClient))
    }
}
```

The view creates the ViewModel once and owns it for its lifetime. The underscore syntax (`self._vm`) accesses the `State` wrapper directly, which is needed to set the initial value in `init`.


### `@Binding` -- Two-way reference to someone else's state

Use `@Binding` when a child view needs to read AND write a value that a parent owns. `@Binding` does not store anything -- it points back to the parent's `@State`.

```swift
// Parent owns the state
struct ParentView: View {
    @State private var searchText = ""

    var body: some View {
        // $searchText creates a Binding that the child can read and write
        SearchField(text: $searchText)
    }
}

// Child reads and writes through the binding
struct SearchField: View {
    @Binding var text: String

    var body: some View {
        TextField("Search", text: $text)  // typing updates the parent's @State
    }
}
```

The `$` prefix creates a `Binding` from a `@State`. The child writes to `text`, and the parent's `searchText` updates, which re-renders both views.

Common mistake: using `@State` instead of `@Binding` in a child view. That creates a separate copy that does not sync back to the parent.


### `@Bindable` -- Creating bindings from `@Observable` objects

`@Bindable` bridges `@Observable` ViewModels with SwiftUI controls that need `Binding` values (like `TextField` or `searchable`).

```swift
struct BusinessListView: View {
    @State private var vm: BusinessListViewModel

    var body: some View {
        @Bindable var vm = vm  // create a bindable reference inside body

        List { ... }
            .searchable(text: $vm.searchText)  // $vm.searchText is now a Binding
    }
}
```

Without `@Bindable`, you cannot use `$vm.searchText` because `@State` wrapping an `@Observable` object does not automatically provide bindings to the object's properties.


### `@Environment` -- Shared dependencies from parent views

Use `@Environment` to read values injected by a parent view. We use this for the theme and the API client.

```swift
struct BusinessDetailView: View {
    @Environment(ThemeManager.self) private var themeManager  // @Observable object
    @Environment(\.apiClient) private var apiClient            // custom EnvironmentKey

    var body: some View {
        Text("Hello")
            .foregroundStyle(themeManager.colors.accent)
    }
}
```

Two syntaxes:
- `@Environment(Type.self)` -- for `@Observable` classes (like `ThemeManager`). The parent injects with `.environment(themeManager)`.
- `@Environment(\.keyPath)` -- for custom `EnvironmentKey` values (like `apiClient`). The parent injects with `.environment(\.apiClient, apiClient)`.

`@Environment` values flow downward through the view hierarchy. Any descendant can read them without explicit parameter passing.


### `@Observable` -- Making a class reactive

`@Observable` is a macro that makes a class's properties trigger SwiftUI re-renders when they change. Applied to ViewModels.

```swift
@Observable
final class BusinessDetailViewModel {
    var businessDetails: BusinessDetails?  // SwiftUI re-renders when this changes
    var isLoading = false                  // and this
    var error: Error?                      // and this
}
```

Without `@Observable`, changing `isLoading` would not cause the view to update.


## Summary Table

| Wrapper | Who owns it? | Direction | Use when... |
|---|---|---|---|
| `let` | Parent | Read only | Displaying data you did not create |
| `@State` | This view | Read/write | This view is the source of truth |
| `@Binding` | Parent's `@State` | Read/write | Child needs to modify parent's state |
| `@Bindable` | `@Observable` object | Read/write | Need `$` bindings from a ViewModel |
| `@Environment` | Ancestor view | Read only | Shared dependencies (theme, API client) |
| `@Observable` | (class macro) | -- | Making a ViewModel reactive |
| `@MainActor` | (class annotation) | -- | Ensuring main-thread safety for UI state |


## Common Patterns in This Project


### Loading data when a view appears

```swift
.task {
    await vm.fetchCategories()
}
```

`.task` runs once when the view appears and cancels automatically when the view disappears.


### Debounced search

```swift
.task(id: vm.searchText) {
    if !vm.searchText.isEmpty {
        do { try await Task.sleep(for: .milliseconds(300)) } catch { return }
    }
    await vm.fetchBusinesses()
}
```

`.task(id:)` re-runs whenever the `id` value changes. The `Task.sleep` acts as a debounce -- if the user types another character within 300ms, the previous task cancels (throwing `CancellationError`), and a new one starts.


### Passing the API client through views

The API client is created once in `SookeCommunityApp` and injected via environment. Views that create ViewModels read it and pass it down.

```swift
// App root injects it
ContentView()
    .environment(\.apiClient, apiClient)

// ContentView reads it and passes to children
struct ContentView: View {
    @Environment(\.apiClient) private var apiClient

    var body: some View {
        BusinessListView(apiClient: apiClient)
    }
}

// BusinessListView passes it to detail view via navigation
.navigationDestination(for: Business.self) { business in
    BusinessDetailView(business: business, apiClient: apiClient)
}
```


### Showing loading and error states

```swift
.overlay {
    if vm.isLoading && vm.businessDetails == nil {
        ProgressView()
    }
}
.alert("Error", isPresented: .constant(vm.error != nil)) {
    Button("OK") { }
} message: {
    if let error = vm.error {
        Text(error.localizedDescription)
    }
}
```


## References

- `ios/SookeCommunity/Services/APIClient.swift` -- HTTP client
- `ios/SookeCommunity/Services/APIClientEnvironment.swift` -- Environment injection
- `ios/SookeCommunity/ViewModels/BusinessListViewModel.swift` -- ViewModel pattern
- `ios/SookeCommunity/Views/BusinessListView.swift` -- View with search, list, navigation
- `ios/SookeCommunity/Views/BusinessDetailView.swift` -- Detail view with async loading
- `server/internal/repository/business.go` -- Go repository (SQL queries)
- `server/internal/handler/business.go` -- Go HTTP handlers
