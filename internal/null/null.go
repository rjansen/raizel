// Package null shuttles single columns between database/sql's Null{Time,
// Int64,Float64,String,Bool} scan targets and either (a) pointer-typed struct
// fields (*time.Time, *int64, ...), where NULL leaves a nil pointer, or (b)
// non-pointer scalar fields tagged `,nullzero`, where NULL coerces to the
// field's zero value. It is internal plumbing used by the row scanner; the
// public raizel API does not expose it directly.
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

// The zero* holders shuttle a single column into a NON-pointer scalar field
// (string, int64, ...) tagged `,nullzero`. A present value is written to the
// field; a NULL column leaves the field at its zero value. Unlike the pointer
// holders above, NULL→zero coercion here is opt-in per field, so it never
// silently masks an unexpected NULL on an untagged column.
type zeroStringHolder struct{ v sql.NullString }

func (h *zeroStringHolder) ScanDest() any { return &h.v }
func (h *zeroStringHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		f.SetString(h.v.String)
	} else {
		f.SetString("")
	}
}

type zeroInt64Holder struct{ v sql.NullInt64 }

func (h *zeroInt64Holder) ScanDest() any { return &h.v }
func (h *zeroInt64Holder) Assign(f reflect.Value) {
	if h.v.Valid {
		f.SetInt(h.v.Int64)
	} else {
		f.SetInt(0)
	}
}

type zeroFloat64Holder struct{ v sql.NullFloat64 }

func (h *zeroFloat64Holder) ScanDest() any { return &h.v }
func (h *zeroFloat64Holder) Assign(f reflect.Value) {
	if h.v.Valid {
		f.SetFloat(h.v.Float64)
	} else {
		f.SetFloat(0)
	}
}

type zeroBoolHolder struct{ v sql.NullBool }

func (h *zeroBoolHolder) ScanDest() any { return &h.v }
func (h *zeroBoolHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		f.SetBool(h.v.Bool)
	} else {
		f.SetBool(false)
	}
}

type zeroTimeHolder struct{ v sql.NullTime }

func (h *zeroTimeHolder) ScanDest() any { return &h.v }
func (h *zeroTimeHolder) Assign(f reflect.Value) {
	if h.v.Valid {
		f.Set(reflect.ValueOf(h.v.Time))
	} else {
		f.Set(reflect.ValueOf(time.Time{}))
	}
}

// ZeroHolderFor returns a fresh zero-coercing holder for the given non-pointer
// scalar type. Unknown types fail loudly, matching HolderFor.
func ZeroHolderFor(baseType reflect.Type) (Holder, error) {
	switch baseType {
	case timeType:
		return &zeroTimeHolder{}, nil
	}
	switch baseType.Kind() {
	case reflect.String:
		return &zeroStringHolder{}, nil
	case reflect.Int64:
		return &zeroInt64Holder{}, nil
	case reflect.Float64:
		return &zeroFloat64Holder{}, nil
	case reflect.Bool:
		return &zeroBoolHolder{}, nil
	}
	return nil, fmt.Errorf("raizel: unsupported nullzero type %s — supported: string, int64, float64, bool, time.Time", baseType)
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
