---
title: "ADR-008: Use NAS-Hosted Postgres over Docker for Local Development"
---

We chose to connect to a Postgres instance on the developer's NAS instead of running Postgres in a local Docker container for development.

## Status

Accepted

## Context

The project plan originally specified "PostgreSQL running in Docker" for local development, with a `docker-compose.yml` to manage the container. The intent was to make the database portable and easy to spin up.

However, the developer already runs a Postgres instance on a Synology NAS on the local network. This instance is always available when the developer is working. Setting up Docker for Postgres would add complexity without a clear benefit in this situation.

The key question: does the convenience of an existing, always-on Postgres instance outweigh the portability of a Dockerized database?

## Decision

We use the NAS-hosted Postgres for local development because:

- The database is already running and configured. Zero additional setup.
- The developer manages it via TablePlus, which is already installed.
- The `.env` config pattern means switching between NAS, Docker, or Railway is a single environment variable change (`DATABASE_URL`). The Go code does not know or care where Postgres is running.
- Docker adds a dependency and resource usage (CPU, memory, disk) that provides no benefit when a dedicated database server is already available on the network.
- Production still uses Railway managed Postgres as planned. The environment variable swap is the same regardless of whether local dev uses Docker or NAS.

## Consequences

- There is no `docker-compose.yml` in the repository. Contributors without access to the NAS must set up their own Postgres instance (local install, Docker, or hosted). The `.env.example` file documents the required environment variables.
- The NAS must be reachable on the local network for the developer to run integration tests. If the NAS is offline, integration tests are skipped (they check for `TEST_DATABASE_URL` and skip if unset).
- If a second developer joins the project and does not have NAS access, adding a `docker-compose.yml` at that point is straightforward. The application code requires no changes -- only the `DATABASE_URL` value differs.
- Backup and maintenance of the development database is the developer's responsibility (NAS-level backups, Postgres `pg_dump`, etc.).
