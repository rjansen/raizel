# raizel

Generic Go database-access helper over `database/sql`. Single package at the
module root. Targets SQLite, PostgreSQL, Oracle (ATP/Free).

> **Status:** stub — generated during revamp. Final pass uses the
> `/create-claude-md` skill in Round 7.

## Commands

- `make test` — unit tests against `:memory:` SQLite
- `make integration` — `docker compose up -d` then build-tag-gated tests against postgres + oracle
- `make check` — gofmt, go vet, golangci-lint
- `make coverage` — coverage report

## Architecture

Single package `raizel` at the repo root. Files split by concern:

- `dialect.go` — `Dialect` enum + placeholder rewriting
- `querier.go` — `Querier` interface (`*sql.DB` or `*sql.Tx`)
- `errors.go` — sentinels (`ErrNotFound`)
- `scanner.go` — reflect-based row scanner, per-type cache
- `null.go` — nullable-pointer holders
- `named.go` — `:name` tokenizer + struct-tag value extraction
- `raizel.go` — public API: `Query[T]`, `QueryOne[T]`, `Exec`, `ExecNamed[T]`

## Code Style

- No external runtime dependencies — only `database/sql` + `reflect`
- Test deps allowed: SQLite driver for unit tests, postgres/oracle drivers under `//go:build integration`
- Per-type field-map cache via `sync.Map` keyed on `reflect.Type`
- Column lookups always lowercase the driver-reported column name
- Nullable `*T` round-trips to SQL NULL via `nullHolder` indirections
- Empty result slice is `[]T{}`, never `nil`

## Important

- NEVER add a runtime dependency outside the standard library
- NEVER bypass `Querier` by typing on `*sql.DB` directly — that breaks transaction support
- NEVER swallow errors from `rows.Err()` after the scan loop
- The `db:"-"` tag means "skip"; the `db:"col,opt"` form takes only the first segment as the column name
