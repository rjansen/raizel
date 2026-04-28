// Package null shuttles single columns between database/sql's Null{Time,
// Int64,Float64,String,Bool} scan targets and pointer-typed struct fields
// (*time.Time, *int64, ...). It is internal plumbing used by the row
// scanner; the public raizel API does not expose it directly.
package null

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

var timeType = reflect.TypeFor[time.Time]()

// Holder shuttles a single column between (a) database/sql's
// Null{Time,Int64,Float64,String,Bool} types and (b) a nullable struct
// field of pointer kind (*time.Time, *int64, ...). When the column scans
// NULL the holder leaves the field's zero pointer in place; otherwise it
// allocates a new value and points the field at it.
type Holder interface {
	ScanDest() any
	Assign(field reflect.Value)
}

type nullTimeHolder struct{ v sql.NullTime }

func (h *nullTimeHolder) ScanDest() any { return &h.v }
func (h *nullTimeHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		t := h.v.Time
		f.Set(reflect.ValueOf(&t))
	}
}

type nullInt64Holder struct{ v sql.NullInt64 }

func (h *nullInt64Holder) ScanDest() any { return &h.v }
func (h *nullInt64Holder) Assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Int64
		f.Set(reflect.ValueOf(&x))
	}
}

type nullFloat64Holder struct{ v sql.NullFloat64 }

func (h *nullFloat64Holder) ScanDest() any { return &h.v }
func (h *nullFloat64Holder) Assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Float64
		f.Set(reflect.ValueOf(&x))
	}
}

type nullStringHolder struct{ v sql.NullString }

func (h *nullStringHolder) ScanDest() any { return &h.v }
func (h *nullStringHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.String
		f.Set(reflect.ValueOf(&x))
	}
}

type nullBoolHolder struct{ v sql.NullBool }

func (h *nullBoolHolder) ScanDest() any { return &h.v }
func (h *nullBoolHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Bool
		f.Set(reflect.ValueOf(&x))
	}
}

// HolderFor returns a fresh holder for the given pointer base type.
// Unknown types fail loudly rather than scanning into the wrong shuttle —
// silent fallbacks were one of the bugs raizel exists to avoid.
func HolderFor(baseType reflect.Type) (Holder, error) {
	switch baseType {
	case timeType:
		return &nullTimeHolder{}, nil
	}
	switch baseType.Kind() {
	case reflect.Int64:
		return &nullInt64Holder{}, nil
	case reflect.Float64:
		return &nullFloat64Holder{}, nil
	case reflect.String:
		return &nullStringHolder{}, nil
	case reflect.Bool:
		return &nullBoolHolder{}, nil
	}
	return nil, fmt.Errorf("raizel: unsupported nullable type *%s — supported: *time.Time, *int64, *float64, *string, *bool", baseType)
}
