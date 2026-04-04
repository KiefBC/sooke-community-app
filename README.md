# Sooke Community App


<div align="center">
  <video src="https://github.com/user-attachments/assets/08efda41-16d6-4182-a9b6-3815fa9db69d" autoplay loop muted playsinline></video>
</div>


A mobile-first community app for Sooke, BC. Built with SwiftUI (iOS), Go, and PostgreSQL.

## Prerequisites

- Go (see `server/go.mod` for version)
- Xcode 26.4+ with iOS 26 SDK
- XcodeGen: `brew install xcodegen`
- PostgreSQL
- goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

## iOS App

```bash
cd ios
xcodegen generate
open SookeCommunity.xcodeproj
```

XcodeGen generates the Xcode project from `project.yml`. Run `xcodegen generate` after adding or removing Swift files. The only external dependency is Kingfisher (image caching), pulled automatically via SPM.

### Tests

```bash
cd ios
xcodebuild test -scheme SookeCommunity -destination 'platform=iOS Simulator,name=iPhone 17 Pro'
```

Tests use the Swift Testing framework (`@Suite`, `@Test`). See `docs/swift-testing-guide.md` for patterns and conventions.

## Environment Setup

```bash
cp .env.example .env
```

Fill in `DATABASE_URL` and `TEST_DATABASE_URL`. The `.env` file is gitignored.

## Go Server

```bash
cd server
go mod download
go run ./cmd/api
```

The server runs on port 8080 by default. Override with the `PORT` env var.

### Build

```bash
cd server
go build -o bin/api ./cmd/api
```

### Tests

```bash
cd server
go test ./...
```

Integration tests require `TEST_DATABASE_URL`:

```bash
cd server
TEST_DATABASE_URL="<your_test_db_url>" go test ./...
```

## Database Migrations

Migrations auto-run at API startup. Manual commands for development:

```bash
# Run all pending migrations
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations up

# Roll back the last migration
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations down

# Check migration status
GOOSE_DRIVER=postgres GOOSE_DBSTRING="$DATABASE_URL" goose -dir server/migrations status

# Create a new migration
goose -dir server/migrations create <name> sql
```

## Seed Data

```bash
cd server
go run ./cmd/seed
```

Inserts sample Sooke businesses, categories, and event types. Idempotent.

## API Endpoints

All routes are prefixed with `/api/v1/`.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Health check -- returns app and database status |
| GET | `/api/v1/businesses` | List businesses (supports `?search=`, `?category=`, `?page=`, `?per_page=`) |
| GET | `/api/v1/businesses/{slug}` | Get a single business by slug |
| GET | `/api/v1/categories` | List all business categories |

