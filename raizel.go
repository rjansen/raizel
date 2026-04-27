// Package raizel is a small, generic helper that maps database/sql query
// results into typed Go values without code generation. It supports
// SQLite, PostgreSQL, and Oracle through any standard database/sql driver.
//
// The public API is four functions:
//
//	Query[T]     run a SELECT and scan every row into a freshly-allocated T
//	QueryOne[T]  run a SELECT and scan the first row, or return ErrNotFound
//	Exec         run a non-query statement with positional arguments
//	ExecNamed[T] run a non-query whose `:name` placeholders are bound from
//	             struct fields tagged `db:"name"` (added in Round 4)
//
// Both `*sql.DB` and `*sql.Tx` satisfy the [Querier] interface, so the
// same helpers work in or out of a transaction.
package raizel

import (
	"context"
	"database/sql"
	"fmt"
)

// Query runs the SQL and scans every row into a value of type T. T may be
// a tagged struct or a scalar (single-column) value. The returned slice is
// never nil — an empty result set yields []T{} of length zero.
func Query[T any](ctx context.Context, q Querier, query string, args ...any) ([]T, error) {
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("raizel.Query: %w", err)
	}
	defer rows.Close()

	scanner, err := newRowScanner[T](rows)
	if err != nil {
		return nil, fmt.Errorf("raizel.Query: %w", err)
	}

	out := []T{}
	for rows.Next() {
		v, err := scanner.scan(rows)
		if err != nil {
			return nil, fmt.Errorf("raizel.Query scan: %w", err)
		}
		out = append(out, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("raizel.Query rows: %w", err)
	}
	return out, nil
}

// QueryOne runs the SQL and scans the first row into a freshly-allocated
// T. It returns ErrNotFound when the query produces no rows.
func QueryOne[T any](ctx context.Context, q Querier, query string, args ...any) (T, error) {
	var zero T
	rows, err := q.QueryContext(ctx, query, args...)
	if err != nil {
		return zero, fmt.Errorf("raizel.QueryOne: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return zero, fmt.Errorf("raizel.QueryOne rows: %w", err)
		}
		return zero, ErrNotFound
	}

	scanner, err := newRowScanner[T](rows)
	if err != nil {
		return zero, fmt.Errorf("raizel.QueryOne: %w", err)
	}
	v, err := scanner.scan(rows)
	if err != nil {
		return zero, fmt.Errorf("raizel.QueryOne scan: %w", err)
	}
	return v, nil
}

// Exec runs a non-query statement with positional parameters. It is a
// thin wrapper around q.ExecContext that wraps driver errors with a
// raizel-prefixed context for easier root-cause analysis.
func Exec(ctx context.Context, q Querier, query string, args ...any) (sql.Result, error) {
	res, err := q.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("raizel.Exec: %w", err)
	}
	return res, nil
}
