---
title: "Planning Discussion Log"
---


This document is a condensed record of the planning conversations that shaped the Sooke Community App project. It is organized by topic, not chronologically. Each entry captures the question, the discussion, and the decision reached.

This document is append-only. Do not edit past entries. New discussions are added at the end.

---

## Entry 1: Mobile Framework Selection

**Date:** 2026-03-08

**Question:** Should we use Tauri v2 for mobile, or consider alternatives?

**Discussion:**
The original plan specified Tauri v2. The developer has extensive Tauri v1 desktop experience and knows Rust well. However, Tauri v2's mobile support shipped in 2024 and is still maturing. Concerns included:

- WebView inconsistencies between iOS and Android.
- Plugin ecosystem gaps for mobile.
- Debugging on-device is harder than web development.
- Small community means fewer resources for mobile-specific issues.

We evaluated five alternatives:

1. **Tauri v2** -- Rust-based, developer knows it, but immature mobile support.
2. **React Native / Expo** -- Battle-tested, but requires learning React. Developer is not a fan of React.
3. **Capacitor + Svelte** -- Keeps Svelte, mature mobile runtime, large plugin ecosystem.
4. **Flutter** -- Cross-platform, but requires learning Dart. Steep learning curve.
5. **PWA** -- Simplest option, but limited iOS support and no app store presence.

The developer's priorities: keep Svelte, leverage Go on the backend, avoid React, and use a proven mobile runtime.

**Decision:** Capacitor + SvelteKit. See [ADR-001](/decisions/001-use-capacitor-over-tauri/).

**Reasoning:** Capacitor lets the developer keep Svelte (no new UI framework), has a proven mobile runtime, and the app is primarily a UI over a Go API -- there is no heavy on-device computation that would benefit from Rust.

---

## Entry 2: Go HTTP Framework

**Date:** 2026-03-08

**Question:** Fiber or Chi for the Go backend?

**Discussion:**
The original plan listed "Fiber (or Chi)" without a firm decision. We compared the two:

- Fiber uses `fasthttp`, a custom HTTP engine. It has an Express-like API but does not support standard `net/http` middleware.
- Chi is built on Go's standard `net/http`. Any Go HTTP middleware works with Chi out of the box.
- At the scale of this app, Fiber's performance advantage is irrelevant.
- Using Fiber would require adapting standard Go libraries (Clerk JWT validation, CORS, rate limiting) to work with `fasthttp`.

The developer valued not being tied to an ecosystem and appreciated that "boring and predictable" is a strength, not a weakness.

**Decision:** Chi. See [ADR-003](/decisions/003-use-chi-over-fiber/).

---

## Entry 3: Maps Provider

**Date:** 2026-03-08

**Question:** Google Maps, Leaflet, MapLibre, or something else?

**Discussion:**
The original plan specified Google Maps JS SDK. We identified concerns:

- API key restrictions are complex in a webview (referrer restrictions do not apply to local webviews).
- Google Maps is overkill for pin-dropping on a small-town map.
- The free tier ($200/month) is generous, but there are free alternatives.

We evaluated three options:

1. **Google Maps JS SDK** -- polished but overkill, API key complexity.
2. **Leaflet + OpenStreetMap** -- free, no API key, but less polished than vector tile maps.
3. **MapLibre GL JS** -- open-source, vector tiles, modern look. Needs a tile provider.

The developer looked at MapLibre and liked the look of it.

For tile providers, we discussed MapTiler (free tier: 100k tile loads/month) vs self-hosting tiles. Self-hosting involves downloading OpenStreetMap data, generating tiles with tools like `tileserver-gl`, and hosting them on your own infrastructure. This is unnecessary at the app's scale but documented for future reference.

**Decision:** MapLibre + MapTiler free tier. See [ADR-004](/decisions/004-use-maplibre-over-google-maps/).

---

## Entry 4: Auth Provider and Role Model

**Date:** 2026-03-08

**Question:** Which auth provider? How do roles work? How do we verify business owners?

**Discussion:**
We evaluated Clerk, Firebase Auth, Supabase Auth, Auth0, and roll-your-own. The developer ruled out roll-your-own immediately due to security concerns.

Firebase Auth is completely free with no user cap. Clerk has a 10,000 MAU free tier but better developer experience and prebuilt UI components. Given that most users will not create accounts (browsing is anonymous), the 10k cap is not a concern.

We defined four roles:

1. **Anonymous visitor** -- no account, browse only.
2. **General user** -- social login, can submit events and receive notifications.
3. **Business owner** -- social login + manually promoted by Super Admin.
4. **Super Admin** -- social login + hardcoded role.

The key security concern: how to prevent someone from impersonating a business owner. The solution separates two concerns:

- **Identity verification:** Social login confirms who someone is (Google/Apple verifies their identity).
- **Business ownership verification:** The Super Admin manually verifies and promotes the user after out-of-band confirmation (email from the business's public address, in-person verification, etc.).

Businesses can only be created by the Super Admin. No self-registration. This eliminates fake business listings entirely.

Accounts are optional. The app is fully usable without logging in. Users only need accounts for actions that require identity (submitting events, notifications).

**Decision:** Clerk with social login. Manual role promotion by Super Admin. See [ADR-002](/decisions/002-use-clerk-for-auth/).

---

## Entry 5: Event System Design

**Date:** 2026-03-08

**Question:** How do events work? Who can create them? How are they approved? How are event locations associated with businesses?

**Discussion:**
Community members should be able to submit events, but all events require approval to prevent spam and fake content.

The proximity problem: if a user creates an event near a business, how do we determine which business "owns" the area? Using geographic proximity is ambiguous when multiple businesses are near each other.

The solution: do not use proximity at all. The event submitter explicitly chooses the location type:

- **"At a business"** -- select from a dropdown. The business owner gets notified and can approve or reject. If the submitter is the business owner, the event coordinates auto-fill from the business's stored lat/lng.
- **"Public location"** -- drop a pin on the map or type an address. Only the Super Admin can approve.

Event status values: `draft`, `pending_review`, `approved`, `rejected`.

Business owners act as scoped admins for their venue only. They can approve or reject events at their business but have no control over other businesses or public events.

Spam prevention measures:
- Businesses are Super Admin-created only.
- Business owners are manually verified.
- Event submissions require a logged-in account.
- Rate limiting (max 5 pending events per user).
- Only approved events appear in the public feed.

**Decision:** Explicit location selection (not proximity). Business owners approve events at their venue. Super Admin approves public location events.

---

## Entry 6: Tag System

**Date:** 2026-03-08

**Question:** Should users be able to create their own tags, or should tags be curated?

**Discussion:**
Free-form tags offer flexibility but create moderation challenges. The Go library `goaway` provides profanity detection with leet-speak support, but profanity filters cannot distinguish context (e.g., "gay" as identity vs. slur). False positives and false negatives are both problematic for a community app representing a real town.

Curated tags eliminate the moderation problem entirely. The Super Admin defines a fixed set of categories for businesses and event types. If users need a new tag, they request it and the Super Admin reviews.

**Decision:** Curated tags for MVP. Document `goaway` as a future option if free-form tags are needed later. See [ADR-006](/decisions/006-curated-tags-over-freeform/).

---

## Entry 7: Image Storage

**Date:** 2026-03-08

**Question:** S3 or another option for storing images?

**Discussion:**
The developer initially mentioned S3 or "some other AWS data holder." We identified Cloudflare R2 as a better fit:

- S3-compatible API (same SDK, same code).
- No egress fees (S3 charges for every image download).
- Free tier: 10GB storage, 10 million reads/month.
- The developer already uses Cloudflare for self-hosting.

**Decision:** Cloudflare R2. See [ADR-005](/decisions/005-use-cloudflare-r2-for-images/).

---

## Entry 8: Database Design Decisions

**Date:** 2026-03-08

**Question:** What should the schema include? How should IDs work? What about business hours and coordinates?

**Discussion:**
Key schema decisions:

- **Slug-based IDs:** Every public-facing entity (businesses, events) has both a numeric primary key and a unique slug (e.g., `joes-coffee-shop`). Slugs are used in API responses and prepare for future deep linking. Adding slugs now is cheap; retrofitting them later is painful.
- **Business hours:** Use a `business_hours` table rather than a JSON blob. Business hours are complex (different hours per day, seasonal changes, holiday closures).
- **Coordinates:** Store lat/lng on businesses and events for map pins. Businesses set their coordinates once. Events at a business reference the business's coordinates. Events at public locations store their own coordinates.
- **Event types and business categories:** Separate tables for the curated tag system.
- **Device tokens:** Table for linking users to push notification tokens (needed in a future phase).

**Decision:** Schema includes businesses, menus, menu_items, events, event_types, business_categories, users, roles, business_hours, and device_tokens. Slugs on all public entities.

---

## Entry 9: API Design

**Date:** 2026-03-08

**Question:** API versioning? Search and filtering?

**Discussion:**
API versioning: prefix all routes with `/api/v1/`. When breaking changes are needed, create `/v2` routes alongside `/v1`. This costs nothing to add upfront.

Search and filtering: required from day one. The number of businesses in Sooke is not precisely known, but the target is community businesses (not chain stores like McDonald's). Postgres `ILIKE` or full-text search (`tsvector`) is sufficient at this scale.

**Decision:** `/api/v1/` prefix. Search and filtering as a priority feature.

---

## Entry 10: Documentation Tooling

**Date:** 2026-03-08

**Question:** How should we host and organize project documentation?

**Discussion:**
The developer values clear documentation from the start and wants to avoid documentation debt. Documentation should be accessible to future team members and business owners.

We evaluated:

1. **Plain markdown in `docs/`** -- zero tooling, works everywhere, but no search, no sidebar.
2. **mdBook** -- Rust-based, simple, self-hostable.
3. **VitePress** -- modern, polished, self-hostable.
4. **Mintlify** -- polished, but not self-hostable. Vendor lock-in.
5. **Starlight (Astro)** -- modern, documentation-focused, self-hostable on Cloudflare Pages.

Mintlify was ruled out because it is a hosted SaaS platform and cannot be self-hosted. The developer prefers owning infrastructure.

The developer liked Starlight's look and its purpose-built documentation features.

**Decision:** Starlight, hosted on Cloudflare Pages. See [ADR-007](/decisions/007-documentation-tooling/).

---

## Entry 11: Development Approach

**Date:** 2026-03-08

**Question:** How should we structure development phases?

**Discussion:**
The developer prefers small iterative phases. "Small changes over time lead to big changes." Each phase is a milestone containing sub-tasks. Every sub-task must include tests and be verified before it is considered complete.

We defined 18 phases, from project scaffolding through app store submission. Each phase focuses on a specific area (e.g., "Database Schema and Migrations", "Clerk Auth Integration") and contains detailed sub-task checklists.

Testing at every layer is important: unit tests (Go and Svelte), integration tests (API and DB), API contract tests, and E2E tests (Playwright). Go tests use the table-driven pattern.

**Decision:** 18 iterative phases with sub-task checklists. Tests required for every sub-task.

---

## Entry 12: Admin Dashboard

**Date:** 2026-03-08

**Question:** Separate web app or embedded in the mobile app?

**Discussion:**
A simple SvelteKit web app, separate from the mobile app. Same Go API, Super Admin auth. CRUD operations for businesses, events, users, and tags. No telemetry or analytics at MVP.

Hosted on Cloudflare Pages since the developer already uses Cloudflare.

Advanced features (analytics, moderation tools, audit logs) are deferred to a future phase.

**Decision:** Separate SvelteKit web app on Cloudflare Pages. Simple CRUD for MVP.
