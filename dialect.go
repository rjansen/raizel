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
