---
title: "ADR-013: Native iOS/Android over Capacitor"
---

We are replacing Capacitor + SvelteKit with native platform apps: SwiftUI for iOS and Kotlin + Jetpack Compose for Android.

## Status

Accepted (supersedes ADR-001)

## Context

After building the initial frontend with SvelteKit + Capacitor, the development experience felt like building a website in a native wrapper rather than a native app. Key issues:

- Zero native Capacitor plugins were used. The app was purely a web app in a WebView.
- No custom native code existed in the iOS or Android directories.
- The UI lacked native scroll physics, transitions, and platform-specific design patterns.
- Native device capabilities (push notifications, camera, offline storage) require plugins that add abstraction over what's already available natively.

The Go backend is entirely API-driven and does not care what consumes it. The frontend framework is the only thing changing.

## Decision

We chose native platform development because:

- SwiftUI and Jetpack Compose provide real platform UI with native widgets, animations, and design language.
- Native SDKs for maps (MapLibre), auth (Firebase), and push notifications (FCM) are more capable than their JavaScript equivalents.
- The development experience uses platform-native tooling (Xcode, Android Studio) with proper debugging, profiling, and previewing.
- The target audience is mostly iPhone users. A native iOS app provides the best experience for them.
- iOS is built first (SwiftUI). Android (Kotlin + Compose) is ported later using the iOS app as a blueprint.

The Go backend + PostgreSQL database stays unchanged. All business logic remains in Go.

## Consequences

- The SvelteKit + Capacitor frontend is archived on branch `archive/sveltekit-capacitor`.
- TypeScript/Svelte is no longer part of the client-side stack for the public app. The iOS app is Swift, the future Android app is Kotlin.
- The admin dashboard remains a web application (separate concern, desktop audience).
- Two native codebases will eventually need maintenance. The Go API is the shared layer.
- The visual design (5 Sooke themes, card layouts, navigation patterns) carries forward as reference for the SwiftUI implementation.
