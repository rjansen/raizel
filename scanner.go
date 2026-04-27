package raizel

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"
	"time"
)

// timeType is recognised as a scalar (not a struct) by the scanner.
var timeType = reflect.TypeFor[time.Time]()

// fieldIndex locates one column-bound struct field. Round 2 uses only the
// shallow `index` value; Round 3 will extend this struct for nested paths
// and nullable indirection.
type fieldIndex struct {
	index int // struct field position; -1 marks an unmapped column
}

// fieldCache memoises the column→fieldIndex map per struct type. Re-running
// the reflect walk for every query becomes the dominant cost on hot paths
// otherwise.
var fieldCache sync.Map // map[reflect.Type]map[string]fieldIndex

// fieldsByColumn returns t's `db` tag → fieldIndex map. Tag values are
// lowercased so column lookups are case-insensitive (Oracle-friendly).
func fieldsByColumn(t reflect.Type) map[string]fieldIndex {
	if v, ok := fieldCache.Load(t); ok {
		return v.(map[string]fieldIndex)
	}
	out := make(map[string]fieldIndex, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("db")
		if tag == "" || tag == "-" {
			continue
		}
		if c := strings.IndexByte(tag, ','); c >= 0 {
			tag = tag[:c]
		}
		out[strings.ToLower(tag)] = fieldIndex{index: i}
	}
	fieldCache.Store(t, out)
	return out
}

// rowScanner holds the per-query scan plan. It is constructed once after
// rows.Columns() and reused across every row in the result set.
type rowScanner[T any] struct {
	isStruct bool
	fields   []fieldIndex // aligned with rows.Columns(); index<0 means discard
}

func newRowScanner[T any](rows *sql.Rows) (*rowScanner[T], error) {
	s := &rowScanner[T]{isStruct: isStructType[T]()}
	if !s.isStruct {
		return s, nil
	}
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	s.fields = make([]fieldIndex, len(cols))
	fmap := fieldsByColumn(reflect.TypeFor[T]())
	for i, c := range cols {
		if fi, ok := fmap[strings.ToLower(c)]; ok {
			s.fields[i] = fi
		} else {
			s.fields[i] = fieldIndex{index: -1}
		}
	}
	return s, nil
}

func (s *rowScanner[T]) scan(rows *sql.Rows) (T, error) {
	var dst T
	if !s.isStruct {
		if err := rows.Scan(&dst); err != nil {
			return dst, err
		}
		return dst, nil
	}
	rv := reflect.ValueOf(&dst).Elem()
	args := make([]any, len(s.fields))
	for i, fi := range s.fields {
		if fi.index < 0 {
			args[i] = new(any) // unmapped — discard into a fresh sink
			continue
		}
		args[i] = rv.Field(fi.index).Addr().Interface()
	}
	if err := rows.Scan(args...); err != nil {
		return dst, err
	}
	return dst, nil
}

// isStructType reports whether T should be scanned as a tagged struct
// rather than as a single-column scalar. time.Time is treated as a scalar.
func isStructType[T any]() bool {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct && t != timeType
}
