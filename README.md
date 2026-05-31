# raizel

[![check](https://github.com/rjansen/raizel/actions/workflows/check.yml/badge.svg)](https://github.com/rjansen/raizel/actions/workflows/check.yml)
[![test](https://github.com/rjansen/raizel/actions/workflows/test.yml/badge.svg)](https://github.com/rjansen/raizel/actions/workflows/test.yml)

A small, generic Go database-access helper that maps query results into
strongly-typed structs without code generation. Built on top of
`database/sql`, it works with any standard driver — and is tested against
**SQLite**, **PostgreSQL**, and **Oracle** (ATP and Free).

## Why

`database/sql` is great but verbose: every read needs a hand-written
`rows.Scan` loop, every nullable column needs a `sql.Null*` shuttle,
every parameter is positional and dialect-specific. Code generators like
sqlc solve this but require a build step and only support a fixed set of
dialects.

raizel takes a different angle: **runtime reflection with a generic API**.
You write SQL by hand, you tag your structs with `db:"col"`, and you call
generic helpers that figure out the rest:

```go
type User struct {
    ID        int64      `db:"id"`
    Email     string     `db:"email"`
    Labels    []string   `db:"labels"`     // native array / JSON, per dialect
    DeletedAt *time.Time `db:"deleted_at"` // NULL ↔ nil
}

rz, err := raizel.Open(raizel.DialectPostgres, dsn) // dialect bound once
users, err := raizel.Query[User](ctx, rz, "SELECT id, email, labels, deleted_at FROM users")
```

## Features

- `Query[T]` / `QueryOne[T]` / `Exec` / `ExecNamed[T]` — one helper per intent
- The dialect is configured **once** on a `*raizel.DB` [Handle] ([Open] maps the dialect to its driver; [Wrap] adopts a pool you tuned) — no per-call dialect noise
- A `*raizel.DB` or its `*raizel.Tx` (`rz.Begin(ctx)`) both satisfy `Handle`, so the same code runs in or out of a transaction
- Reflect-based row scanning with `db:"col"` tags, cached per type
- Nullable pointer fields (`*time.Time`, `*int64`, `*float64`, `*string`, `*bool`) round-trip to/from SQL NULL transparently
- Slice fields (`[]string`, `[]int64`, …) round-trip as **native arrays on PostgreSQL** (`TEXT[]`) and as **JSON text on Oracle/SQLite** — transparently, same struct
- Nested struct scanning via dot-notation column aliases (`SELECT t.id "team.id"` → `Member.Team.ID`)
- Named parameters in SQL (`:user_id`) extracted from struct tags and rewritten to the right placeholder per dialect (`?` / `$1` / `:1`)
- Multi-model `ExecNamed` batches wrap themselves in a transaction automatically
- Case-insensitive column lookup so Oracle's uppercase defaults Just Work
- Tag options supported (`db:"col,opt"`)
- Zero runtime dependencies — only `database/sql` + `reflect`

## Install

```sh
go get github.com/rjansen/raizel
```

## Quick start by dialect

### SQLite (any driver — example uses pure-Go `modernc.org/sqlite`)

```go
import (
    "context"
    "database/sql"
    "time"

    "github.com/rjansen/raizel"
    _ "modernc.org/sqlite"
)

type Bill struct {
    ID     int64      `db:"id"`
    Title  string     `db:"title"`
    Amount float64    `db:"amount"`
    PaidAt *time.Time `db:"paid_at"`
}

func main() {
    // raizel.Open maps the dialect to its driver ("sqlite") and binds the
    // dialect to the handle. The matching driver must be blank-imported.
    rz, _ := raizel.Open(raizel.DialectSQLite, "bills.db")
    defer rz.Close()
    ctx := context.Background()

    bill, err := raizel.QueryOne[Bill](ctx, rz,
        "SELECT id, title, amount, paid_at FROM bills WHERE id = ?", 42)
    if err == raizel.ErrNotFound { /* handle */ }

    when := time.Now()
    _, err = raizel.ExecNamed(ctx, rz,
        "UPDATE bills SET amount = :amount, paid_at = :paid_at WHERE id = :id",
        Bill{ID: 42, Amount: 1500, PaidAt: &when})
}
```

To adopt a pool you opened and tuned yourself (custom session params,
secrets injection), use `raizel.Wrap(sqlDB, dialect)` instead of `Open`.

### PostgreSQL (`pgx` via `database/sql`)

```go
import (
    "github.com/rjansen/raizel"
    _ "github.com/jackc/pgx/v5/stdlib"
)

rz, _ := raizel.Open(raizel.DialectPostgres, "postgres://user:pass@host:5432/db?sslmode=disable")

// The handle's dialect rewrites :name to $1, $2, ... for postgres.
_, err := raizel.ExecNamed(ctx, rz,
    `INSERT INTO bills (title, amount, paid_at) VALUES (:title, :amount, :paid_at)`,
    Bill{Title: "rent", Amount: 1500, PaidAt: &when})
```

### Oracle (ATP or Free, pure-Go `go-ora`)

```go
import (
    "github.com/rjansen/raizel"
    _ "github.com/sijms/go-ora/v2"
)

rz, _ := raizel.Open(raizel.DialectOracle, "oracle://user:pass@host:1521/FREEPDB1")

// :name → :1, :2, ... — Oracle's positional bind syntax.
_, err := raizel.ExecNamed(ctx, rz,
    `INSERT INTO bills (title, amount, paid_at) VALUES (:title, :amount, :paid_at)`,
    Bill{Title: "rent", Amount: 1500, PaidAt: &when})
```

## Testing

- `make test` — unit tests against in-memory SQLite (fast, no Docker)
- `make integration` — spins up the docker-compose stack and runs the same scenarios against PostgreSQL and Oracle Free
- `make check` — `gofmt`, `go vet`, `golangci-lint`

See [`docs/TESTING.md`](docs/TESTING.md) for layout, build tags, and CI details.

## Architecture

See [`docs/ARCHITECTURE.md`](docs/ARCHITECTURE.md) for the scanner pipeline,
nullable shuttling, named-param tokenizer, and cache/concurrency notes.

## Status

Active development on the `rjansen/revamp_database_access` branch. The
historical `firestore`/`cassandra`/`spanner`/`sql` layout has been
replaced with this single focused package.

## License

MIT — see [LICENSE](./LICENSE).
