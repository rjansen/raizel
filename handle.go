package raizel

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// DB is a database/sql pool bound to a [Dialect]. Create one once at
// startup with [Open] (raizel owns the pool) or [Wrap] (you own it), then
// pass it to the generic helpers — the dialect rides along, so call sites
// stay free of per-query dialect noise.
type DB struct {
	db      *sql.DB
	dialect Dialect
}

// Open maps the dialect to its database/sql driver, opens a connection
// pool, and returns a dialect-bound handle. The matching driver must be
// registered (blank imported) by the calling binary — see
// [Dialect.driverName]. Pool tuning is applied from opts.
//
// Open does not Ping; call [DB.PingContext] if you need to verify
// connectivity eagerly.
func Open(dialect Dialect, dsn string, opts ...Option) (*DB, error) {
	sqlDB, err := sql.Open(dialect.driverName(), dsn)
	if err != nil {
		return nil, fmt.Errorf("raizel.Open %s: %w", dialect, err)
	}
	for _, opt := range opts {
		opt(sqlDB)
	}
	return &DB{db: sqlDB, dialect: dialect}, nil
}

// Wrap adopts an existing *sql.DB whose lifecycle and pool settings the
// caller manages. Use it when you need driver-specific setup (custom
// session params, secrets injection) that [Open] does not cover.
func Wrap(db *sql.DB, dialect Dialect) *DB {
	return &DB{db: db, dialect: dialect}
}

// SQL returns the underlying *sql.DB for pool tuning, health checks, or any
// operation outside raizel's helpers.
func (d *DB) SQL() *sql.DB { return d.db }

// Dialect reports the SQL dialect this handle speaks.
func (d *DB) Dialect() Dialect { return d.dialect }

// P renders the dialect's positional placeholder for the n-th parameter
// (1-indexed) — a shorthand for d.Dialect().Placeholder(n) when composing
// SQL by hand.
func (d *DB) P(n int) string { return d.dialect.Placeholder(n) }

// Close closes the underlying pool.
func (d *DB) Close() error { return d.db.Close() }

// PingContext verifies a connection can be established.
func (d *DB) PingContext(ctx context.Context) error { return d.db.PingContext(ctx) }

// QueryContext implements [Querier].
func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.db.QueryContext(ctx, query, args...)
}

// QueryRowContext implements [Querier].
func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return d.db.QueryRowContext(ctx, query, args...)
}

// ExecContext implements [Querier].
func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

// Begin starts a transaction that carries this handle's dialect forward.
func (d *DB) Begin(ctx context.Context) (*Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("raizel.Begin: %w", err)
	}
	return &Tx{tx: tx, dialect: d.dialect}, nil
}

// Tx is a dialect-bound transaction. It satisfies [Handle], so the generic
// helpers work against it identically to a [DB].
type Tx struct {
	tx      *sql.Tx
	dialect Dialect
}

// Dialect reports the SQL dialect this transaction speaks.
func (t *Tx) Dialect() Dialect { return t.dialect }

// SQL returns the underlying *sql.Tx.
func (t *Tx) SQL() *sql.Tx { return t.tx }

// Commit commits the transaction.
func (t *Tx) Commit() error { return t.tx.Commit() }

// Rollback aborts the transaction.
func (t *Tx) Rollback() error { return t.tx.Rollback() }

// QueryContext implements [Querier].
func (t *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return t.tx.QueryContext(ctx, query, args...)
}

// QueryRowContext implements [Querier].
func (t *Tx) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return t.tx.QueryRowContext(ctx, query, args...)
}

// ExecContext implements [Querier].
func (t *Tx) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return t.tx.ExecContext(ctx, query, args...)
}

// Option configures the *sql.DB pool opened by [Open].
type Option func(*sql.DB)

// WithMaxOpenConns caps the number of open connections.
func WithMaxOpenConns(n int) Option { return func(db *sql.DB) { db.SetMaxOpenConns(n) } }

// WithMaxIdleConns caps the idle connection pool.
func WithMaxIdleConns(n int) Option { return func(db *sql.DB) { db.SetMaxIdleConns(n) } }

// WithConnMaxLifetime bounds how long a connection may be reused.
func WithConnMaxLifetime(d time.Duration) Option {
	return func(db *sql.DB) { db.SetConnMaxLifetime(d) }
}

// WithConnMaxIdleTime bounds how long a connection may sit idle.
func WithConnMaxIdleTime(d time.Duration) Option {
	return func(db *sql.DB) { db.SetConnMaxIdleTime(d) }
}
