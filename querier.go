package raizel

import (
	"context"
	"database/sql"
)

// Querier is the minimal subset of *sql.DB that raizel needs. *sql.DB and
// *sql.Tx both satisfy it, so the same helper code works in or out of a
// transaction.
type Querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}
