# raizel

A small, generic Go database-access helper that maps query results into
strongly-typed structs without code generation. Built on top of
`database/sql`, it works with any standard driver — and is tested against
**SQLite**, **PostgreSQL**, and **Oracle** (ATP and Free).

## Why

`database/sql` is great but verbose: every read needs a hand-written
`rows.Scan` loop, every nullable column needs a `sql.Null*` shuttle, every
parameter is positional and dialect-specific. Code generators like sqlc
solve this but require a build step and only support a fixed set of
dialects.

`raizel` takes a different angle: **runtime reflection with a generic API**.
You write SQL by hand, you tag your structs with `db:"col"`, and you call
generic helpers that figure out the rest:

```go
type User struct {
    ID        int64      `db:"id"`
    Email     string     `db:"email"`
    DeletedAt *time.Time `db:"deleted_at"` // nullable, scans NULL → nil
}

users, err := raizel.Query[User](ctx, db, "SELECT id, email, deleted_at FROM users")
```

## Features

- `Query[T]` / `QueryOne[T]` / `Exec` / `ExecNamed[T]` — one helper per intent
- Works with `*sql.DB` **or** `*sql.Tx` via the `Querier` interface — same code in or out of a transaction
- Reflect-based row scanning with `db:"col"` tags, cached per type
- Nullable pointer fields (`*time.Time`, `*int64`, `*float64`, `*string`, `*bool`) round-trip to/from SQL NULL transparently
- Nested struct scanning via dot-notation column aliases (`SELECT t.id "team.id"` → `Member.Team.ID`)
- Named parameters in SQL (`:user_id`) extracted from struct tags and rewritten to the right placeholder per dialect (`?` / `$1` / `:1`)
- Multi-model `ExecNamed` batches wrap themselves in a transaction automatically
- Case-insensitive column lookup so Oracle's uppercase defaults Just Work
- Tag options supported (`db:"col,opt"`)

## Status

Active development on the `rjansen/revamp_database_access` branch. The
historical `firestore`/`cassandra`/`spanner`/`sql` layout has been replaced
with this single focused package.

## Quick start

```go
import (
    "context"
    "database/sql"

    "github.com/rjansen/raizel"
    _ "github.com/mattn/go-sqlite3"
)

type Bill struct {
    ID     int64    `db:"id"`
    Title  string   `db:"title"`
    Amount float64  `db:"amount"`
    PaidAt *time.Time `db:"paid_at"`
}

func main() {
    db, _ := sql.Open("sqlite3", "bills.db")
    ctx := context.Background()

    bill, err := raizel.QueryOne[Bill](ctx, db,
        "SELECT id, title, amount, paid_at FROM bills WHERE id = ?", 42)
    if err == raizel.ErrNotFound { /* ... */ }

    _, err = raizel.ExecNamed(ctx, db, raizel.DialectSQLite,
        "UPDATE bills SET amount = :amount, paid_at = :paid_at WHERE id = :id", bill)
}
```

See [`CLAUDE.md`](./CLAUDE.md) for architecture and contributor guidelines.

## Testing

- `make test` — unit tests against in-memory SQLite (fast, no Docker)
- `make integration` — spins up the docker-compose stack and runs the same scenarios against PostgreSQL and Oracle Free
- `make check` — `gofmt`, `go vet`, `golangci-lint`

## License

MIT — see [LICENSE](./LICENSE).
