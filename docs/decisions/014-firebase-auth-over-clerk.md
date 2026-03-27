---
title: "ADR-014: Firebase Auth over Clerk"
---


We are replacing Clerk with Firebase Auth as the authentication provider.

## Status

Accepted (supersedes ADR-002)

## Context

ADR-002 chose Clerk for its developer experience and prebuilt web UI components. With the migration to native iOS/Android (ADR-013), the auth requirements changed:

- Clerk has no native Android SDK. The iOS SDK exists but is newer and less documented.
- Native apps need native auth flows, not web-based UI components.
- Apple requires Sign in with Apple for any iOS app offering social login.
- Push notifications are a planned feature, and Firebase Cloud Messaging (FCM) integrates naturally with Firebase Auth.

## Decision

We chose Firebase Auth because:

- First-party native SDKs for both iOS (Swift) and Android (Kotlin), actively maintained by Google.
- Apple Sign in with Apple is a first-class integration.
- Unlimited free authentication users (no MAU cap).
- Firebase Cloud Messaging comes with the same SDK for push notifications.
- The firebase-admin-go SDK provides straightforward JWT validation for the Go backend.
- The largest ecosystem of iOS-specific documentation, tutorials, and community resources.

The Go backend change is minimal: replace Clerk JWT middleware with Firebase JWT middleware. The `clerk_id` column in the users table becomes `firebase_uid`.

## Consequences

- The Go backend uses firebase-admin-go for JWT validation in Chi middleware.
- The users table `clerk_id` column is renamed to `firebase_uid`.
- Sign in with Apple is mandatory (App Store requirement when offering social login).
- Firebase Console is used for auth user management instead of Clerk Dashboard.
- If Firebase's ecosystem becomes undesirable, the auth layer is abstracted behind an AuthService protocol on iOS and middleware on Go, making provider swaps manageable.
