---
title: "Database Schema Reference"
---

This document defines every table in the Sooke Community App database, what each column is for, and why each design decision was made.

---

## Overview

The schema has 9 tables, created in dependency order by goose migrations. Every table that references another table is created after the table it depends on.

```
roles
  -> users
       -> businesses (owner_id)
            -> business_hours
            -> menus
                 -> menu_items
       -> events (submitted_by)
business_categories
  -> businesses (category_id)
event_types
  -> events (event_type_id)
```

---

## Tables

### roles

Defines the permission tiers for the application. Created first because `users` references it.

```sql
CREATE TABLE roles (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `id` | BIGSERIAL | Auto-incrementing primary key. BIGSERIAL over SERIAL to avoid the 2.1B ceiling. |
| `name` | TEXT NOT NULL UNIQUE | The role name (`general_user`, `business_owner`, `super_admin`). UNIQUE prevents duplicates. |
| `created_at` | TIMESTAMPTZ | When the row was created. TIMESTAMPTZ stores timezone-aware timestamps to avoid UTC bugs. |
| `updated_at` | TIMESTAMPTZ | When the row was last modified. |

**Expected rows:** `general_user`, `business_owner`, `super_admin`. These are seeded, not migrated.

---

### users

Links a Clerk identity to app-specific data. Without this table, the API has no way to know who is making a request or what role they have.

```sql
CREATE TABLE users (
  id BIGSERIAL PRIMARY KEY,
  clerk_id TEXT NOT NULL UNIQUE,
  role_id BIGINT NOT NULL REFERENCES roles(id),
  email TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `clerk_id` | TEXT NOT NULL UNIQUE | The user's ID from Clerk. Every authenticated API request looks up the user by this field. The UNIQUE constraint creates an implicit index for fast lookups. |
| `role_id` | BIGINT NOT NULL REFERENCES roles(id) | Links the user to their role. NOT NULL because every user must have a role. |
| `email` | TEXT NOT NULL UNIQUE | The user's email from Clerk. UNIQUE prevents duplicate accounts. |
| `display_name` | TEXT NOT NULL | What the user sees in the UI. |

---

### business_categories

The curated tag list for businesses. A separate table (not a Postgres ENUM) so the Super Admin can add categories at runtime through the admin dashboard without requiring a database migration.

```sql
CREATE TABLE business_categories (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  slug TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `name` | TEXT NOT NULL UNIQUE | What the user sees ("Restaurant", "Cafe", "Retail"). |
| `slug` | TEXT NOT NULL UNIQUE | What the API uses for filtering (`/api/v1/businesses?category=restaurant`). |

---

### businesses

The core entity. Everything in the app orbits around local businesses.

```sql
CREATE TABLE businesses (
  id BIGSERIAL PRIMARY KEY,
  owner_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
  category_id BIGINT NOT NULL REFERENCES business_categories(id),
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  description TEXT,
  phone TEXT,
  email TEXT,
  website TEXT,
  address TEXT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lng DOUBLE PRECISION NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `owner_id` | BIGINT, nullable, ON DELETE SET NULL | The Super Admin creates the business first. The owner is assigned later when a real person is verified. Nullable because not every business has a claimed owner. SET NULL so the business survives if the user account is deleted. |
| `category_id` | BIGINT NOT NULL | Every business needs a category for filtering to work. |
| `slug` | TEXT NOT NULL UNIQUE | URL-friendly identifier (`joes-coffee-shop`). Used in API responses and prepared for future deep linking. |
| `description` | TEXT, nullable | Optional blurb. Some businesses will not have one initially. |
| `phone`, `email`, `website` | TEXT, nullable | Not every business has all three contact methods. |
| `address` | TEXT NOT NULL | Street address for display. Every business has a physical location. |
| `lat`, `lng` | DOUBLE PRECISION NOT NULL | Map pin coordinates. DOUBLE PRECISION gives sub-millimeter accuracy at city scale. PostGIS is overkill for pin-dropping. NOT NULL because every business needs a map pin. |

---

### business_hours

A separate table instead of a JSON blob because hours are complex (per-day, seasonal, holidays) and need to be queryable ("what is open right now?").

```sql
CREATE TABLE business_hours (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  day_of_week SMALLINT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
  open_time TIME NOT NULL,
  close_time TIME NOT NULL,
  is_closed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE (business_id, day_of_week)
);
```

| Column | Type | Why |
|--------|------|-----|
| `business_id` | BIGINT NOT NULL, ON DELETE CASCADE | If the business is removed, its hours go with it. No orphaned rows. |
| `day_of_week` | SMALLINT, CHECK 0--6 | 0 = Sunday, 6 = Saturday. The CHECK constraint prevents invalid values like `9`. |
| `open_time`, `close_time` | TIME NOT NULL | Opening and closing times for that day. |
| `is_closed` | BOOLEAN, default FALSE | Handles days the business is closed. When true, `open_time`/`close_time` are ignored by the application. |
| UNIQUE (business_id, day_of_week) | Constraint | One row per business per day. Prevents duplicate entries for the same day. |

---

### menus

A named container for menu items. One business can have multiple menus (lunch, dinner, drinks).

```sql
CREATE TABLE menus (
  id BIGSERIAL PRIMARY KEY,
  business_id BIGINT NOT NULL REFERENCES businesses(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `business_id` | BIGINT NOT NULL, ON DELETE CASCADE | If the business is deleted, its menus go with it. |
| `name` | TEXT NOT NULL | "Lunch Menu", "Dinner Menu", "Drinks", etc. |

---

### menu_items

Individual items within a menu.

```sql
CREATE TABLE menu_items (
  id BIGSERIAL PRIMARY KEY,
  menu_id BIGINT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  description TEXT,
  price NUMERIC(10,2) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `menu_id` | BIGINT NOT NULL, ON DELETE CASCADE | If the menu is deleted, its items go with it. Also cascades transitively -- if the business is deleted, menus are deleted, which deletes items. |
| `description` | TEXT, nullable | Optional. Some items are self-explanatory ("Coffee"). |
| `price` | NUMERIC(10,2) NOT NULL | Exact decimal arithmetic. FLOAT would give $12.989999 instead of $12.99. Supports values up to $99,999,999.99. |

---

### event_types

The curated tag list for events. Same pattern as `business_categories` -- a separate table so the Super Admin can manage event types at runtime.

```sql
CREATE TABLE event_types (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  slug TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

| Column | Type | Why |
|--------|------|-----|
| `name` | TEXT NOT NULL UNIQUE | What the user sees ("Live Music", "Market", "Workshop"). |
| `slug` | TEXT NOT NULL UNIQUE | What the API uses for filtering. |

---

### events

The second core entity. Events are submitted by users and require approval before they are visible.

```sql
CREATE TABLE events (
  id BIGSERIAL PRIMARY KEY,
  event_type_id BIGINT NOT NULL REFERENCES event_types(id),
  submitted_by BIGINT NOT NULL REFERENCES users(id),
  business_id BIGINT REFERENCES businesses(id) ON DELETE SET NULL,
  name TEXT NOT NULL,
  slug TEXT NOT NULL UNIQUE,
  description TEXT,
  lat DOUBLE PRECISION,
  lng DOUBLE PRECISION,
  starts_at TIMESTAMPTZ NOT NULL,
  ends_at TIMESTAMPTZ,
  status TEXT NOT NULL DEFAULT 'draft'
    CHECK (status IN ('draft', 'pending_review', 'approved', 'rejected')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT event_location_check CHECK (
    (business_id IS NOT NULL AND lat IS NULL AND lng IS NULL) OR
    (business_id IS NULL AND lat IS NOT NULL AND lng IS NOT NULL)
  )
);
```

| Column | Type | Why |
|--------|------|-----|
| `event_type_id` | BIGINT NOT NULL | Every event must have a type for filtering. |
| `submitted_by` | BIGINT NOT NULL | Tracks who created the event. Needed for rate limiting (max 5 pending per user) and notifications. |
| `business_id` | BIGINT, nullable, ON DELETE SET NULL | Nullable because events at public locations do not have a business. SET NULL so the event survives if the business is later deleted. |
| `lat`, `lng` | DOUBLE PRECISION, nullable | Nullable because events at a business use the business's coordinates instead. |
| `starts_at` | TIMESTAMPTZ NOT NULL | When the event starts. Required. |
| `ends_at` | TIMESTAMPTZ, nullable | When the event ends. Optional because some events do not have a fixed end time. |
| `status` | TEXT with CHECK | Uses a CHECK constraint instead of a Postgres ENUM. ENUMs are painful to modify (requires `ALTER TYPE`). A CHECK is easier to extend by dropping and re-adding the constraint. |
| `event_location_check` | Constraint | The XOR rule: an event has EITHER a `business_id` OR its own lat/lng, never both, never neither. This is a business rule enforced at the database level so no buggy API code can violate it. |

---

## Design Decisions Summary

| Decision | Why |
|----------|-----|
| BIGSERIAL over SERIAL | Avoids the 2.1 billion row ceiling. Costs nothing extra. |
| TIMESTAMPTZ over TIMESTAMP | Always stores timezone-aware timestamps. Prevents bugs when servers run in different timezones. |
| DOUBLE PRECISION over PostGIS | Sub-millimeter accuracy at city scale. PostGIS adds complexity with no benefit for pin-dropping. |
| NUMERIC(10,2) over FLOAT | Exact decimal arithmetic. $12.99 stays $12.99. |
| TEXT CHECK over ENUM | Easier to modify. Adding a new status is a constraint swap, not an `ALTER TYPE`. |
| CASCADE on child tables | business_hours, menus, menu_items are deleted when their parent is deleted. No orphaned rows. |
| SET NULL on optional references | Events survive business deletion. Businesses survive user deletion. Data is preserved. |
| Separate tables for categories/types | Super Admin can add new tags at runtime without a migration. |
| Nullable owner_id | Businesses are created by the admin before any owner is assigned. |

---

## Related Documents

- [ADR-009: Use goose for migrations](/decisions/009-use-goose-for-migrations/) -- why we chose goose as the migration tool.
- [ADR-008: NAS Postgres over Docker](/decisions/008-nas-postgres-over-docker/) -- why we use NAS-hosted Postgres for local development.
- [Project Plan](/project-plan/) -- full development phases and sub-task checklists.
