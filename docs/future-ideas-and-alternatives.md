---
title: "Future Ideas and Alternatives"
---


This document records deferred features, documented alternatives, and ideas to revisit in later phases. It is the reference for "why not X?" and "what about Y?" questions. The project plan references this document where relevant.

---

## Deferred Features

### Offline Caching (On-Device)

**Status:** Deferred. Nice to have, not required for MVP.

**Approach when implemented:**
- Use SQLite on-device via Capacitor's SQLite plugin.
- Cache the last successful API response per screen.
- Display cached data with a "last updated X ago" indicator.
- Refresh data when the device is online.
- Do not build a full sync engine. Keep it simple.

**Considerations:**
- Define what happens when cached data is stale (e.g., a restaurant changed its hours).
- Decide how much data to cache (all businesses? recently viewed only?).
- Choose a cache invalidation strategy (ETags, last-modified timestamps, version numbers).

### Deep Linking

**Status:** Deferred. Slugs are designed into the schema now to make this easier later.

**Approach when implemented:**
- A URL like `sooke.app/business/joes-coffee` opens the app to that business's page (if installed) or falls back to a web page (if not).
- Requires a web domain, platform-specific config (Apple Universal Links, Android App Links), and routing logic in the app.
- Slug-based IDs in the database are already prepared for this.

### Push Notifications

**Status:** Deferred to a future phase. Significant backend and native layer work.

**Approach when implemented:**
- iOS: Apple Push Notification Service (APNs).
- Android: Firebase Cloud Messaging (FCM).
- Go backend sends push messages when events are approved, reminders before events, etc.
- Capacitor has an official Push Notifications plugin.
- Schema already includes a `device_tokens` table for linking users to their push tokens.

### Advanced Admin Dashboard

**Status:** MVP admin dashboard is a simple CRUD interface. Advanced features deferred.

**Future scope:**
- Telemetry and analytics (user counts, popular businesses, event engagement).
- Moderation tools and audit logs.
- Bulk operations (import businesses from CSV, etc.).
- Dashboard metrics and charts.

### Terms of Service and Privacy Policy

**Status:** Required before app store submission. Deferred until Phase 18.

**Notes:**
- Both Apple App Store and Google Play require these documents.
- Consider using a template service or legal template as a starting point.
- Must cover data collection, social login data usage, and cookie/tracking policies.

### User-Submitted Event Moderation (Advanced)

**Status:** Basic approval flow is in the MVP. Advanced moderation deferred.

**Future scope:**
- Flagging system for inappropriate events.
- Community reporting.
- Automated content screening.
- Moderation logs and audit trail.

---

## Alternative Technologies Considered

### Mobile Framework: Tauri v2

**Decision:** We chose Capacitor. See [ADR-001](/decisions/001-use-capacitor-over-tauri/).

**Why it was considered:**
- The developer has extensive Tauri v1 experience.
- Tauri uses Rust, which the developer knows well.
- Lightweight, small binary sizes, native webview.

**Why it was not chosen:**
- Tauri v2 mobile support shipped in 2024 and is still maturing.
- Smaller community means fewer answers for mobile-specific bugs.
- Plugin ecosystem for mobile is thin compared to Capacitor.
- The app does not require heavy on-device computation that would benefit from Rust.

**Revisit if:** Tauri v2 mobile support matures significantly and the plugin ecosystem grows. Check back in 2026-2027.

### Mobile Framework: React Native / Expo

**Why it was considered:**
- Battle-tested for mobile (Instagram, Discord, Shopify use it).
- Expo simplifies builds, OTA updates, and app store submission.
- Massive ecosystem.

**Why it was not chosen:**
- Requires learning React. The developer knows Svelte and prefers it.
- JavaScript/TypeScript only. No Rust or Go on the client.
- Heavier runtime than Capacitor.

### Mobile Framework: Flutter

**Why it was considered:**
- Truly cross-platform with native-feeling UI.
- Strong Google Maps support.
- Good offline/SQLite story.

**Why it was not chosen:**
- Dart is a new language for the developer.
- Steepest learning curve of all options.
- UI paradigm is completely different from web development.

### Mobile Framework: Progressive Web App (PWA)

**Why it was considered:**
- Simplest option. Just a website with a manifest.
- No app store review process.
- SvelteKit handles this well.

**Why it was not chosen:**
- iOS Safari has limited PWA support.
- No app store presence reduces discoverability for a local community app.
- Feels less "native" to users.

### Auth: Firebase Auth

**Decision:** We chose Clerk. See [ADR-002](/decisions/002-use-clerk-for-auth/).

**Why it was considered:**
- Completely free with no user cap.
- Go has a well-supported Firebase Admin SDK for JWT validation.
- Social login is straightforward.

**Why it was not chosen:**
- Clerk has better developer experience and prebuilt UI components.
- Clerk's free tier (10,000 MAU) is more than sufficient for Sooke.

**Revisit if:** Clerk's free tier becomes limiting or pricing changes unfavorably.

### Auth: Roll Your Own

**Why it was not chosen:**
- Implementing auth from scratch is a significant security risk.
- Operational security concerns outweigh the benefits of full control.
- Social login integration requires significant work to do correctly.

### Go Framework: Fiber

**Decision:** We chose Chi. See [ADR-003](/decisions/003-use-chi-over-fiber/).

**Why it was considered:**
- Express.js-like API, familiar to web developers.
- Faster in synthetic benchmarks.

**Why it was not chosen:**
- Uses `fasthttp` instead of Go's standard `net/http`.
- Standard Go middleware does not work without adaptation.
- Creates ecosystem lock-in.

### Maps: Google Maps JS SDK

**Decision:** We chose MapLibre. See [ADR-004](/decisions/004-use-maplibre-over-google-maps/).

**Why it was considered:**
- Industry standard. Polished map tiles.
- Free tier provides $200 USD/month in credits.

**Why it was not chosen:**
- API key management is complex in a webview (referrer restrictions).
- Overkill for pin-dropping on a small-town map.
- MapLibre with MapTiler's free tier is sufficient and simpler.

**Revisit if:** MapTiler's free tier becomes limiting or map quality is insufficient.

### Maps: Leaflet + OpenStreetMap

**Why it was considered:**
- Open source, free, no API key needed.
- Lightweight and simple.

**Why it was not chosen:**
- MapLibre offers better-looking vector tiles and smoother zooming.
- MapLibre is more modern and has a larger active community.

### Image Storage: AWS S3

**Decision:** We chose Cloudflare R2. See [ADR-005](/decisions/005-use-cloudflare-r2-for-images/).

**Why it was considered:**
- Industry standard for object storage.
- Mature, well-documented.

**Why it was not chosen:**
- Egress fees. S3 charges for every image download.
- R2 is S3-compatible (same SDK, same code) with no egress fees.
- The developer already uses Cloudflare for self-hosting.

### Tag System: Free-Form User Tags with Profanity Filter

**Decision:** We chose curated tags. See [ADR-006](/decisions/006-curated-tags-over-freeform/).

**Why it was considered:**
- More flexible. Users can create tags the admin did not anticipate.
- Go library `goaway` (github.com/TwiN/go-away) provides profanity detection with leet-speak support.

**Why it was not chosen:**
- Profanity filters are blunt instruments. They cannot distinguish context.
- False positives are embarrassing. False negatives are harmful.
- For a community app representing a real town, curated tags are safer.

**Revisit if:** The curated tag list proves too restrictive and users frequently request new tags. Consider free-form tags with `goaway` filtering and admin review queue.

### Self-Hosting Map Tiles

**Why it was considered:**
- Full control over map data. No dependency on a tile provider.

**What it involves:**
- Download OpenStreetMap data for the Sooke region.
- Generate vector tiles using tools like `tileserver-gl`.
- Host the tile server on your own infrastructure.
- Large storage requirements. Significant setup complexity.

**Why it was not chosen:**
- Completely unnecessary at this scale. MapTiler's free tier is more than sufficient.

**Revisit if:** MapTiler changes pricing or the app grows beyond the free tier.

### Caching: Redis

**Why it was not chosen:**
- Not needed at MVP scale. In-memory caching with TTL on the Go backend is sufficient.

**Revisit if:** The app grows to a scale where in-memory caching is insufficient or the API runs on multiple instances (shared cache needed).

### Documentation: Mintlify

**Why it was considered:**
- Polished, modern look. API documentation with interactive "try it" buttons.
- Algolia-powered search.

**Why it was not chosen:**
- Not self-hostable. Documentation lives on Mintlify's infrastructure.
- Vendor lock-in conflicts with the project's preference for owning infrastructure.

### Documentation: mdBook

**Why it was considered:**
- Rust-based, which aligns with the developer's Rust experience.
- Simple, clean output. Zero vendor dependency.

**Why it was not chosen:**
- Less polished than Starlight. Limited customization.
- Rust dependency just for documentation tooling.

### Documentation: VitePress

**Why it was considered:**
- Modern, polished, self-hostable.
- Large community (Vue/Vite ecosystem).

**Why it was not chosen:**
- Starlight is purpose-built for documentation and offers comparable features.
- Starlight's Astro foundation provides better documentation-specific features out of the box.
