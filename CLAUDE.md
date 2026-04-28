# raizel

Generic Go database-access helper on top of `database/sql`. Public API at the module root, implementation under `internal/`. Supports SQLite, PostgreSQL, and Oracle (ATP / Free).

## Commands

- `make test` — unit tests against `:memory:` SQLite (no Docker)
- `make integration` — `docker compose up`, build-tag-gated tests against postgres + oracle
- `make check` — gofmt + go vet + golangci-lint
- `make fmt` / `make lint` — auto-fix variants
- `make coverage` — coverage report

## Architecture

Public API (module root, `package raizel`):

- `raizel.go` — `Query[T]`, `QueryOne[T]`, `Exec`, `ExecNamed[T]`
- `querier.go` — `Querier` interface (`*sql.DB` or `*sql.Tx`)
- `dialect.go` — `Dialect` enum + placeholder rewriting
- `errors.go` — sentinels (`ErrNotFound`)

Implementation (under `internal/`, unimportable outside the module):

- `internal/scanner/` — reflect-based row scanner with per-type cache
- `internal/named/` — `:name` tokenizer + struct-tag value extractor (dialect-agnostic; takes a placeholder callback)
- `internal/null/` — nullable-pointer holders

## Code Style

- No external runtime deps — only `database/sql` + `reflect`
- DB drivers (sqlite, pgx, go-ora) are test-only
- Column lookups always lowercase the driver-reported name (Oracle-friendly)
- Empty result slice is `[]T{}`, never `nil`

## Important

- Public API stays at the module root; implementation helpers live under `internal/<subpkg>/`. New unexported helpers go in `internal/`, not the root — the `internal/` boundary is what enforces the public surface, so keep it small and intentional.
- NEVER add a runtime dependency outside the standard library
- NEVER type-assert to `*sql.DB` directly — go through `Querier` so transactions work
- NEVER swallow `rows.Err()` after the scan loop
- See `docs/ARCHITECTURE.md` for the scanner pipeline and dialect details
- See `docs/TESTING.md` for test layout, build tags, and integration setup
