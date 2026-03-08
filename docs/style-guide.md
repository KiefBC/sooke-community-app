---
title: "Documentation Style Guide"
---


This document defines the formatting, writing, and structural rules for all documentation in the Sooke Community App project. All contributors must follow these rules to keep documentation consistent and clear.

---

## Formatting Rules

- Do not use em-dashes. Use two hyphens (`--`) or restructure the sentence.
- Do not use emojis in documentation or code comments.
- Use ATX-style headers (`## Header`, not underline style).
- Leave one blank line before and after headers.
- Use pipe (`|`) tables for structured comparisons. Do not use HTML tables.
- Code blocks must specify a language (` ```go `, ` ```svelte `, ` ```sql `, etc.).
- Wrap file paths and terminal commands in backticks (`` `docs/project-plan.md` ``).
- Use hyphens (`-`) for unordered lists, not asterisks.
- Use `[ ]` and `[x]` for task/checklist items.

---

## Writing Rules

- Use active voice over passive voice.
- Write short sentences. One idea per sentence.
- Use "we" for project decisions.
- Use "the user" for app users.
- Use "the developer" or "Super Admin" for the project maintainer.
- Define acronyms on first use. Example: "Architecture Decision Record (ADR)."
- Do not use marketing language or superlatives ("blazing fast", "cutting edge", etc.).
- Document the "why" alongside the "what."
- Be direct. Do not hedge with phrases like "it might be worth considering."

---

## Structure Rules

- Every document starts with a title (`# Title`) and a one-sentence summary.
- Architecture Decision Records (ADRs) follow a fixed template: Context, Decision, Consequences.
- `docs/project-plan.md` is the source of truth for current decisions.
- `docs/future-ideas-and-alternatives.md` is the source of truth for deferred and rejected options.
- `docs/planning-discussion.md` is append-only. Do not edit past entries.

---

## File Naming Rules

- Use lowercase letters only.
- Use hyphens for spaces (`project-plan.md`, not `ProjectPlan.md` or `project_plan.md`).
- Number ADRs with zero-padded three-digit prefixes: `001-use-capacitor-over-tauri.md`.
- Do not use spaces in file names.

---

## ADR Template

All ADRs must follow this structure:

```markdown
# ADR-NNN: Title

Summary sentence describing the decision.

## Status

Accepted | Superseded | Deprecated

## Context

What problem or question prompted this decision?

## Decision

What did we decide, and why?

## Consequences

What are the tradeoffs, risks, or follow-up actions?
```

---

## References

- This style guide applies to all files in the `docs/` directory.
- Code comments in the application should also follow these writing rules where practical.
- The project plan references this guide as the authority on documentation standards.
