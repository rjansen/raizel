package raizel_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/rjansen/raizel"

	_ "modernc.org/sqlite"
)

// --- Nullable round-trip ---

type nullableUser struct {
	ID        int64      `db:"id"`
	Name      string     `db:"name"`
	DeletedAt *time.Time `db:"deleted_at"`
	Score     *int64     `db:"score"`
	Bonus     *float64   `db:"bonus"`
	Nickname  *string    `db:"nickname"`
	Active    *bool      `db:"active"`
}

func seedNullable(t *testing.T, db *raizel.DB) {
	t.Helper()
	if _, err := db.SQL().Exec(`CREATE TABLE n_users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		deleted_at DATETIME,
		score INTEGER,
		bonus REAL,
		nickname TEXT,
		active INTEGER
	)`); err != nil {
		t.Fatalf("create: %v", err)
	}
}

func TestQuery_AllNullableFieldsNull(t *testing.T) {
	db := openDB(t)
	seedNullable(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO n_users(name) VALUES ('Alice')"); err != nil {
		t.Fatal(err)
	}
	u, err := raizel.QueryOne[nullableUser](ctx, db,
		"SELECT id, name, deleted_at, score, bonus, nickname, active FROM n_users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Name != "Alice" {
		t.Errorf("name: %q", u.Name)
	}
	if u.DeletedAt != nil || u.Score != nil || u.Bonus != nil || u.Nickname != nil || u.Active != nil {
		t.Errorf("expected all-nil pointers, got %+v", u)
	}
}

func TestQuery_AllNullableFieldsSet(t *testing.T) {
	db := openDB(t)
	seedNullable(t, db)
	ctx := context.Background()

	deleted := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	if _, err := db.SQL().Exec(
		"INSERT INTO n_users(name, deleted_at, score, bonus, nickname, active) VALUES (?,?,?,?,?,?)",
		"Bob", deleted, int64(99), 1.5, "bb", true); err != nil {
		t.Fatal(err)
	}
	u, err := raizel.QueryOne[nullableUser](ctx, db,
		"SELECT id, name, deleted_at, score, bonus, nickname, active FROM n_users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.DeletedAt == nil || !u.DeletedAt.Equal(deleted) {
		t.Errorf("deleted_at: got %v", u.DeletedAt)
	}
	if u.Score == nil || *u.Score != 99 {
		t.Errorf("score: got %v", u.Score)
	}
	if u.Bonus == nil || *u.Bonus != 1.5 {
		t.Errorf("bonus: got %v", u.Bonus)
	}
	if u.Nickname == nil || *u.Nickname != "bb" {
		t.Errorf("nickname: got %v", u.Nickname)
	}
	if u.Active == nil || !*u.Active {
		t.Errorf("active: got %v", u.Active)
	}
}

func TestQuery_NullableMixed(t *testing.T) {
	db := openDB(t)
	seedNullable(t, db)
	ctx := context.Background()

	if _, err := db.SQL().Exec("INSERT INTO n_users(name, score) VALUES ('Eve', 7)"); err != nil {
		t.Fatal(err)
	}
	u, err := raizel.QueryOne[nullableUser](ctx, db,
		"SELECT id, name, deleted_at, score, bonus, nickname, active FROM n_users")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if u.Score == nil || *u.Score != 7 {
		t.Errorf("score: got %v", u.Score)
	}
	if u.DeletedAt != nil || u.Bonus != nil || u.Nickname != nil || u.Active != nil {
		t.Errorf("expected nil for unset columns, got %+v", u)
	}
}

// --- nullzero: NULL coerced to zero value on non-pointer fields ---

type nullzeroRow struct {
	ID      int64     `db:"id"`
	Body    string    `db:"body,nullzero"`
	Score   int64     `db:"score,nullzero"`
	Bonus   float64   `db:"bonus,nullzero"`
	Active  bool      `db:"active,nullzero"`
	SeenAt  time.Time `db:"seen_at,nullzero"`
	Subject string    `db:"subject"` // untagged: NULL must fail loudly
}

func seedNullzero(t *testing.T, db *raizel.DB) {
	t.Helper()
	if _, err := db.SQL().Exec(`CREATE TABLE nz (
		id INTEGER PRIMARY KEY,
		body TEXT,
		score INTEGER,
		bonus REAL,
		active INTEGER,
		seen_at DATETIME,
		subject TEXT
	)`); err != nil {
		t.Fatalf("create: %v", err)
	}
}

func TestNullzero_NullColumnsCoerceToZero(t *testing.T) {
	db := openDB(t)
	seedNullzero(t, db)
	ctx := context.Background()

	// Every nullzero column NULL; the untagged subject is non-NULL.
	if _, err := db.SQL().Exec(
		"INSERT INTO nz(id, body, score, bonus, active, seen_at, subject) VALUES (1, NULL, NULL, NULL, NULL, NULL, 'hi')"); err != nil {
		t.Fatal(err)
	}
	r, err := raizel.QueryOne[nullzeroRow](ctx, db,
		"SELECT id, body, score, bonus, active, seen_at, subject FROM nz")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if r.Body != "" || r.Score != 0 || r.Bonus != 0 || r.Active != false || !r.SeenAt.IsZero() {
		t.Errorf("expected zero values for NULL columns, got %+v", r)
	}
	if r.Subject != "hi" {
		t.Errorf("subject: got %q", r.Subject)
	}
}

func TestNullzero_PresentValuesPreserved(t *testing.T) {
	db := openDB(t)
	seedNullzero(t, db)
	ctx := context.Background()

	seen := time.Date(2026, 6, 6, 12, 0, 0, 0, time.UTC)
	if _, err := db.SQL().Exec(
		"INSERT INTO nz(id, body, score, bonus, active, seen_at, subject) VALUES (1, 'hello', 9, 2.5, 1, ?, 'sub')",
		seen); err != nil {
		t.Fatal(err)
	}
	r, err := raizel.QueryOne[nullzeroRow](ctx, db,
		"SELECT id, body, score, bonus, active, seen_at, subject FROM nz")
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if r.Body != "hello" || r.Score != 9 || r.Bonus != 2.5 || !r.Active || !r.SeenAt.Equal(seen) {
		t.Errorf("present values not preserved: %+v", r)
	}
}

func TestNullzero_UntaggedNullStillFailsLoudly(t *testing.T) {
	db := openDB(t)
	seedNullzero(t, db)
	ctx := context.Background()

	// subject is untagged and NULL → scanning into a plain string must error.
	if _, err := db.SQL().Exec(
		"INSERT INTO nz(id, body, subject) VALUES (1, 'b', NULL)"); err != nil {
		t.Fatal(err)
	}
	_, err := raizel.QueryOne[nullzeroRow](ctx, db,
		"SELECT id, body, score, bonus, active, seen_at, subject FROM nz")
	if err == nil {
		t.Fatal("expected error scanning NULL into untagged string field")
	}
}

// --- Nullable holder errors ---

type unsupportedNullable struct {
	ID    int64  `db:"id"`
	Count *int32 `db:"count"` // *int32 is not in the supported set
}

func TestNullable_UnsupportedTypeFailsLoudly(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()
	if _, err := db.SQL().Exec("CREATE TABLE x (id INTEGER, count INTEGER)"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec("INSERT INTO x VALUES (1, NULL)"); err != nil {
		t.Fatal(err)
	}
	_, err := raizel.QueryOne[unsupportedNullable](ctx, db, "SELECT id, count FROM x")
	if err == nil {
		t.Fatal("expected error for *int32 nullable")
	}
	if !strings.Contains(err.Error(), "unsupported nullable") {
		t.Errorf("expected 'unsupported nullable' in error, got %q", err.Error())
	}
}

// --- Nested struct via dot-notation aliases ---

type team struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type member struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Team team   `db:"team"` // children will be addressed as "team.id", "team.name"
}

func TestQuery_NestedStruct(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()
	if _, err := db.SQL().Exec(`
		CREATE TABLE teams (id INTEGER PRIMARY KEY, name TEXT);
		CREATE TABLE members (id INTEGER PRIMARY KEY, name TEXT, team_id INTEGER);
	`); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec("INSERT INTO teams VALUES (1, 'Alpha')"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec("INSERT INTO members VALUES (10, 'Alice', 1)"); err != nil {
		t.Fatal(err)
	}
	m, err := raizel.QueryOne[member](ctx, db,
		`SELECT m.id, m.name, t.id "team.id", t.name "team.name"
		   FROM members m JOIN teams t ON m.team_id = t.id`)
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if m.ID != 10 || m.Name != "Alice" {
		t.Errorf("member: %+v", m)
	}
	if m.Team.ID != 1 || m.Team.Name != "Alpha" {
		t.Errorf("team: %+v", m.Team)
	}
}

// Nested struct child has a nullable field — the dbm bug we fixed
// silently wiped its nullability when copying. This test would catch a
// regression of that.
type outerWithNullableInner struct {
	ID    int64       `db:"id"`
	Inner innerHasOpt `db:"inner"`
}

type innerHasOpt struct {
	Tag       string     `db:"tag"`
	StoppedAt *time.Time `db:"stopped_at"` // nullable leaf inside nested struct
}

func TestQuery_NestedNullableLeafPreserved(t *testing.T) {
	db := openDB(t)
	ctx := context.Background()
	if _, err := db.SQL().Exec("CREATE TABLE outer_t (id INTEGER, tag TEXT, stopped_at DATETIME)"); err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().Exec("INSERT INTO outer_t VALUES (1, 'live', NULL)"); err != nil {
		t.Fatal(err)
	}
	o, err := raizel.QueryOne[outerWithNullableInner](ctx, db,
		`SELECT id, tag "inner.tag", stopped_at "inner.stopped_at" FROM outer_t`)
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if o.Inner.Tag != "live" {
		t.Errorf("tag: %q", o.Inner.Tag)
	}
	if o.Inner.StoppedAt != nil {
		t.Errorf("expected nil StoppedAt, got %v", o.Inner.StoppedAt)
	}

	stop := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if _, err := db.SQL().Exec("INSERT INTO outer_t VALUES (2, 'done', ?)", stop); err != nil {
		t.Fatal(err)
	}
	o2, err := raizel.QueryOne[outerWithNullableInner](ctx, db,
		`SELECT id, tag "inner.tag", stopped_at "inner.stopped_at" FROM outer_t WHERE id = 2`)
	if err != nil {
		t.Fatalf("QueryOne: %v", err)
	}
	if o2.Inner.StoppedAt == nil || !o2.Inner.StoppedAt.Equal(stop) {
		t.Errorf("StoppedAt: %v", o2.Inner.StoppedAt)
	}
}

// --- RETURNING + partial model ---

func TestQueryOne_ReturningClause(t *testing.T) {
	db := openDB(t)
	seedNullable(t, db)
	ctx := context.Background()

	u, err := raizel.QueryOne[nullableUser](ctx, db,
		`INSERT INTO n_users (name, score) VALUES (?, ?)
		 RETURNING id, name, deleted_at, score, bonus, nickname, active`,
		"Diana", 77)
	if err != nil {
		t.Fatalf("QueryOne RETURNING: %v", err)
	}
	if u.ID <= 0 || u.Name != "Diana" {
		t.Errorf("got %+v", u)
	}
	if u.Score == nil || *u.Score != 77 {
		t.Errorf("score: %v", u.Score)
	}
}
