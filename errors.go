package raizel

import "errors"

// ErrNotFound is returned by QueryOne when the query produces no rows.
//
// We use a custom sentinel rather than re-exporting sql.ErrNoRows because
// raizel's QueryOne does not call (*sql.Row).Scan — it iterates a
// (*sql.Rows) so the errors.Is chain to sql.ErrNoRows would not fire
// naturally.
var ErrNotFound = errors.New("raizel: not found")
