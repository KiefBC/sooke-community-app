# Sooke Community App

A mobile-first community app for Sooke, BC. Built with SvelteKit, Capacitor, Go, and PostgreSQL.

## Prerequisites

- Go (see `server/go.mod` for version)
- Node.js 18+
- PostgreSQL
- goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`

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

## Health Check

```bash
curl http://localhost:8080/api/v1/health
```

## Documentation Site (Starlight)

```bash
cd starlight-docs
npm install
npm run preview
```
