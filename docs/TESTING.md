# Testing

raizel ships two test layers:

- **Unit tests** run against in-memory SQLite (`modernc.org/sqlite`,
  pure Go). They are fast, self-contained, and require no Docker.
  `go test ./...` and `make test` both run this layer.
- **Integration tests** are gated behind the `integration` build tag and
  exercise the same scenarios against PostgreSQL and Oracle running in
  docker-compose.

## Layout

| File | Layer | Purpose |
|------|-------|---------|
| `raizel_test.go` | unit | Core API: structs, scalars, transactions, tag handling, case-insensitive columns, error wrapping. |
| `null_test.go` | unit | Nullable round-trip, unsupported-type errors, nested-struct scanning, RETURNING clauses. |
| `named_test.go` | unit | Named-param rewriting per dialect (uses a `recordingQuerier` fake to assert query strings without a live DB), batch transaction semantics, end-to-end SQLite. |
| `integration_test.go` | integration | Live PostgreSQL + Oracle exercising CRUD and batch-rollback. |

## Running unit tests

```sh
make test            # short mode, race detector off
go test ./...        # equivalent
go test -run X ./... # single test
```

## Running integration tests

```sh
make integration     # boots docker-compose, runs, tears down
```

Or manually:

```sh
docker compose up -d postgres oracle
# Wait for both to be healthy (oracle takes ~60 s on first boot)
POSTGRES_DSN='postgres://raizel:raizel@localhost:5432/raizel?sslmode=disable' \
ORACLE_DSN='oracle://raizel:raizel@localhost:1521/FREEPDB1' \
go test -tags=integration -v ./...
docker compose down
```

A dialect whose DSN env var is unset is **skipped** (not failed), so a
developer can iterate on PostgreSQL alone without booting Oracle.

Integration test tables are prefixed `raizel_test_` to avoid colliding
with whatever schema the dev DB might already host. Each test pre-drops
its table before creating it, and registers a `t.Cleanup` to drop it on
exit.

## CI

| Workflow | Trigger | What it does |
|----------|---------|--------------|
| `.github/workflows/check.yml` | every push + PR | gofmt -l, go vet (default + integration tag), golangci-lint. |
| `.github/workflows/test.yml` | every push + PR | `go test -short -race -coverprofile=…` and prints the coverage summary. |
| `.github/workflows/integration.yml` | manual (`workflow_dispatch`) | Boots postgres + oracle as service containers and runs `go test -tags=integration -v ./...`. |

## Adding new tests

- For a new helper or scanner behaviour: extend `raizel_test.go` /
  `null_test.go` / `named_test.go`. Keep them SQLite-only — the unit
  layer must not require Docker.
- For a new dialect-specific behaviour: add a case to
  `integration_test.go`. Use `c.dialect.Placeholder(n)` to build
  positional placeholders portably; never hardcode `?`/`$1`/`:1`.
- The `recordingQuerier` fake in `named_test.go` is the right tool when
  you want to assert the exact SQL string raizel sends to the driver,
  for any dialect, without a live DB.
