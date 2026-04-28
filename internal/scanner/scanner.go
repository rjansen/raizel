// Package scanner builds and runs the per-query plan that maps
// database/sql columns into typed Go values. It is internal plumbing
// driven by the public raizel API.
package scanner

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/rjansen/raizel/internal/null"
)

// timeType is recognised as a scalar (not a struct) by the scanner.
var timeType = reflect.TypeFor[time.Time]()

// fieldIndex locates one column-bound struct field. `path` walks possibly
// nested struct types — empty path marks an unmapped column. `nullable`
// is true when the leaf field is a pointer; `baseType` is then the
// element type used to pick a null.Holder.
type fieldIndex struct {
	path     []int
	nullable bool
	baseType reflect.Type
}

func (fi fieldIndex) unmapped() bool { return len(fi.path) == 0 }

// fieldCache memoises the column→fieldIndex map per struct type. Re-running
// the reflect walk for every query becomes the dominant cost on hot paths
// otherwise.
var fieldCache sync.Map // map[reflect.Type]map[string]fieldIndex

// fieldsByColumn returns t's `db` tag → fieldIndex map. Tag values are
// lowercased so column lookups are case-insensitive (Oracle-friendly).
// Nested structs contribute dot-notation keys ("team.id") so a JOIN can
// scan into a parent struct in one shot.
func fieldsByColumn(t reflect.Type) map[string]fieldIndex {
	if v, ok := fieldCache.Load(t); ok {
		return v.(map[string]fieldIndex)
	}
	out := buildFieldMap(t, nil, "")
	fieldCache.Store(t, out)
	return out
}

func buildFieldMap(t reflect.Type, basePath []int, prefix string) map[string]fieldIndex {
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

		path := make([]int, 0, len(basePath)+1)
		path = append(path, basePath...)
		path = append(path, i)

		ft := f.Type
		nullable := ft.Kind() == reflect.Pointer
		if nullable {
			ft = ft.Elem()
		}

		// Nested struct (excluding time.Time which we treat as a scalar):
		// recurse and merge with a "tag." prefix. Nested-leaf nullability
		// is preserved — that was a bug in the implementation we forked.
		if ft.Kind() == reflect.Struct && ft != timeType {
			nested := buildFieldMap(ft, path, tag+".")
			for k, v := range nested {
				out[strings.ToLower(k)] = v
			}
			continue
		}

		col := tag
		if prefix != "" {
			col = prefix + tag
		}
		out[strings.ToLower(col)] = fieldIndex{
			path:     path,
			nullable: nullable,
			baseType: ft,
		}
	}
	return out
}

// fieldByPath walks an index path through possibly-nested struct values,
// auto-allocating any nil pointer-to-struct it traverses so the leaf
// field is addressable.
func fieldByPath(v reflect.Value, path []int) reflect.Value {
	for _, idx := range path {
		if v.Kind() == reflect.Pointer {
			if v.IsNil() {
				v.Set(reflect.New(v.Type().Elem()))
			}
			v = v.Elem()
		}
		v = v.Field(idx)
	}
	return v
}

// RowScanner holds the per-query scan plan. It is built once after
// rows.Columns() and reused across every row in the result set.
type RowScanner[T any] struct {
	isStruct bool
	fields   []fieldIndex // aligned with rows.Columns(); unmapped() marks discard
}

// New builds a scan plan for type T against the columns reported by rows.
func New[T any](rows *sql.Rows) (*RowScanner[T], error) {
	s := &RowScanner[T]{isStruct: isStructType[T]()}
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
			// Validate nullable holders eagerly: a misconfigured *T raises
			// an error before we touch driver memory.
			if fi.nullable {
				if _, err := null.HolderFor(fi.baseType); err != nil {
					return nil, err
				}
			}
		}
		// else: zero-value fieldIndex with nil path → marked unmapped
	}
	return s, nil
}

// Scan reads the next row through the prepared plan and returns a freshly
// allocated T.
func (s *RowScanner[T]) Scan(rows *sql.Rows) (T, error) {
	var dst T
	if !s.isStruct {
		if err := rows.Scan(&dst); err != nil {
			return dst, err
		}
		return dst, nil
	}

	rv := reflect.ValueOf(&dst).Elem()
	args := make([]any, len(s.fields))
	holders := make([]null.Holder, len(s.fields))

	for i, fi := range s.fields {
		if fi.unmapped() {
			args[i] = new(any) // discard each unmapped column into a fresh sink
			continue
		}
		if fi.nullable {
			h, _ := null.HolderFor(fi.baseType) // validated at scanner construction
			holders[i] = h
			args[i] = h.ScanDest()
			continue
		}
		f := fieldByPath(rv, fi.path)
		args[i] = f.Addr().Interface()
	}

	if err := rows.Scan(args...); err != nil {
		return dst, err
	}

	for i, h := range holders {
		if h == nil {
			continue
		}
		f := fieldByPath(rv, s.fields[i].path)
		h.Assign(f)
	}
	return dst, nil
}

// isStructType reports whether T should be scanned as a tagged struct
// rather than as a single-column scalar. time.Time is a struct in the
// reflect sense but is treated as a scalar here.
func isStructType[T any]() bool {
	t := reflect.TypeFor[T]()
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct && t != timeType
}
