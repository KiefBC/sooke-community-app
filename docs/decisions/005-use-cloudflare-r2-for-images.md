---
title: "ADR-005: Use Cloudflare R2 over AWS S3 for Image Storage"
---


We chose Cloudflare R2 as the object storage provider for images and uploaded media.

## Status

Accepted

## Context

The app needs to store images for business listings (logos, photos) and potentially event images in later phases. The standard approach is to upload images to an object storage service and store the URL in Postgres.

The developer initially considered AWS S3. We identified Cloudflare R2 as a better fit.

## Decision

We chose Cloudflare R2 because:

- R2 uses an S3-compatible API. The same SDK and code patterns work for both. Migration to or from S3 requires minimal code changes.
- R2 has no egress fees. S3 charges for every image download. For a community app where many users view the same business photos, egress costs add up on S3 but are zero on R2.
- R2's free tier is generous: 10GB storage and 10 million reads per month. For a small community app, the cost may remain zero indefinitely.
- The developer already uses Cloudflare for self-hosting and is familiar with the platform.

## Consequences

- Images are uploaded to R2, and the resulting URL is stored in Postgres.
- Images are served directly from R2 (or via Cloudflare's CDN).
- If we ever need to migrate to S3 or another S3-compatible service, the code changes are minimal because R2 uses the same API.
- Do not store images in the database. Always use object storage.
