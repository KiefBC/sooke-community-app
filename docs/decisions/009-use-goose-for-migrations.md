---
title: "ADR-009: Use Goose for Database Migrations"
---

We chose goose as the database migration tool for managing schema versioning, up/down migrations, and tracking applied changes.

## Status

Accepted

## Context

Phase 2 introduces the initial database schema -- nine tables covering businesses, events, users, roles, menus, and related entities. We need a migration tool that supports:

- Versioned up and down migrations in plain SQL.
- Tracking which migrations have been applied.
- Running migrations programmatically from Go at application startup.
- A CLI for manual migration management during development.

Three options were considered:

1. **goose** (`github.com/pressly/goose/v3`) -- SQL or Go migration files, works directly with `database/sql`, supports Go `embed.FS`, built-in CLI.
2. **golang-migrate** (`github.com/golang-migrate/migrate/v4`) -- SQL migration files, requires separate database driver and source driver packages.
3. **Raw SQL files without a runner** -- manual execution, no version tracking.

## Decision

We use goose because:

- It works directly with `database/sql` and requires only a single Go dependency. The project already registers the pgx driver, so goose needs no additional database driver package.
- It supports `embed.FS`, which lets us bundle migration files into the compiled binary. This means the production container does not need a separate `migrations/` directory at runtime.
- The Go API is clean and minimal: `goose.SetDialect("postgres")` followed by `goose.Up(db, dir)`. This integrates naturally into the existing `database` package.
- The built-in CLI (`goose -dir server/migrations up/down/status`) is useful during development for manual migration management and status inspection.
- It tracks applied migrations in a `goose_db_version` table. `Up()` is idempotent -- it only applies pending migrations and skips already-applied ones.
- golang-migrate was rejected because it requires separate packages for the database driver (`database/pgx/v5`) and file source (`source/file`), adding dependency complexity without additional benefit. Its Go API is also less ergonomic for our use case.
- Raw SQL without a runner was rejected because it provides no version tracking, no rollback automation, and no way to programmatically apply migrations at startup.

## Consequences

- `github.com/pressly/goose/v3` is added as a Go dependency in `server/go.mod`.
- Developers install the goose CLI via `go install github.com/pressly/goose/v3/cmd/goose@latest` for manual migration management.
- Migration files live in `server/migrations/` and follow goose's naming convention (sequential numbering with up/down suffixes).
- Migrations auto-run at API startup. This keeps development and production schemas in sync without manual steps.
- Goose creates and manages a `goose_db_version` table in the database. This table must not be modified manually.
- If we later need Go-based migrations (data transformations, backfills), goose supports `.go` migration files alongside `.sql` files without changing the tool.
