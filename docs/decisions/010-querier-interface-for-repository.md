---
title: "ADR-010: Use Querier Interface for Repository Functions"
---

We chose to define a `Querier` interface in the repository package so that repository functions accept either a database connection or a transaction.

## Status

Accepted

## Context

Phase 3 introduces the repository layer -- functions that run SQL queries and return Go structs. These functions need a database handle to execute queries against.

The initial approach used `*sql.DB` directly in function signatures. This works for production code, but creates a problem in tests: `*sql.DB` always operates outside a transaction, so test data inserted inside a transaction is invisible to the repository function. This prevents rollback-based test isolation -- each test would leave data behind that affects subsequent tests.

Both `*sql.DB` and `*sql.Tx` implement the same query methods (`QueryContext`, `QueryRowContext`, `ExecContext`), but Go does not automatically unify them. Without a shared interface, every repository function would need two versions or the tests would need a separate cleanup strategy.

## Decision

We define a `Querier` interface in the repository package:

```go
type Querier interface {
    QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
    QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
    ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
```

All repository functions accept `Querier` instead of `*sql.DB`. Both `*sql.DB` and `*sql.Tx` satisfy this interface without any wrapper code.

- **Production handlers** pass `*sql.DB` as before -- no change in behavior.
- **Tests** pass `*sql.Tx`, insert seed data, call the repository function, assert results, then `tx.Rollback()`. Each test is fully isolated.

## Consequences

- Every repository function uses `Querier` in its signature instead of `*sql.DB`. This is a one-word change per function (`db *sql.DB` becomes `q Querier`).
- Handlers still receive `*sql.DB` from the router and pass it to repository functions. No handler changes are needed.
- Tests can now use exact equality assertions (`==`) instead of "at least" comparisons, because each test only sees its own transaction's data.
- Future phases that add write operations (Phase 10 -- business owner editing, Phase 11 -- event submission) can wrap multiple repository calls in a single transaction for atomicity, passing the same `tx` to each function.
- The interface lives in the repository package, not in a separate `internal/db` package. If a second package needs it later, we can extract it then.
