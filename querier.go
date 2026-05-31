package raizel

import (
	"context"
	"database/sql"
)

// Querier is the minimal subset of *sql.DB / *sql.Tx that raizel's helpers
// invoke. It is satisfied by the standard library types as well as by
// raizel's own [DB] and [Tx] handles.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// Handle is a [Querier] that also knows which SQL [Dialect] it speaks. The
// generic helpers (Query, QueryOne, Exec, ExecNamed) take a Handle so the
// dialect is configured once — when the pool is opened — rather than passed
// on every call. Both *[DB] and *[Tx] satisfy it, so the same helper code
// works in or out of a transaction.
type Handle interface {
	Querier
	Dialect() Dialect
}
