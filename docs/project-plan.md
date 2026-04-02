---
title: "Sooke Community App -- Project Plan"
---

This document is the source of truth for the current state of the Sooke Community App project.

---

## Overview

A mobile community app for Sooke, BC, Canada -- a small coastal town with no existing local community app. The app serves Sooke residents and visitors with local business listings, restaurant menus, community events, and a map of local points of interest.

This is a personal project with no monetization goal. The developer is building and maintaining it solo, at least initially.

---

## Goals

- Give Sooke residents a single place to discover local events, restaurants, and businesses.
- Allow verified business owners to manage their own listings and menus.
- Provide a map view of local businesses and event locations.
- Work well on both iOS and Android.
- Support searching and filtering for businesses and events from day one.

---

## Tech Stack

### Mobile App

- **iOS:** SwiftUI (Swift 6.0, iOS 18.0+). Native app built with Xcode and XcodeGen.
- **Android (future):** Kotlin + Jetpack Compose. Ported from the iOS app as a blueprint.
- **Architecture:** MVVM with `@Observable` ViewModels, SwiftUI views, and an API service layer.
- **See:** [ADR-013](/decisions/013-native-ios-android-over-capacitor/) for why we chose native over Capacitor (supersedes ADR-001).

### Backend API

- **Language:** Go
- **Framework:** [Chi](https://github.com/go-chi/chi) -- lightweight, idiomatic, built on `net/http`.
- **Role:** REST API serving events, business listings, menus. Handles auth token validation.
- **See:** [ADR-003](/decisions/003-use-chi-over-fiber/) for why we chose Chi over Fiber.

### Database

- **Local dev:** PostgreSQL on the developer's NAS. No Docker required.
- **Production:** Railway managed PostgreSQL (same schema, environment variable swap only).
- **See:** [ADR-008](/decisions/008-nas-postgres-over-docker/) for why we chose NAS-hosted Postgres over Docker for local development.
- **Schema:** businesses, menus, menu_items, events, event_types, business_categories, users, roles, business_hours, device_tokens.
- **IDs:** Every public-facing entity has both a numeric primary key and a unique slug (e.g., `joes-coffee-shop`). Slugs are used in API responses and prepared for future deep linking.

### Auth

- **Provider:** [Firebase Auth](https://firebase.google.com/docs/auth) -- native SDKs for iOS and Android.
- **Social login:** Google, Apple (required by App Store), Facebook. Available to all users.
- **JWT validation:** Go backend validates Firebase JWTs using firebase-admin-go in Chi middleware on every protected route.
- **See:** [ADR-014](/decisions/014-firebase-auth-over-clerk/) for why we chose Firebase Auth over Clerk (supersedes ADR-002).

### Maps

- **Library:** [MapLibre Native iOS SDK](https://github.com/maplibre/maplibre-native) -- open-source native map renderer.
- **Tile provider:** [MapTiler](https://www.maptiler.com/) free tier (100k tile loads/month).
- **Usage:** Pin businesses and event locations on a map of Sooke.
- **See:** [ADR-004](/decisions/004-use-maplibre-over-google-maps/) for why we chose MapLibre over Google Maps.

### Image Storage

- **Provider:** [Cloudflare R2](https://developers.cloudflare.com/r2/) -- S3-compatible, no egress fees.
- **Usage:** Business logos, photos, and any uploaded media. Store the URL in Postgres, serve from R2.
- **See:** [ADR-005](/decisions/005-use-cloudflare-r2-for-images/) for why we chose R2 over S3.

### Hosting (Production)

- **Platform:** [Railway](https://railway.app/)
- **Services:** Go API container + managed PostgreSQL instance.
- **Config:** All connection strings and secrets via environment variables.

### Admin Dashboard

- **Framework:** SvelteKit (separate web app from the mobile app).
- **Hosting:** Cloudflare Pages.
- **Role:** CRUD operations for businesses, events, users, and tags. Super Admin access only.

### Documentation

- **Framework:** [Starlight](https://starlight.astro.build/) (Astro-based documentation site).
- **Hosting:** Cloudflare Pages.
- **See:** [ADR-007](/decisions/007-documentation-tooling/) for why we chose Starlight.

---

## User Roles

| Role                    | How they authenticate                           | Capabilities                                                                                                             |
| ----------------------- | ----------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| Anonymous visitor       | No account needed                               | Browse businesses, menus, events, map. Read-only.                                                                        |
| General user            | Social login (Google/Apple/Facebook)            | Everything anonymous can do + submit events for review + subscribe to notifications.                                     |
| Business owner          | Social login + manually promoted by Super Admin | Everything general user can do + edit their own business listing, menus, hours. Approve or reject events at their venue. |
| Super Admin (developer) | Social login + hardcoded role                   | Everything. Add/remove businesses, promote users, approve/reject any event, manage all content.                          |

### Business Owner Verification Flow

1. The Super Admin creates the business listing. No user can self-register a business.
2. A real business owner contacts the Super Admin out-of-band (email, in-person, phone).
3. The Super Admin verifies their identity and business ownership.
4. The business owner creates an account using social login (same as any other user).
5. The Super Admin promotes their account to "business owner" scoped to their specific business via the admin dashboard. This links their user ID to the business ID in the database.

Social login provides identity verification. Manual promotion provides business ownership verification. These are two separate concerns handled by two separate mechanisms.

### Accounts Are Optional

The app is fully usable without logging in. Users only need an account to submit events, receive notifications, or manage a business. Most users will never create an account.

---

## Event System

### Event Submission

Any logged-in user can submit an event. Events require approval before they are visible.

### Event Location

When creating an event, the submitter chooses one of two location types:

- **"At a business"** -- select a business from a dropdown. This creates a `business_id` foreign key on the event. The business owner can approve or reject it. If the submitter is the business owner, the event location auto-fills from the business's stored coordinates.
- **"Public location"** -- drop a pin on the map or type an address. No business association. Only the Super Admin can approve.

### Event Approval Flow

1. User submits an event. Status: `pending_review`.
2. If the event is at a business, the business owner is notified (in-app + email) and can approve or reject.
3. If the event is at a public location, the Super Admin reviews and approves or rejects.
4. Approved events are visible to everyone.
5. Rejected events notify the submitter with an optional reason. The submitter can resubmit at a different location if appropriate.

### Event Status Values

`draft`, `pending_review`, `approved`, `rejected`

### Spam Prevention

- Businesses are only created by the Super Admin. No self-registration.
- Business owners are manually verified and promoted.
- Event submissions require a logged-in account (social login adds friction).
- Rate limiting on event submissions (e.g., max 5 pending events per user).
- Only approved events are visible in the public feed.

---

## Tag System

Both businesses and events use a curated tag system. Tags are managed by the Super Admin. Users select from the predefined list and cannot create custom tags.

**Business categories (examples):** Restaurant, Cafe, Bar, Retail, Grocery, Outdoor Recreation, Health and Wellness, Arts and Culture, Professional Services, Accommodation.

**Event types (examples):** Live Music, Market, Workshop, Community Meeting, Sports, Festival, Fundraiser, Kids and Family, Outdoor.

If a user or business owner needs a tag that does not exist, they can request it. The Super Admin reviews and adds it to the master list if appropriate.

See [ADR-006](/decisions/006-curated-tags-over-freeform/) for why we chose curated tags over free-form.

---

## Testing Strategy

We test every layer of the application. The testing pyramid guides our approach: many unit tests, fewer integration tests, fewer E2E tests.

| Layer                   | What it tests                                        | Tools                                               |
| ----------------------- | ---------------------------------------------------- | --------------------------------------------------- |
| Unit tests (Go)         | Individual functions, service logic, validation      | Go standard `testing` package, table-driven tests   |
| Unit tests (Swift)      | ViewModels, models, services, theme                  | Swift Testing framework (`@Suite`, `@Test`)         |
| Integration tests (API) | HTTP handlers with real DB, middleware chains        | Go `testing` + `httptest` + test Postgres container |
| Integration tests (DB)  | Migrations, queries, constraints, foreign keys       | Go `testing` + test Postgres container              |
| API contract tests      | Response shapes, status codes, error formats         | Go `testing` or Hurl                                |

Every milestone issue specifies which test layers are required. No issue is complete until its tests pass.

---

## Core Features (MVP)

### For Anonymous Visitors and General Users

- [ ] Browse local businesses (restaurants, shops, services)
- [ ] Search and filter businesses by name and category
- [ ] View restaurant menus
- [ ] Browse upcoming community events
- [ ] Search and filter events by type
- [ ] Map view with pins for businesses and events
- [ ] Submit community events for review (requires account)

### For Business Owners

- [ ] Edit their own business listing (name, description, hours, contact, location)
- [ ] Add, edit, and remove menu items and prices
- [ ] Approve or reject events submitted at their venue
- [ ] Create events at their own venue (auto-fills location)

### For Super Admin

- [ ] Add new businesses to the directory
- [ ] Create and manage community events
- [ ] Verify and promote business owner accounts
- [ ] Manage curated tag lists (business categories, event types)
- [ ] Admin dashboard (separate SvelteKit web app)

---

## Development Milestones

Each milestone is tracked on GitHub. Each milestone contains one or more consolidated issues. Every issue includes tests in its acceptance criteria -- nothing is complete until tests pass.

Issues link: [github.com/KiefBC/sooke-community-app/issues](https://github.com/KiefBC/sooke-community-app/issues)

### Project Scaffolding (Complete)

- [x] Initialize Capacitor + SvelteKit project
- [x] Verify app builds and runs on iOS simulator
- [x] Verify app builds and runs on Android emulator
- [x] Scaffold Go + Chi API with health check endpoint
- [x] Write test for health check endpoint
- [x] Set up PostgreSQL connection (NAS-hosted -- see ADR-008)
- [x] Verify API connects to Postgres
- [x] Set up `.env` config pattern for local development
- [x] Write integration test for DB connection
- [x] Set up Starlight documentation site and deploy to Cloudflare Pages

### Genesis Protocol -- Database

Set up the Postgres schema, migrations, and seed data that everything else builds on.

- [ ] Database schema, migrations, and seed data ([#107](https://github.com/KiefBC/sooke-community-app/issues/107)) -- schema design, slug fields, lat/lng coordinates, reversible migrations with Goose, constraint tests, and sample Sooke business seed data

### First Contact -- Business Listings

Stand up the business listings API and frontend so users can browse Sooke businesses.

- [ ] Business listings API ([#108](https://github.com/KiefBC/sooke-community-app/issues/108)) -- `GET /api/v1/businesses` (list with search/filter), `GET /api/v1/businesses/:slug` (detail), unit/integration tests, contract tests
- [ ] Business list and detail UI ([#109](https://github.com/KiefBC/sooke-community-app/issues/109)) -- Svelte list and detail components, loading/error states, Vitest component tests
- [ ] Business category filtering ([#110](https://github.com/KiefBC/sooke-community-app/issues/110)) -- `GET /api/v1/categories` endpoint, category filter on business list, filter UI component, tests for both API and UI

### Here Be Dragons -- Maps

Integrate MapLibre GL JS and place business pins on a map of Sooke.

- [ ] MapLibre integration with business pins ([#111](https://github.com/KiefBC/sooke-community-app/issues/111)) -- MapLibre + MapTiler setup, business pins at lat/lng, clickable popups linking to detail, E2E tests, iOS/Android webview verification

### The Daily Bugle -- Events

Build out the events system -- API, frontend, map pins, and the event-business location link.

- [ ] Events API and filtering ([#112](https://github.com/KiefBC/sooke-community-app/issues/112)) -- `GET /api/v1/events` (list with search/filter), `GET /api/v1/events/:slug` (detail), event type filtering, unit/integration/contract tests
- [ ] Event list, detail, and map pins UI ([#113](https://github.com/KiefBC/sooke-community-app/issues/113)) -- Svelte event components, event pins on map (distinct from business pins), Vitest component tests
- [ ] Event-business location association ([#114](https://github.com/KiefBC/sooke-community-app/issues/114)) -- foreign key relationship, event form with "At a business" / "Public location" toggle, auto-fill coordinates from business, association tests

### Who Goes There? -- Auth

Integrate Clerk for social login on the frontend and JWT validation on the backend.

- [ ] Clerk auth integration ([#115](https://github.com/KiefBC/sooke-community-app/issues/115)) -- Clerk account setup, SvelteKit integration, Go Chi JWT middleware, protect write endpoints, middleware tests (valid/expired/missing token), iOS/Android login verification

### The Chuunin Exams -- Roles

Implement the role system that controls who can do what.

- [ ] Role-based access control ([#116](https://github.com/KiefBC/sooke-community-app/issues/116)) -- role model in DB (super_admin, business_owner, general_user), role-checking middleware, scoped business owner permissions, user promotion endpoint, RBAC tests

### The Gatekeepers -- Editing and Approval

Give business owners control of their listings and build the event approval workflow.

- [ ] Business owner editing UI ([#117](https://github.com/KiefBC/sooke-community-app/issues/117)) -- edit form for business listing, menu management UI, business hours editing, E2E tests, permission boundary tests
- [ ] Event submission and approval workflow ([#118](https://github.com/KiefBC/sooke-community-app/issues/118)) -- `POST /api/v1/events`, approve/reject endpoints, rate limiting, submission form, approval queue UI, notifications, full-flow tests

### The Batcave -- Admin Dashboard

Build the Super Admin dashboard as a separate SvelteKit app on Cloudflare Pages.

- [ ] Admin dashboard scaffolding ([#119](https://github.com/KiefBC/sooke-community-app/issues/119)) -- SvelteKit project init, Cloudflare Pages deploy, Super Admin auth check, navigation and layout
- [ ] Admin dashboard CRUD pages ([#120](https://github.com/KiefBC/sooke-community-app/issues/120)) -- business management, event management, user management, tag management, CRUD tests

### Sharingan Activated -- Images

Add image upload support backed by Cloudflare R2.

- [ ] Image upload with Cloudflare R2 ([#121](https://github.com/KiefBC/sooke-community-app/issues/121)) -- R2 bucket setup, upload endpoint (jpg/png/webp, max 5MB), store URLs in Postgres, upload UI with drag-and-drop, upload tests

### The Palantir -- Search

Add full-text search and lock down API versioning.

- [ ] Full-text search and API versioning ([#122](https://github.com/KiefBC/sooke-community-app/issues/122)) -- Postgres full-text search for businesses and events, route audit for `/api/v1/` prefix, versioning docs, search tests

### Evangelion Launch -- Deploy

Containerize the API and deploy everything to Railway.

- [ ] Dockerize and deploy to Railway ([#123](https://github.com/KiefBC/sooke-community-app/issues/123)) -- multi-stage Dockerfile, Railway deploy, migrate Postgres to Railway, update app to production API URL, end-to-end production verification

### Plus Ultra! -- Ship It

Polish, test on real devices, and submit to app stores.

- [ ] Production polish and device testing ([#124](https://github.com/KiefBC/sooke-community-app/issues/124)) -- UI polish and loading states, real iOS device testing, real Android device testing, performance testing on lower-end Android, platform-specific fixes
- [ ] App store submission ([#125](https://github.com/KiefBC/sooke-community-app/issues/125)) -- app store listings and screenshots, Terms of Service and Privacy Policy, Google Play submission, Apple App Store submission

---

## Environment Config

- Use `.env` files for local development.
- Use Railway environment variables for production.
- Never commit secrets to version control.
- Never hardcode connection strings, API keys, or Clerk secrets.

---

## Developer Environment

| Tool                       | Detail                               |
| -------------------------- | ------------------------------------ |
| Primary dev machine        | MacBook Pro M4 Pro 48GB              |
| iOS builds                 | On MBP (Xcode required)              |
| Android testing            | Developer's personal Android phone   |
| Windows PC                 | Available for cross-platform testing |
| Languages                  | Go, TypeScript, Svelte, some Rust    |
| Prior Capacitor experience | None (new for this project)          |
| Prior SvelteKit experience | Familiar with Svelte, some SvelteKit |
| Prior Railway experience   | Basic familiarity                    |

---

## Guard Rails

These are hard constraints for development. Do not deviate from these without explicit discussion and a new ADR.

- Do not scaffold anything not listed in the current milestone.
- Complete the current milestone before starting the next.
- Every issue must include tests before it is considered complete.
- Go tests use the table-driven pattern.
- Always use environment variables for DB connections, API keys, and Clerk secrets. Never hardcode.
- Postgres is the only server-side data store. Do not introduce other stores without discussion.
- Business logic lives in Go, not in the frontend.
- Svelte components should be small and composable.
- MapLibre is used via Typescript in the Capacitor webview. Do not use a native maps plugin.
- Clerk JWT validation happens in Chi middleware on every protected route.
- Businesses are only created by the Super Admin. No self-registration.
- No free-form tags at MVP. Use the curated tag list managed by Super Admin.
- Chi for HTTP routing. Never Fiber.
- API routes are prefixed with `/api/v1/`.
- All documentation follows the rules in [style-guide.md](/style-guide/).

---

## Related Documents

- [Style Guide](/style-guide/) -- documentation formatting and writing rules.
- [Common Commands](/common-commands/) -- how to run, build, and test each part of the project.
- [Future Ideas and Alternatives](/future-ideas-and-alternatives/) -- deferred features and documented alternatives.
- [Planning Discussion](/planning-discussion/) -- condensed log of planning decisions and reasoning.
- [Architecture Decision Records](/decisions/) -- individual records for each major technical decision.
