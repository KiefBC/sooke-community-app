---
title: "Common Commands"
---

This document lists the commands needed to run, build, and manage the project. It is updated as new tools and services are added.

---

## Documentation (Starlight)

The documentation site uses [Starlight](https://starlight.astro.build/), an Astro-based static site generator. The Starlight project lives in the `starlight-docs/` directory. The source markdown files are maintained in `docs/` and copied into `starlight-docs/src/content/docs/` with frontmatter added.

### Prerequisites

- Node.js 18 or later
- npm (included with Node.js)

### Install dependencies

```bash
cd starlight-docs
npm install
```

### Run the documentation site locally

```bash
cd starlight-docs
npm run dev
```

This starts a local development server with hot reload. Open the URL shown in the terminal (typically `http://localhost:4321`).

### Build the documentation site for production

```bash
cd starlight-docs
npm run build
```

This generates a static site in `starlight-docs/dist/`. The output can be deployed to Cloudflare Pages or any static hosting provider.

### Preview the production build locally

```bash
cd starlight-docs
npm run preview
```

This serves the built static site locally so you can verify the production output before deploying.

---

## Go Server

The backend API is written in Go using [Chi](https://github.com/go-chi/chi). The API lives in the server/ directory.

### Prerequisites

- A supported Go toolchain (see `server/go.mod` for the minimum required version)

### Install dependencies

```bash
cd server
go mod download
```

### Run API Locally

```bash
cd server
go run ./cmd/api
```

The server defaults to port 8080. Set the `PORT` environment variable to change the port.

### Build the binary

```bash
cd server
go build -o bin/api ./cmd/api
```

This will produce an executable binary at `server/bin/api` that can be run on the target platform.

### Run tests

```bash
cd server
go test ./... -p 1
```

This runs all unit tests. The `-p 1` flag runs test packages serially -- this is required because multiple packages (`database`, `repository`) share a single test database and reset its schema in `TestMain`. Without `-p 1`, concurrent schema drops cause race conditions.

Integration tests that require a database are skipped when `TEST_DATABASE_URL` is not set.

### Run tests with database (integration)

```bash
cd server
TEST_DATABASE_URL="<your_test_database_url>" go test ./... -p 1
```

This runs the full test suite including integration tests that verify database connectivity.

### Run tests for a single package

```bash
cd server
go test ./internal/repository/... -v
```

When running a single package, `-p 1` is not needed. Use `-v` for verbose output showing individual test case names.

---

## PostgreSQL

The development database runs on the developer's NAS. See [ADR-008](/decisions/008-nas-postgres-over-docker/) for why we chose this over Docker. Connection details (host, port, credentials) are configured via environment variables in `.env`.

The database client used for manual inspection is [TablePlus](https://tableplus.com/) on macOS.

### Environment setup

1. Copy `.env.example` to `.env` at the project root.
2. Fill in the real values for `DATABASE_URL` and `TEST_DATABASE_URL`.
3. The `.env` file is gitignored and must never be committed.

The Go API loads `.env` automatically at startup via `godotenv`. In production (Railway), environment variables are set directly -- no `.env` file is needed.

### Verify the connection

Start the API and hit the health endpoint:

```bash
cd server
go run ./cmd/api
```

```bash
curl http://localhost:8080/api/v1/health
```

Expected response when the database is reachable:

```json
{ "status": "ok", "db_status": "connected" }
```

---

## Database Migrations and Seed Data

We use [goose](https://github.com/pressly/goose) for schema migrations. See [ADR-009](/decisions/009-use-goose-for-migrations/) for why we chose goose.

Migration files live in `server/migrations/`. Migrations auto-run at API startup, so manual execution is usually unnecessary. The commands below are for development and troubleshooting.

### Prerequisites

- goose CLI: `go install github.com/pressly/goose/v3/cmd/goose@latest`

### Run all pending migrations

```bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations up
```

### Roll back the last migration

```bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations down
```

### Check migration status

```bash
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations status
```

### Create a new migration

```bash
goose -dir server/migrations create <name> sql
```

This creates a pair of timestamped `.sql` files (up and down) in `server/migrations/`.

### Seed the database with sample data

```bash
cd server
go run ./cmd/seed
```

This inserts sample Sooke businesses, categories, and event types for development. Requires `DATABASE_URL` to be set. The seed runner is idempotent -- it can be run multiple times without creating duplicates.

---

## Sections to Add Later

The following sections will be added as each tool or service is set up:

- SvelteKit mobile app (run, test, build)
- Capacitor (iOS build, Android build, sync)
- Admin dashboard (run, test, build, deploy)
- Cloudflare R2 (upload, configure)
- Clerk (configure, test tokens)
