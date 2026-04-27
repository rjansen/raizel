# Architecture

raizel is one Go package at the module root. The four public functions —
`Query[T]`, `QueryOne[T]`, `Exec`, `ExecNamed[T]` — sit on top of three
internal subsystems: a row scanner, a nullable-pointer shuttle, and a
named-parameter rewriter.

## File map

| File | Role |
|------|------|
| `raizel.go` | Public API: `Query[T]`, `QueryOne[T]`, `Exec`, `ExecNamed[T]`, the `txBeginner` detection used to auto-wrap batches. |
| `querier.go` | `Querier` interface — the minimal subset of `*sql.DB` (also satisfied by `*sql.Tx`). |
| `dialect.go` | `Dialect` enum + `Placeholder(n)` for `?` / `$N` / `:N`. |
| `errors.go` | `ErrNotFound` returned by `QueryOne` when no row matches. |
| `scanner.go` | `fieldIndex`, the per-type field-map cache (`sync.Map`), `rowScanner[T]`, `fieldByPath`. |
| `null.go` | `nullHolder` interface and shuttles for `*time.Time`, `*int64`, `*float64`, `*string`, `*bool`. |
| `named.go` | `rewriteNamed` tokenizer (`:name` → dialect placeholder; `::cast` is preserved) and `tagValues` extractor. |

## Scanner pipeline

1. `Query` / `QueryOne` calls `q.QueryContext` and gets `*sql.Rows`.
2. `newRowScanner[T]` is built once: it asks the rows for their column
   names, walks `T`'s `db` tags via `fieldsByColumn` (cached per type),
   and produces a `[]fieldIndex` aligned with the column order.
3. Each `fieldIndex` carries (a) the index path through possibly-nested
   structs, (b) a `nullable` flag, (c) the leaf base type used to pick a
   `nullHolder`. An empty path marks an unmapped column — its slot
   discards into a fresh `*any` sink so the driver still sees the right
   number of destinations.
4. `scan` allocates one destination per column. Non-nullable fields scan
   straight into the addressable struct field; nullable fields scan into
   a `nullHolder`'s `Null*` shuttle and a second pass calls
   `holder.assign(field)` to set the pointer when the column was non-NULL.
5. `fieldByPath` walks possibly-nested struct values and auto-allocates
   any nil pointer-to-struct it crosses, so a JOIN can scan straight into
   `Member.Team.ID` even when `Team` was zero on entry.

## Nullable shuttling

`nullHolderFor` returns a fresh holder per scan; pre-flight validation in
`newRowScanner` rejects unsupported pointer base types (e.g. `*int32`)
before any driver memory is touched. Adding a new supported type means
adding a new `nullHolder` implementation and a switch arm in
`nullHolderFor` — the cache and scanner do not need to change.

## Dialect-aware named params

`ExecNamed[T]` walks the SQL once with a small state machine:

- `:name` (where `name` is `[A-Za-z0-9_]+`) is replaced with
  `dialect.Placeholder(n)`; the matching value is read from the model's
  `db:"name"` tag.
- `::` is recognised as a Postgres type cast and copied through.
- Lone `:` followed by a non-name char is left as-is.

When two or more models are passed against a `*sql.DB`, the function
opens a transaction via `txBeginner.BeginTx`, runs every model on the
returned `*sql.Tx`, and commits on success or rolls back on the first
error. When the caller already supplies a `*sql.Tx`, it is reused
verbatim — no nested transactions.

## Caches and concurrency

- `fieldCache` (`sync.Map`) is keyed by `reflect.Type`. Reads dominate
  writes after warmup; the map is goroutine-safe.
- Each `rowScanner` is built per-query and therefore not shared across
  goroutines. The cache it consults is.
- No package-global state mutates after init.
