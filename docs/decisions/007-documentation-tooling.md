---
title: "ADR-007: Use Starlight for Documentation"
---


We chose Starlight (Astro-based) as the documentation framework, hosted on Cloudflare Pages.

## Status

Accepted

## Context

The developer values clear documentation from the start and wants to avoid documentation debt. Documentation should be accessible to future team members, business owners, and contributors.

We evaluated five options:

1. **Plain markdown in `docs/`** -- zero tooling, works everywhere, but no search, no sidebar, harder to navigate as docs grow.
2. **mdBook** -- Rust-based, simple, self-hostable. Less polished output. Limited customization.
3. **VitePress** -- modern, polished, self-hostable. Large community (Vue/Vite ecosystem).
4. **Mintlify** -- polished, API docs features, Algolia search. But not self-hostable. Vendor lock-in.
5. **Starlight** -- Astro-based, purpose-built for documentation, self-hostable, modern.

Mintlify was ruled out because it is a hosted SaaS platform. Documentation lives on Mintlify's infrastructure and cannot be exported as a static site. This conflicts with the project's preference for owning infrastructure.

## Decision

We chose Starlight because:

- It is purpose-built for documentation sites with features like sidebar navigation, search, versioning, and i18n support built in.
- It produces static output that can be hosted on Cloudflare Pages (which the developer already uses).
- No vendor lock-in. Documentation source files are standard markdown.
- The developer reviewed Starlight and liked the visual quality and feature set.
- Astro's build system is fast and the framework is actively maintained with a growing community.

## Consequences

- Starlight adds a Node dependency for the documentation site.
- Documentation is written in markdown and built into a static site on each deploy.
- The documentation site is separate from the mobile app and admin dashboard.
- All three web properties (docs, admin dashboard, and potentially a future marketing site) can be hosted on Cloudflare Pages.
- If Starlight becomes unmaintained, VitePress is a comparable alternative with a larger community.
