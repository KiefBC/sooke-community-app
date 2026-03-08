---
title: "ADR-003: Use Chi over Fiber for Go HTTP Routing"
---


We chose Chi as the Go HTTP routing framework instead of Fiber.

## Status

Accepted

## Context

The Go backend needs an HTTP router for the REST API. The two primary candidates were Fiber and Chi.

Fiber is built on `fasthttp`, a custom HTTP engine that does not implement Go's standard `net/http` interfaces. It has an Express.js-like API that feels familiar to web developers coming from Node.js. It is faster than `net/http` in synthetic benchmarks.

Chi is built on Go's standard `net/http`. It is lightweight, idiomatic, and compatible with any Go middleware or library that targets `net/http`.

## Decision

We chose Chi because:

- It is built on `net/http`. Any standard Go middleware works without adaptation. This matters for Clerk JWT validation, CORS handling, rate limiting, and any future middleware.
- It is idiomatic Go. When searching "how to do X in Go HTTP," the answers work with Chi directly.
- At the scale of this app (small-town community, hundreds of users at most), Fiber's performance advantage is irrelevant.
- Chi avoids ecosystem lock-in. Fiber's `fasthttp` base means standard Go libraries often need wrappers or adapters.
- "Boring and predictable" is a strength for a solo developer maintaining a project long-term.

## Consequences

- Standard Go HTTP patterns and middleware work out of the box.
- Fiber's Express-like convenience (e.g., `c.JSON()`, `c.Params()`) is not available. Chi uses standard `http.ResponseWriter` and `*http.Request` patterns.
- If the API ever needs to handle extremely high throughput (unlikely at this scale), Fiber could be reconsidered. But at that point, the bottleneck would be the database, not the HTTP router.
