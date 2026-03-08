---
title: "ADR-002: Use Clerk for Authentication"
---


We chose Clerk as the authentication provider for the Sooke Community App.

## Status

Accepted

## Context

The app needs authentication for three purposes:

1. Identifying users who submit events (social login).
2. Scoping business owner permissions to their specific business.
3. Identifying the Super Admin for administrative actions.

Most users will never create an account. Browsing is anonymous. Accounts are only needed for event submission, notifications, and business management.

We evaluated five options: Clerk, Firebase Auth, Supabase Auth, Auth0, and roll-your-own.

Roll-your-own was immediately rejected due to operational security concerns. Implementing auth from scratch is a significant risk for a solo developer.

## Decision

We chose Clerk because:

- Better developer experience and prebuilt UI components compared to Firebase Auth.
- Social login (Google, Apple, Facebook) is straightforward to configure.
- Free tier supports 10,000 monthly active users, which is more than sufficient for a small-town community app where most users browse anonymously.
- Clerk JWTs can be validated in Go using standard JWT libraries. No vendor-specific SDK is needed on the backend.
- Clerk metadata can store role information (business owner, Super Admin) alongside the user record.

Firebase Auth was the runner-up. It is completely free with no user cap and has a well-supported Go Admin SDK. It remains documented as an alternative if Clerk's pricing changes or the free tier becomes limiting.

## Consequences

- The Go backend validates Clerk JWTs in Chi middleware on every protected route.
- Role promotion (General User to Business Owner) is done by the Super Admin via the admin dashboard, updating both Clerk metadata and the backend database.
- If Clerk's free tier becomes limiting or pricing changes unfavorably, Firebase Auth is a documented fallback. See [future-ideas-and-alternatives.md](/future-ideas-and-alternatives/).
- The frontend uses Clerk's prebuilt UI components for login/signup flows.
