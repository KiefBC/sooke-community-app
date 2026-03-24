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

- **Framework:** [Capacitor](https://capacitorjs.com/) -- wraps a web app in a native shell for iOS and Android.
- **Frontend UI:** [Svelte 5](https://svelte.dev/) with [SvelteKit](https://kit.svelte.dev/) -- file-based routing, lightweight, fast iteration.
- **Language:** TypeScript (UI and frontend logic).
- **Targets:** iOS + Android via Capacitor native projects.
- **See:** [ADR-001](/decisions/001-use-capacitor-over-tauri/) for why we chose Capacitor over Tauri v2.

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

- **Provider:** [Clerk](https://clerk.com/) -- handles all auth UI and session management.
- **Social login:** Google, Facebook, Apple, etc. Available to all users.
- **JWT validation:** Go backend validates Clerk JWTs in Chi middleware on every protected route.
- **See:** [ADR-002](/decisions/002-use-clerk-for-auth/) for why we chose Clerk.

### Maps

- **Library:** [MapLibre GL JS](https://maplibre.org/) -- open-source map renderer embedded in the webview.
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
| Unit tests (Svelte)     | Component rendering, reactive logic, form validation | Vitest + Svelte Testing Library                     |
| Integration tests (API) | HTTP handlers with real DB, middleware chains        | Go `testing` + `httptest` + test Postgres container |
| Integration tests (DB)  | Migrations, queries, constraints, foreign keys       | Go `testing` + test Postgres container              |
| API contract tests      | Response shapes, status codes, error formats         | Go `testing` or Hurl                                |
| E2E tests               | Full user flows through the app                      | Playwright                                          |

Every phase sub-task specifies which test layers are required. No sub-task is complete until its tests pass.

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

## Development Phases

Each phase is a milestone. Each phase contains sub-tasks. Each sub-task must include tests and be verified before it is considered complete.

### Phase 1 -- Project Scaffolding

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

### Phase 2 -- Database Schema and Migrations

- [ ] Design and implement initial Postgres schema (businesses, menus, menu_items, events, users, roles, business_hours, business_categories, event_types)
- [ ] Add slug fields to businesses and events
- [ ] Add lat/lng coordinate fields to businesses and events
- [ ] Write migration scripts (up and down)
- [ ] Write constraint and migration tests
- [ ] Seed database with sample Sooke businesses for development

### Phase 3 -- Business Listings API and UI

- [ ] Implement `GET /api/v1/businesses` (list, with search and filtering)
- [ ] Implement `GET /api/v1/businesses/:slug` (detail)
- [ ] Write unit and integration tests for business endpoints
- [ ] Write API contract tests for response shapes
- [ ] Build business list UI component in Svelte
- [ ] Build business detail UI component in Svelte
- [ ] Write Svelte component tests with Vitest

### Phase 4 -- Business Categories and Tags

- [ ] Implement `GET /api/v1/categories` endpoint
- [ ] Add category filtering to business list endpoint
- [ ] Write tests for category endpoints and filtering
- [ ] Build category filter UI component
- [ ] Write component tests

### Phase 5 -- Map Integration

- [ ] Integrate MapLibre GL JS into the SvelteKit app
- [ ] Display business pins on the map using stored coordinates
- [ ] Build map view with clickable pins linking to business detail
- [ ] Write E2E test for map interaction
- [ ] Test map rendering on iOS and Android

### Phase 6 -- Events API and UI

- [ ] Implement `GET /api/v1/events` (list, with search and filtering)
- [ ] Implement `GET /api/v1/events/:slug` (detail)
- [ ] Add event type filtering
- [ ] Write unit, integration, and contract tests for event endpoints
- [ ] Build event list and detail UI components
- [ ] Display event pins on the map
- [ ] Write component tests

### Phase 7 -- Event-Business Location Association

- [ ] Implement event-business foreign key relationship in API
- [ ] Build event creation form with "At a business" / "Public location" toggle
- [ ] Auto-fill coordinates from business location when applicable
- [ ] Write tests for event-business association logic

### Phase 8 -- Clerk Auth Integration

- [ ] Set up Clerk account and configure social login providers
- [ ] Integrate Clerk into SvelteKit frontend
- [ ] Implement JWT validation middleware in Go Chi
- [ ] Write tests for auth middleware (valid token, expired token, missing token)
- [ ] Protect write endpoints behind auth middleware
- [ ] Verify login flow on iOS and Android

### Phase 9 -- Role System

- [ ] Implement role model in database (Super Admin, Business Owner, General User)
- [ ] Build role-checking middleware in Go
- [ ] Scope business owner permissions to their specific business
- [ ] Write tests for role-based access control
- [ ] Build user promotion flow for Super Admin

### Phase 10 -- Business Owner Editing UI

- [ ] Build edit form for business listing (scoped to owner's business)
- [ ] Build menu management UI (add, edit, remove items)
- [ ] Build business hours editing UI
- [ ] Write E2E tests for owner editing flow
- [ ] Test permission boundaries (owner cannot edit other businesses)

### Phase 11 -- Event Submission and Approval

- [ ] Implement `POST /api/v1/events` (create, sets status to `pending_review`)
- [ ] Implement approval/rejection endpoints for business owners and Super Admin
- [ ] Add rate limiting on event submissions
- [ ] Build event submission form UI
- [ ] Build approval queue UI for business owners
- [ ] Implement notification (in-app + email) for pending approvals
- [ ] Write tests for full approval flow

### Phase 12 -- Admin Dashboard Scaffolding

- [ ] Initialize SvelteKit admin dashboard project
- [ ] Deploy to Cloudflare Pages
- [ ] Implement Super Admin auth check
- [ ] Build navigation and layout

### Phase 13 -- Admin CRUD

- [ ] Build business management pages (list, create, edit, delete)
- [ ] Build event management pages (list, edit, delete, approval queue)
- [ ] Build user management pages (list, promote to business owner, link to business)
- [ ] Build tag management pages (add, edit, remove categories and event types)
- [ ] Write tests for admin operations

### Phase 14 -- Cloudflare R2 Integration

- [ ] Set up R2 bucket and access credentials
- [ ] Implement image upload endpoint in Go API
- [ ] Store image URLs in Postgres
- [ ] Build image upload UI for business listings
- [ ] Write tests for upload flow

### Phase 15 -- Search and API Versioning

- [ ] Implement full-text search for businesses (Postgres `ILIKE` or `tsvector`)
- [ ] Implement full-text search for events
- [ ] Verify all routes are prefixed with `/api/v1/`
- [ ] Document API versioning strategy
- [ ] Write search tests

### Phase 16 -- Dockerize and Deploy

- [ ] Write Dockerfile for Go API
- [ ] Push to Railway
- [ ] Migrate local Postgres data to Railway managed Postgres
- [ ] Update app to point at production API
- [ ] Verify production deployment end-to-end

### Phase 17 -- Production Polish and Device Testing

- [ ] UI polish, loading states, error handling
- [ ] Test on real iOS device (MBP M4 Pro + Xcode)
- [ ] Test on real Android device
- [ ] Performance testing on lower-end Android device
- [ ] Fix platform-specific issues

### Phase 18 -- App Store Submission

- [ ] Prepare app store listings (screenshots, descriptions)
- [ ] Write Terms of Service and Privacy Policy
- [ ] Submit to Google Play
- [ ] Submit to Apple App Store

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

- Do not scaffold anything not listed in the current phase.
- Complete the current phase before starting the next.
- Every sub-task must include tests before it is considered complete.
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
