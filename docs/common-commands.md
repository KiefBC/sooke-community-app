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

- Go 1.24 or later

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

## Sections to Add Later

The following sections will be added as each tool or service is set up:

- Go API (run, test, build, lint)
- SvelteKit mobile app (run, test, build)
- Capacitor (iOS build, Android build, sync)
- PostgreSQL (Docker, migrations, seed data)
- Admin dashboard (run, test, build, deploy)
- Cloudflare R2 (upload, configure)
- Clerk (configure, test tokens)
