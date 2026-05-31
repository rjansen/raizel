package raizel

import "strconv"

// Dialect identifies the SQL dialect a database driver speaks. It controls
// how named parameters are rewritten into the driver's positional
// placeholder form.
type Dialect int

const (
	// DialectSQLite uses '?' positional placeholders. Same shape as MySQL.
	DialectSQLite Dialect = iota
	// DialectPostgres uses '$1', '$2', ... positional placeholders.
	DialectPostgres
	// DialectOracle uses ':1', ':2', ... positional placeholders.
	DialectOracle
)

// Placeholder returns the driver-specific marker for the n-th parameter
// (1-indexed).
func (d Dialect) Placeholder(n int) string {
	switch d {
	case DialectPostgres:
		return "$" + strconv.Itoa(n)
	case DialectOracle:
		return ":" + strconv.Itoa(n)
	default:
		return "?"
	}
}

// String returns the canonical lowercase name of the dialect.
func (d Dialect) String() string {
	switch d {
	case DialectPostgres:
		return "postgres"
	case DialectOracle:
		return "oracle"
	default:
		return "sqlite"
	}
}

// NativeArrays reports whether the dialect's driver maps a Go slice
// (e.g. []string) directly onto a native array column. PostgreSQL does
// (TEXT[] via pgx); Oracle and SQLite do not, so raizel encodes slices as
// JSON text for those engines on the way in and out.
func (d Dialect) NativeArrays() bool { return d == DialectPostgres }

// driverName maps the dialect to the database/sql driver name raizel.Open
// passes to sql.Open. The matching driver must be registered (blank
// imported) by the calling binary:
//
//	postgres -> "pgx"    (github.com/jackc/pgx/v5/stdlib)
//	oracle   -> "oracle" (github.com/sijms/go-ora/v2)
//	sqlite   -> "sqlite" (modernc.org/sqlite)
func (d Dialect) driverName() string {
	switch d {
	case DialectPostgres:
		return "pgx"
	case DialectOracle:
		return "oracle"
	default:
		return "sqlite"
	}
}
