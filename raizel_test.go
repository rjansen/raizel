package raizel_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/rjansen/raizel"

	_ "modernc.org/sqlite"
)

// --- helpers ---

func openDB(t *testing.T) *raizel.DB {
	t.Helper()
	sqlDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	return raizel.Wrap(sqlDB, raizel.DialectSQLite)
}

func seedUsers(t *testing.T, db *raizel.DB) {
	t.Helper()
	_, err := db.SQL().Exec(`CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		score INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
}

type user struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
	Score int64  `db:"score"`
}

// --- Query tests ---

func TestQuery_Struct(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES (?,?,?)", "Alice", "a@b.com", 10); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES (?,?,?)", "Bob", "b@b.com", 20); err != nil {
		t.Fatal(err)
	}

	users, err := raizel.Query[user](ctx, db, "SELECT id, name, email, score FROM users ORDER BY id")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("got %d users, want 2", len(users))
	}
	if users[0].Name != "Alice" || users[1].Score != 20 {
		t.Errorf("unexpected: %+v", users)
	}
}

func TestQuery_ScalarString(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email) VALUES ('A','a@b.com'),('B','b@b.com')"); err != nil {
		t.Fatal(err)
	}
	names, err := raizel.Query[string](ctx, db, "SELECT name FROM users ORDER BY name")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(names) != 2 || names[0] != "A" || names[1] != "B" {
		t.Errorf("got %v", names)
	}
}

func TestQuery_ScalarInt(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES ('A','a',10),('B','b',20)"); err != nil {
		t.Fatal(err)
	}
	scores, err := raizel.Query[int64](ctx, db, "SELECT score FROM users ORDER BY score")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if len(scores) != 2 || scores[0] != 10 || scores[1] != 20 {
		t.Errorf("got %v", scores)
	}
}

func TestQuery_Empty_ReturnsNonNilSlice(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	users, err := raizel.Query[user](ctx, db, "SELECT id, name, email, score FROM users")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	if users == nil {
		t.Fatal("expected non-nil empty slice")
	}
	if len(users) != 0 {
		t.Errorf("got len=%d, want 0", len(users))
	}
}

// --- QueryOne tests ---

func TestQueryOne_Struct(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES ('Alice','a@b.com',99)"); err != nil {
		t.Fatal(err)
	}
	u, err := raizel.QueryOne[user](ctx, db,
		"SELECT id, name, email, score FROM users WHERE name = ?", "Alice")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Name != "Alice" || u.Score != 99 {
		t.Errorf("got %+v", u)
	}
}

func TestQueryOne_NotFound(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	_, err := raizel.QueryOne[user](ctx, db,
		"SELECT id, name, email, score FROM users WHERE id = ?", 999)
	if !errors.Is(err, raizel.ErrNotFound) {
		t.Fatalf("got err=%v, want ErrNotFound", err)
	}
}

func TestQueryOne_Scalar(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email) VALUES ('A','a'),('B','b')"); err != nil {
		t.Fatal(err)
	}
	count, err := raizel.QueryOne[int64](ctx, db, "SELECT COUNT(*) FROM users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if count != 2 {
		t.Errorf("got %d, want 2", count)
	}
}

// --- Exec tests ---

func TestExec_Insert(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	res, err := raizel.Exec(ctx, db,
		"INSERT INTO users(name,email,score) VALUES (?,?,?)", "Alice", "a@b.com", 10)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	n, err := res.RowsAffected()
	if err != nil || n != 1 {
		t.Errorf("rows affected: got %d, err=%v", n, err)
	}
}

// --- Querier compatibility (transactions) ---

func TestQuerier_Tx(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	tx, err := db.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := raizel.Exec(ctx, tx,
		"INSERT INTO users(name,email) VALUES (?,?)", "Tx", "t@x.com"); err != nil {
		t.Fatalf("Exec via tx: %v", err)
	}
	users, err := raizel.Query[user](ctx, tx, "SELECT id, name, email, score FROM users")
	if err != nil {
		t.Fatalf("Query via tx: %v", err)
	}
	if len(users) != 1 || users[0].Name != "Tx" {
		t.Errorf("tx visibility broken: %+v", users)
	}
	if err := tx.Rollback(); err != nil {
		t.Fatal(err)
	}

	// After rollback, the row is gone.
	count, _ := raizel.QueryOne[int64](ctx, db, "SELECT COUNT(*) FROM users")
	if count != 0 {
		t.Errorf("rollback failed: count=%d", count)
	}
}

// --- Tag handling ---

type taggedUser struct {
	ID       int64  `db:"id,omitempty"` // tag option must be ignored
	FullName string `db:"name"`
	Internal string `db:"-"` // skipped
	Untagged string // no tag — skipped
	Email    string `db:"email"`
}

func TestStruct_TagOptions(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES ('Tagged','t@e.com',1)"); err != nil {
		t.Fatal(err)
	}
	u, err := raizel.QueryOne[taggedUser](ctx, db, "SELECT id, name, email FROM users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.FullName != "Tagged" || u.Email != "t@e.com" || u.ID == 0 {
		t.Errorf("got %+v", u)
	}
	if u.Internal != "" || u.Untagged != "" {
		t.Errorf("expected skipped fields zero, got %+v", u)
	}
}

func TestStruct_CaseInsensitiveColumns(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email) VALUES ('A','a@b.com')"); err != nil {
		t.Fatal(err)
	}
	// Force uppercase aliases so the column-name lookup must lowercase.
	u, err := raizel.QueryOne[user](ctx, db,
		`SELECT id AS "ID", name AS "NAME", email AS "EMAIL", score AS "SCORE" FROM users`)
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Name != "A" || u.Email != "a@b.com" {
		t.Errorf("got %+v", u)
	}
}

func TestStruct_PartialColumns(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES ('A','a@b.com',42)"); err != nil {
		t.Fatal(err)
	}
	// Only id + name selected — other fields should remain zero-valued.
	u, err := raizel.QueryOne[user](ctx, db, "SELECT id, name FROM users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Name != "A" {
		t.Errorf("name: got %q", u.Name)
	}
	if u.Email != "" || u.Score != 0 {
		t.Errorf("expected zero email/score, got %+v", u)
	}
}

func TestStruct_UnmappedColumnsDiscarded(t *testing.T) {
	db := openDB(t)
	seedUsers(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO users(name,email,score) VALUES ('A','a@b.com',1)"); err != nil {
		t.Fatal(err)
	}
	// Add an extra column that has no struct field — must not error.
	u, err := raizel.QueryOne[user](ctx, db,
		"SELECT id, name, email, score, 'extra' AS not_in_struct FROM users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Name != "A" {
		t.Errorf("got %+v", u)
	}
}

// --- Error wrapping ---

func TestQuery_DriverErrorWrapped(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()

	_, err := raizel.Query[user](ctx, db, "SELECT * FROM does_not_exist")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.HasPrefix(err.Error(), "raizel.Query") {
		t.Errorf("expected raizel-prefixed error, got %q", err.Error())
	}
}
