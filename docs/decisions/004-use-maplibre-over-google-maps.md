---
title: "ADR-004: Use MapLibre over Google Maps"
---


We chose MapLibre GL JS with MapTiler as the tile provider instead of Google Maps JS SDK.

## Status

Accepted

## Context

The app needs a map to pin business locations and event venues in Sooke. The map requirements are straightforward: display a map of the Sooke area, place pins on it, and allow users to tap pins to view details.

The original plan specified Google Maps JS SDK. We identified several concerns:

- API key restrictions are complex in a mobile webview. Google Maps restricts keys by HTTP referrer, but a local webview does not have a standard referrer.
- Google Maps is overkill for pin-dropping on a small-town map.
- The JS SDK can feel sluggish in a mobile webview compared to lighter alternatives.
- While the free tier ($200/month in credits) is generous, there are fully free alternatives.

We evaluated three options: Google Maps JS SDK, Leaflet + OpenStreetMap, and MapLibre GL JS.

## Decision

We chose MapLibre GL JS with MapTiler because:

- MapLibre is open-source and renders vector tiles with smooth zooming and a modern look.
- MapTiler's free tier provides 100,000 tile loads per month, which is more than sufficient for a small-town community app.
- No API key referrer restrictions to work around in the webview.
- MapLibre is just a renderer. The tile source is pluggable. If MapTiler changes pricing, we can swap to another tile provider or self-host tiles without changing application code.
- The developer reviewed MapLibre and liked the visual quality.

Leaflet + OpenStreetMap was considered but MapLibre offers better-looking vector tiles and smoother performance.

## Consequences

- We depend on MapTiler for tile hosting. If the free tier becomes insufficient, we can switch to another provider or self-host tiles. Self-hosting is documented in [future-ideas-and-alternatives.md](/future-ideas-and-alternatives/).
- MapLibre is embedded in the Capacitor webview via JavaScript. No native maps plugin is used.
- Google Maps remains a documented fallback if MapLibre or MapTiler does not meet quality expectations.
