---
title: "ADR-001: Use Capacitor over Tauri v2 for Mobile"
---


We chose Capacitor + SvelteKit as the mobile framework instead of Tauri v2.

## Status

Accepted

## Context

The original plan specified Tauri v2 for the mobile app. The developer has extensive Tauri v1 desktop experience and knows Rust. However, Tauri v2's mobile support shipped in 2024 and is still maturing. Key concerns:

- WebView inconsistencies between iOS (WKWebView) and Android.
- The mobile plugin ecosystem is thin compared to established alternatives.
- Smaller community means fewer resources for mobile-specific issues.
- The developer would be an early adopter, which means also being a beta tester.

The app is primarily a UI over a Go REST API. It displays lists, detail views, maps, and forms. There is no heavy on-device computation that would benefit from Rust.

We evaluated five options: Tauri v2, React Native/Expo, Capacitor + Svelte, Flutter, and PWA.

## Decision

We chose Capacitor + SvelteKit because:

- The developer keeps Svelte, which they already know and prefer.
- Capacitor is a mature, proven mobile runtime (Ionic ecosystem).
- The plugin ecosystem covers everything in the project plan: SQLite, push notifications, geolocation.
- The Go backend stays unchanged. Most business logic lives there, not on the client.
- The same app can also be deployed as a PWA if needed.
- Google Maps JS SDK (or MapLibre) works naturally in the webview.

React Native was rejected because the developer does not want to learn React. Flutter was rejected because Dart is a new language and the UI paradigm is entirely different from web development. PWA was rejected because iOS Safari support is limited and there is no app store presence.

## Consequences

- Rust is no longer part of the client-side stack. All client code is TypeScript/Svelte.
- The developer needs to learn Capacitor's plugin system and native project structure.
- Custom native plugins (if ever needed) would be written in Swift (iOS) and Kotlin (Android), not Rust.
- If Tauri v2 mobile support matures significantly in the future, we could revisit this decision. See [future-ideas-and-alternatives.md](/future-ideas-and-alternatives/).
