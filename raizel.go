// Package raizel is a small, generic helper that maps database/sql query
// results into typed Go values without code generation. It supports
// SQLite, PostgreSQL, and Oracle through any standard database/sql driver.
//
// The public API is four generic functions over a dialect-bound [Handle]
// (a [DB] opened with [Open]/[Wrap], or a [Tx]):
//
//	Query[T]     run a SELECT and scan every row into a freshly-allocated T
//	QueryOne[T]  run a SELECT and scan the first row, or return ErrNotFound
//	Exec         run a non-query statement with positional arguments
//	ExecNamed[T] run a non-query whose `:name` placeholders are bound from
//	             struct fields tagged `db:"name"`
//
// Because the Handle carries the [Dialect], the dialect is configured once
// when the pool is opened rather than threaded through every call. Slice
// fields (e.g. []string) round-trip as native arrays on PostgreSQL and as
// JSON text on Oracle/SQLite, transparently.
package raizel

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	go_ora "github.com/sijms/go-ora/v2"

	"github.com/rjansen/raizel/internal/named"
	"github.com/rjansen/raizel/internal/scanner"
)

// oracleVarchar2BindLimit is the largest byte length go-ora binds as a
// VARCHAR2. A longer string overflows the bind with ORA-01461 even when the
// target column is a CLOB (the type is fixed client-side, before the server
// sees the statement), so such values must be sent as a LOB instead.
const oracleVarchar2BindLimit = 32767

// promoteOracleClobs rewrites oversized string arguments into go_ora.Clob so
// the driver binds them as CLOBs rather than VARCHAR2. Only strings beyond the
// VARCHAR2 bind ceiling are touched — shorter values bind as VARCHAR2 exactly
// as before, and any column that legitimately receives a longer value must be
// a CLOB/LONG, so the promotion is always type-correct. Oracle-only: other
// dialects have no such limit. It mutates args in place.
func promoteOracleClobs(args []any) {
	for i, a := range args {
		if s, ok := a.(string); ok && len(s) > oracleVarchar2BindLimit {
			args[i] = go_ora.Clob{String: s, Valid: true}
		}
	}
}

// Query runs the SQL and scans every row into a value of type T. T may be
// a tagged struct or a scalar (single-column) value. The returned slice is
// never nil — an empty result set yields []T{} of length zero.
func Query[T any](ctx context.Context, h Handle, query string, args ...any) ([]T, error) {
	rows, err := h.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("raizel.Query: %w", err)
	}
	defer rows.Close()

	rs, err := scanner.New[T](rows, !h.Dialect().NativeArrays())
	if err != nil {
		return nil, fmt.Errorf("raizel.Query: %w", err)
	}

	out := []T{}
	for rows.Next() {
		v, err := rs.Scan(rows)
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
func QueryOne[T any](ctx context.Context, h Handle, query string, args ...any) (T, error) {
	var zero T
	rows, err := h.QueryContext(ctx, query, args...)
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

	rs, err := scanner.New[T](rows, !h.Dialect().NativeArrays())
	if err != nil {
		return zero, fmt.Errorf("raizel.QueryOne: %w", err)
	}
	v, err := rs.Scan(rows)
	if err != nil {
		return zero, fmt.Errorf("raizel.QueryOne scan: %w", err)
	}
	return v, nil
}

// Exec runs a non-query statement with positional parameters. It is a
// thin wrapper around the handle's ExecContext that wraps driver errors
// with a raizel-prefixed context for easier root-cause analysis.
//
// Positional slice arguments are passed to the driver as-is; on dialects
// without native arrays, encode them yourself (or prefer ExecNamed, which
// handles slice encoding from struct tags).
func Exec(ctx context.Context, h Handle, query string, args ...any) (sql.Result, error) {
	if h.Dialect() == DialectOracle {
		promoteOracleClobs(args)
	}
	res, err := h.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("raizel.Exec: %w", err)
	}
	return res, nil
}

// ExecNamed runs query once for each model. The query may contain `:name`
// placeholders that are matched against the model's `db` tags and
// rewritten to the handle dialect's positional form (?, $N, :N) before
// being sent to the driver. Slice-typed fields are encoded as native
// arrays (PostgreSQL) or JSON text (Oracle/SQLite) to match the dialect.
//
// When called with two or more models against a *DB, ExecNamed wraps the
// batch in a single transaction so a mid-batch failure rolls back the rows
// that came before. When the caller already passes a *Tx, that transaction
// is reused as-is.
//
// Returns the sql.Result of the last successful execution.
func ExecNamed[T any](ctx context.Context, h Handle, query string, models ...T) (sql.Result, error) {
	if len(models) == 0 {
		return nil, errors.New("raizel.ExecNamed: at least one model is required")
	}

	if len(models) > 1 {
		if dbh, ok := h.(*DB); ok {
			tx, err := dbh.Begin(ctx)
			if err != nil {
				return nil, fmt.Errorf("raizel.ExecNamed begin: %w", err)
			}
			res, execErr := execNamedAll(ctx, tx, query, models)
			if execErr != nil {
				_ = tx.Rollback()
				return nil, execErr
			}
			if err := tx.Commit(); err != nil {
				return nil, fmt.Errorf("raizel.ExecNamed commit: %w", err)
			}
			return res, nil
		}
	}
	return execNamedAll(ctx, h, query, models)
}

func execNamedAll[T any](ctx context.Context, h Handle, query string, models []T) (sql.Result, error) {
	dialect := h.Dialect()
	jsonArrays := !dialect.NativeArrays()
	var last sql.Result
	for _, m := range models {
		rewritten, params, err := named.Rewrite(query, dialect.Placeholder, jsonArrays, m)
		if err != nil {
			return nil, fmt.Errorf("raizel.ExecNamed: %w", err)
		}
		if dialect == DialectOracle {
			promoteOracleClobs(params)
		}
		res, err := h.ExecContext(ctx, rewritten, params...)
		if err != nil {
			return nil, fmt.Errorf("raizel.ExecNamed: %w", err)
		}
		last = res
	}
	return last, nil
}
