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
go test ./...
```

This runs all unit tests. Integration tests that require a database are skipped when `TEST_DATABASE_URL` is not set.

### Run tests with database (integration)

```bash
cd server
TEST_DATABASE_URL="<your_test_database_url>" go test ./...
```

This runs the full test suite including integration tests that verify database connectivity. Replace `<password>` with the actual password for the `sooke_app` user.

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
curl http://localhost:8989/api/v1/health
```

Expected response when the database is reachable:

```json
{"status":"ok","db_status":"connected"}
```

---

## Sections to Add Later

The following sections will be added as each tool or service is set up:

- SvelteKit mobile app (run, test, build)
- Capacitor (iOS build, Android build, sync)
- PostgreSQL migrations and seed data
- Admin dashboard (run, test, build, deploy)
- Cloudflare R2 (upload, configure)
- Clerk (configure, test tokens)
