package raizel_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rjansen/raizel"
)

// --- Fake Querier capturing the rewritten SQL + args ---

// recordingQuerier captures the last ExecContext call. Lets us assert
// the dialect-specific placeholder rewriting without needing a live
// PostgreSQL or Oracle to run unit tests against.
type recordingQuerier struct {
	lastQuery string
	lastArgs  []any
	calls     int
	failOn    int // 1-indexed call number to return an error on; 0 = never
}

type fakeResult struct{}

func (fakeResult) LastInsertId() (int64, error) { return 0, nil }
func (fakeResult) RowsAffected() (int64, error) { return 1, nil }

func (r *recordingQuerier) QueryContext(_ context.Context, _ string, _ ...any) (*sql.Rows, error) {
	return nil, errors.New("not implemented")
}

func (r *recordingQuerier) QueryRowContext(_ context.Context, _ string, _ ...any) *sql.Row {
	return nil
}

func (r *recordingQuerier) ExecContext(_ context.Context, query string, args ...any) (sql.Result, error) {
	r.calls++
	r.lastQuery = query
	r.lastArgs = args
	if r.failOn != 0 && r.calls == r.failOn {
		return nil, fmt.Errorf("synthetic failure on call %d", r.calls)
	}
	return fakeResult{}, nil
}

// --- Rewrite tests (per dialect) ---

type insertable struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	Score     int64      `db:"score"`
	DeletedAt *time.Time `db:"deleted_at"` // nil → SQL NULL
}

func TestExecNamed_SQLite_RewritesToQuestionMark(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	m := insertable{Name: "Alice", Score: 10}
	_, err := raizel.ExecNamed(ctx, q, raizel.DialectSQLite,
		"INSERT INTO t (name, score, deleted_at) VALUES (:name, :score, :deleted_at)", m)
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	want := "INSERT INTO t (name, score, deleted_at) VALUES (?, ?, ?)"
	if q.lastQuery != want {
		t.Errorf("query mismatch:\n got %q\nwant %q", q.lastQuery, want)
	}
	if len(q.lastArgs) != 3 || q.lastArgs[0] != "Alice" || q.lastArgs[1] != int64(10) || q.lastArgs[2] != nil {
		t.Errorf("args: %+v", q.lastArgs)
	}
}

func TestExecNamed_Postgres_RewritesToDollarN(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	m := insertable{Name: "Bob", Score: 99}
	_, err := raizel.ExecNamed(ctx, q, raizel.DialectPostgres,
		"INSERT INTO t (name, score) VALUES (:name, :score)", m)
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	want := "INSERT INTO t (name, score) VALUES ($1, $2)"
	if q.lastQuery != want {
		t.Errorf("query mismatch:\n got %q\nwant %q", q.lastQuery, want)
	}
}

func TestExecNamed_Oracle_RewritesToColonN(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	m := insertable{Name: "Carol", Score: 7}
	_, err := raizel.ExecNamed(ctx, q, raizel.DialectOracle,
		"INSERT INTO t (name, score) VALUES (:name, :score)", m)
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	want := "INSERT INTO t (name, score) VALUES (:1, :2)"
	if q.lastQuery != want {
		t.Errorf("query mismatch:\n got %q\nwant %q", q.lastQuery, want)
	}
}

func TestExecNamed_Postgres_DoubleColonCastPreserved(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	m := insertable{Name: "X", Score: 1}
	// `::int` is a Postgres type cast — must NOT be tokenized as :int.
	_, err := raizel.ExecNamed(ctx, q, raizel.DialectPostgres,
		"SELECT (:score)::int8 + 1 FROM t WHERE name = :name", m)
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	if !strings.Contains(q.lastQuery, "::int8") {
		t.Errorf("type cast lost: %q", q.lastQuery)
	}
	if !strings.Contains(q.lastQuery, "$1") || !strings.Contains(q.lastQuery, "$2") {
		t.Errorf("placeholders missing: %q", q.lastQuery)
	}
}

// --- Error paths ---

func TestExecNamed_NoModels_Errors(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	_, err := raizel.ExecNamed[insertable](ctx, q, raizel.DialectSQLite, "INSERT ...")
	if err == nil {
		t.Fatal("expected error for empty models")
	}
}

func TestExecNamed_UnknownTag_Errors(t *testing.T) {
	q := &recordingQuerier{}
	ctx := context.Background()

	m := insertable{Name: "X"}
	_, err := raizel.ExecNamed(ctx, q, raizel.DialectSQLite,
		"INSERT INTO t (name) VALUES (:nonexistent)", m)
	if err == nil {
		t.Fatal("expected error for unknown name")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("expected error to mention :nonexistent, got %q", err.Error())
	}
}

// --- Live SQLite for transaction behaviour ---

func TestExecNamed_BatchAutoTx_RollsBackOnFailure(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE)`); err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	rows := []row{
		{ID: 1, Name: "ok"},
		{ID: 2, Name: "ok"}, // second insert fails — same name
	}
	_, err := raizel.ExecNamed(ctx, db, raizel.DialectSQLite,
		"INSERT INTO t (id, name) VALUES (:id, :name)", rows[0], row{ID: 2, Name: "ok"})
	if err == nil {
		t.Fatal("expected unique-violation error")
	}

	// Auto-tx must have rolled back the first row too.
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM t").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("rollback didn't fire: count=%d", n)
	}
}

func TestExecNamed_BatchAutoTx_CommitsAllOnSuccess(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.Exec("CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	_, err := raizel.ExecNamed(ctx, db, raizel.DialectSQLite,
		"INSERT INTO t (id, name) VALUES (:id, :name)",
		row{ID: 1, Name: "a"}, row{ID: 2, Name: "b"}, row{ID: 3, Name: "c"})
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}

	count, err := raizel.QueryOne[int64](ctx, db, "SELECT COUNT(*) FROM t")
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Errorf("got %d, want 3", count)
	}
}

// When the caller already supplies a *sql.Tx, ExecNamed must reuse it
// rather than starting a nested one.
func TestExecNamed_RespectsCallerTx(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.Exec("CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if _, err := raizel.ExecNamed(ctx, tx, raizel.DialectSQLite,
		"INSERT INTO t (id, name) VALUES (:id, :name)",
		row{ID: 1, Name: "a"}, row{ID: 2, Name: "b"}); err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	// Roll the caller's tx back — both rows should disappear.
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}
	count, _ := raizel.QueryOne[int64](ctx, db, "SELECT COUNT(*) FROM t")
	if count != 0 {
		t.Errorf("rollback failed: count=%d", count)
	}
}

// --- Live SQLite end-to-end with named params ---

func TestExecNamed_LiveSQLite(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.Exec(`CREATE TABLE bills (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		amount REAL NOT NULL,
		paid_at DATETIME
	)`); err != nil {
		t.Fatal(err)
	}
	type bill struct {
		ID     int64      `db:"id"`
		Title  string     `db:"title"`
		Amount float64    `db:"amount"`
		PaidAt *time.Time `db:"paid_at"`
	}
	when := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	b := bill{ID: 1, Title: "rent", Amount: 1500, PaidAt: &when}
	if _, err := raizel.ExecNamed(ctx, db, raizel.DialectSQLite,
		"INSERT INTO bills (id, title, amount, paid_at) VALUES (:id, :title, :amount, :paid_at)", b); err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	got, err := raizel.QueryOne[bill](ctx, db, "SELECT id, title, amount, paid_at FROM bills WHERE id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "rent" || got.Amount != 1500 || got.PaidAt == nil || !got.PaidAt.Equal(when) {
		t.Errorf("round-trip mismatch: %+v", got)
	}
}
