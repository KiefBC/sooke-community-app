---
title: "ADR-006: Use Curated Tags over Free-Form User Tags"
---


We chose a curated tag system managed by the Super Admin instead of allowing users to create custom tags.

## Status

Accepted

## Context

Both businesses and events need categorization for filtering and discovery. The question was whether users should be able to create their own tags or select from a predefined list.

Free-form tags offer flexibility but create moderation challenges. The Go library `goaway` (github.com/TwiN/go-away) provides profanity detection with leet-speak support, but profanity filters are blunt instruments. They cannot distinguish context. For example, "gay" as identity vs. "gay" as a slur. False positives are embarrassing for a community app. False negatives are harmful.

## Decision

We chose curated tags because:

- A predefined list eliminates the moderation problem entirely. No free-form input means no offensive content.
- The Super Admin defines and manages the tag list via the admin dashboard.
- Business categories and event types are predictable and finite for a small town. A list of 15-20 tags per entity type covers the vast majority of cases.
- If a user or business owner needs a tag that does not exist, they can request it. The Super Admin reviews and adds it if appropriate.

Examples of business categories: Restaurant, Cafe, Bar, Retail, Grocery, Outdoor Recreation, Health and Wellness, Arts and Culture, Professional Services, Accommodation.

Examples of event types: Live Music, Market, Workshop, Community Meeting, Sports, Festival, Fundraiser, Kids and Family, Outdoor.

## Consequences

- Users cannot create custom tags. This may occasionally require Super Admin intervention to add new categories.
- The tag list is stored in the database and managed through the admin dashboard.
- If the curated approach proves too restrictive, we can revisit free-form tags with the `goaway` profanity filter and an admin review queue. This alternative is documented in [future-ideas-and-alternatives.md](/future-ideas-and-alternatives/).
