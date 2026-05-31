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
// PostgreSQL or Oracle to run unit tests against. It satisfies
// raizel.Handle by reporting a configurable dialect.
type recordingQuerier struct {
	dialect   raizel.Dialect
	lastQuery string
	lastArgs  []any
	calls     int
	failOn    int // 1-indexed call number to return an error on; 0 = never
}

func (r *recordingQuerier) Dialect() raizel.Dialect { return r.dialect }

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
	q := &recordingQuerier{dialect: raizel.DialectSQLite}
	ctx := context.Background()

	m := insertable{Name: "Alice", Score: 10}
	_, err := raizel.ExecNamed(ctx, q,
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
	q := &recordingQuerier{dialect: raizel.DialectPostgres}
	ctx := context.Background()

	m := insertable{Name: "Bob", Score: 99}
	_, err := raizel.ExecNamed(ctx, q,
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
	q := &recordingQuerier{dialect: raizel.DialectOracle}
	ctx := context.Background()

	m := insertable{Name: "Carol", Score: 7}
	_, err := raizel.ExecNamed(ctx, q,
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
	q := &recordingQuerier{dialect: raizel.DialectPostgres}
	ctx := context.Background()

	m := insertable{Name: "X", Score: 1}
	// `::int` is a Postgres type cast — must NOT be tokenized as :int.
	_, err := raizel.ExecNamed(ctx, q,
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

// --- Slice fields encode per dialect ---

type taggable struct {
	ID     int64    `db:"id"`
	Labels []string `db:"labels"`
}

func TestExecNamed_SliceField_JSONEncodedForOracle(t *testing.T) {
	q := &recordingQuerier{dialect: raizel.DialectOracle}
	ctx := context.Background()

	_, err := raizel.ExecNamed(ctx, q,
		"INSERT INTO t (id, labels) VALUES (:id, :labels)",
		taggable{ID: 1, Labels: []string{"INBOX", "IMPORTANT"}})
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	if len(q.lastArgs) != 2 {
		t.Fatalf("args: %+v", q.lastArgs)
	}
	if got, ok := q.lastArgs[1].(string); !ok || got != `["INBOX","IMPORTANT"]` {
		t.Errorf("labels arg: got %#v, want JSON string", q.lastArgs[1])
	}
}

func TestExecNamed_NilSliceField_EncodesEmptyJSONArrayForOracle(t *testing.T) {
	q := &recordingQuerier{dialect: raizel.DialectOracle}
	ctx := context.Background()

	_, err := raizel.ExecNamed(ctx, q,
		"INSERT INTO t (id, labels) VALUES (:id, :labels)", taggable{ID: 1})
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	if got, ok := q.lastArgs[1].(string); !ok || got != "[]" {
		t.Errorf("nil labels arg: got %#v, want \"[]\"", q.lastArgs[1])
	}
}

func TestExecNamed_SliceField_PassthroughForPostgres(t *testing.T) {
	q := &recordingQuerier{dialect: raizel.DialectPostgres}
	ctx := context.Background()

	_, err := raizel.ExecNamed(ctx, q,
		"INSERT INTO t (id, labels) VALUES (:id, :labels)",
		taggable{ID: 1, Labels: []string{"INBOX"}})
	if err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}
	got, ok := q.lastArgs[1].([]string)
	if !ok || len(got) != 1 || got[0] != "INBOX" {
		t.Errorf("labels arg: got %#v, want []string{\"INBOX\"} passthrough", q.lastArgs[1])
	}
}

// --- Error paths ---

func TestExecNamed_NoModels_Errors(t *testing.T) {
	q := &recordingQuerier{dialect: raizel.DialectSQLite}
	ctx := context.Background()

	_, err := raizel.ExecNamed[insertable](ctx, q, "INSERT ...")
	if err == nil {
		t.Fatal("expected error for empty models")
	}
}

func TestExecNamed_UnknownTag_Errors(t *testing.T) {
	q := &recordingQuerier{dialect: raizel.DialectSQLite}
	ctx := context.Background()

	m := insertable{Name: "X"}
	_, err := raizel.ExecNamed(ctx, q,
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

	if _, err := db.SQL().Exec(`CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE)`); err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	_, err := raizel.ExecNamed(ctx, db,
		"INSERT INTO t (id, name) VALUES (:id, :name)",
		row{ID: 1, Name: "ok"}, row{ID: 2, Name: "ok"}) // second insert fails — same name
	if err == nil {
		t.Fatal("expected unique-violation error")
	}

	// Auto-tx must have rolled back the first row too.
	var n int
	if err := db.SQL().QueryRow("SELECT COUNT(*) FROM t").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Errorf("rollback didn't fire: count=%d", n)
	}
}

func TestExecNamed_BatchAutoTx_CommitsAllOnSuccess(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.SQL().Exec("CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	_, err := raizel.ExecNamed(ctx, db,
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

// When the caller already supplies a *Tx, ExecNamed must reuse it rather
// than starting a nested one.
func TestExecNamed_RespectsCallerTx(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.SQL().Exec("CREATE TABLE t (id INTEGER PRIMARY KEY, name TEXT)"); err != nil {
		t.Fatal(err)
	}
	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	type row struct {
		ID   int64  `db:"id"`
		Name string `db:"name"`
	}
	if _, err := raizel.ExecNamed(ctx, tx,
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

	if _, err := db.SQL().Exec(`CREATE TABLE bills (
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
	if _, err := raizel.ExecNamed(ctx, db,
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

// --- Live SQLite slice round-trip (JSON-encoded) ---

func TestExecNamed_LiveSQLite_SliceRoundTrip(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	if _, err := db.SQL().Exec(`CREATE TABLE docs (id INTEGER PRIMARY KEY, labels TEXT)`); err != nil {
		t.Fatal(err)
	}
	if _, err := raizel.ExecNamed(ctx, db,
		"INSERT INTO docs (id, labels) VALUES (:id, :labels)",
		taggable{ID: 1, Labels: []string{"INBOX", "IMPORTANT"}},
		taggable{ID: 2, Labels: nil}); err != nil {
		t.Fatalf("ExecNamed: %v", err)
	}

	got, err := raizel.QueryOne[taggable](ctx, db,
		"SELECT id, labels FROM docs WHERE id = ?", 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Labels) != 2 || got.Labels[0] != "INBOX" || got.Labels[1] != "IMPORTANT" {
		t.Errorf("labels round-trip: got %#v", got.Labels)
	}

	empty, err := raizel.QueryOne[taggable](ctx, db,
		"SELECT id, labels FROM docs WHERE id = ?", 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(empty.Labels) != 0 {
		t.Errorf("nil labels round-trip: got %#v, want empty", empty.Labels)
	}
}
