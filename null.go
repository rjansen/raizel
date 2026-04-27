package raizel

import (
	"database/sql"
	"fmt"
	"reflect"
)

// nullHolder shuttles a single column between (a) database/sql's
// Null{Time,Int64,Float64,String,Bool} types and (b) a nullable struct
// field of pointer kind (*time.Time, *int64, ...). When the column scans
// NULL the holder leaves the field's zero pointer in place; otherwise it
// allocates a new value and points the field at it.
type nullHolder interface {
	scanDest() any
	assign(field reflect.Value)
}

type nullTimeHolder struct{ v sql.NullTime }

func (h *nullTimeHolder) scanDest() any { return &h.v }
func (h *nullTimeHolder) assign(f reflect.Value) {
	if h.v.Valid {
		t := h.v.Time
		f.Set(reflect.ValueOf(&t))
	}
}

type nullInt64Holder struct{ v sql.NullInt64 }

func (h *nullInt64Holder) scanDest() any { return &h.v }
func (h *nullInt64Holder) assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Int64
		f.Set(reflect.ValueOf(&x))
	}
}

type nullFloat64Holder struct{ v sql.NullFloat64 }

func (h *nullFloat64Holder) scanDest() any { return &h.v }
func (h *nullFloat64Holder) assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Float64
		f.Set(reflect.ValueOf(&x))
	}
}

type nullStringHolder struct{ v sql.NullString }

func (h *nullStringHolder) scanDest() any { return &h.v }
func (h *nullStringHolder) assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.String
		f.Set(reflect.ValueOf(&x))
	}
}

type nullBoolHolder struct{ v sql.NullBool }

func (h *nullBoolHolder) scanDest() any { return &h.v }
func (h *nullBoolHolder) assign(f reflect.Value) {
	if h.v.Valid {
		x := h.v.Bool
		f.Set(reflect.ValueOf(&x))
	}
}

// nullHolderFor returns a fresh holder for the given pointer base type.
// Unknown types fail loudly rather than scanning into the wrong shuttle —
// silent fallbacks were one of the bugs raizel exists to avoid.
func nullHolderFor(baseType reflect.Type) (nullHolder, error) {
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
