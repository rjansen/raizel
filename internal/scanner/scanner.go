// Package scanner builds and runs the per-query plan that maps
// database/sql columns into typed Go values. It is internal plumbing
// driven by the public raizel API.
package scanner

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
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
// element type used to pick a null.Holder. `nullZero` is true when a
// non-pointer scalar leaf opts into NULL→zero coercion via the `,nullzero`
// tag option; `baseType` is the scalar type. `slice` is true when the leaf
// is an array-bound slice (e.g. []string); `sliceType` is the slice type
// used for JSON decoding on non-native-array dialects.
type fieldIndex struct {
	path      []int
	nullable  bool
	nullZero  bool
	baseType  reflect.Type
	slice     bool
	sliceType reflect.Type
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

// isArraySlice reports whether t is a slice raizel maps onto an array
// column. []byte is excluded — it is a scalar blob, not an array.
func isArraySlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice && t.Elem().Kind() != reflect.Uint8
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
		var opts string
		if c := strings.IndexByte(tag, ','); c >= 0 {
			tag, opts = tag[:c], tag[c+1:]
		}
		nullZero := hasOption(opts, "nullzero")

		path := make([]int, 0, len(basePath)+1)
		path = append(path, basePath...)
		path = append(path, i)

		ft := f.Type
		nullable := ft.Kind() == reflect.Pointer
		if nullable {
			ft = ft.Elem()
		}

		col := tag
		if prefix != "" {
			col = prefix + tag
		}

		// Array-bound slice leaf (e.g. []string): mapped specially so the
		// scanner can JSON-decode on dialects without native arrays.
		if isArraySlice(ft) {
			out[strings.ToLower(col)] = fieldIndex{
				path:      path,
				slice:     true,
				sliceType: ft,
			}
			continue
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

		out[strings.ToLower(col)] = fieldIndex{
			path: path,
			// A pointer field is already nullable; `,nullzero` is only
			// meaningful for non-pointer scalars, so it is ignored there.
			nullable: nullable,
			nullZero: nullZero && !nullable,
			baseType: ft,
		}
	}
	return out
}

// hasOption reports whether a comma-separated `db` tag option list contains
// the named option (e.g. "nullzero" in `db:"body,nullzero"`).
func hasOption(opts, name string) bool {
	for opts != "" {
		var part string
		if c := strings.IndexByte(opts, ','); c >= 0 {
			part, opts = opts[:c], opts[c+1:]
		} else {
			part, opts = opts, ""
		}
		if part == name {
			return true
		}
	}
	return false
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
	isStruct   bool
	jsonArrays bool         // decode slice fields from JSON text (non-native-array dialects)
	fields     []fieldIndex // aligned with rows.Columns(); unmapped() marks discard
}

// New builds a scan plan for type T against the columns reported by rows.
// jsonArrays selects how slice fields are read: false delegates to the
// driver's native array decoding (PostgreSQL), true decodes JSON text
// (Oracle/SQLite).
func New[T any](rows *sql.Rows, jsonArrays bool) (*RowScanner[T], error) {
	s := &RowScanner[T]{isStruct: isStructType[T](), jsonArrays: jsonArrays}
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
			if fi.nullZero {
				if _, err := null.ZeroHolderFor(fi.baseType); err != nil {
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
	sliceText := make([]*sql.NullString, len(s.fields))

	for i, fi := range s.fields {
		switch {
		case fi.unmapped():
			args[i] = new(any) // discard each unmapped column into a fresh sink
		case fi.slice:
			// Array columns come back as text through database/sql (a PG
			// array literal, or JSON for Oracle/SQLite). Capture the raw
			// string and decode after the scan.
			ns := &sql.NullString{}
			sliceText[i] = ns
			args[i] = ns
		case fi.nullable:
			h, _ := null.HolderFor(fi.baseType) // validated at scanner construction
			holders[i] = h
			args[i] = h.ScanDest()
		case fi.nullZero:
			// Non-pointer scalar tagged `,nullzero`: read through a Null
			// shuttle and assign the zero value when the column is NULL.
			h, _ := null.ZeroHolderFor(fi.baseType) // validated at scanner construction
			holders[i] = h
			args[i] = h.ScanDest()
		default:
			// Plain scalar or []byte blob: scan straight into the field.
			f := fieldByPath(rv, fi.path)
			args[i] = f.Addr().Interface()
		}
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
	for i, ns := range sliceText {
		if ns == nil {
			continue
		}
		if err := assignSlice(fieldByPath(rv, s.fields[i].path), s.fields[i].sliceType, ns, s.jsonArrays); err != nil {
			return dst, err
		}
	}
	return dst, nil
}

// assignSlice decodes an array column's text form into a slice field. On
// dialects with JSON-encoded arrays it parses JSON; otherwise it parses a
// PostgreSQL array literal. A NULL or blank column leaves the field's zero
// value (nil slice) in place.
func assignSlice(field reflect.Value, sliceType reflect.Type, ns *sql.NullString, jsonArrays bool) error {
	if !ns.Valid || strings.TrimSpace(ns.String) == "" {
		return nil
	}
	if jsonArrays {
		ptr := reflect.New(sliceType)
		if err := json.Unmarshal([]byte(ns.String), ptr.Interface()); err != nil {
			return fmt.Errorf("raizel: decode JSON array into %s: %w", sliceType, err)
		}
		field.Set(ptr.Elem())
		return nil
	}
	v, err := parsePGArray(ns.String, sliceType)
	if err != nil {
		return err
	}
	field.Set(v)
	return nil
}

// parsePGArray parses a PostgreSQL array literal (e.g. {a,"b,c",NULL})
// into a slice of sliceType, converting each element to the slice's
// element type. Quoting and backslash escapes inside double-quoted
// elements are honoured; an unquoted NULL yields the element zero value.
func parsePGArray(s string, sliceType reflect.Type) (reflect.Value, error) {
	s = strings.TrimSpace(s)
	if len(s) < 2 || s[0] != '{' || s[len(s)-1] != '}' {
		return reflect.Value{}, fmt.Errorf("raizel: malformed postgres array %q", s)
	}
	elemType := sliceType.Elem()
	out := reflect.MakeSlice(sliceType, 0, 0)
	inner := s[1 : len(s)-1]
	if strings.TrimSpace(inner) == "" {
		return out, nil
	}

	var b strings.Builder
	i, n := 0, len(inner)
	for i < n {
		b.Reset()
		quoted := false
		if inner[i] == '"' {
			quoted = true
			i++
			for i < n {
				c := inner[i]
				if c == '\\' && i+1 < n {
					b.WriteByte(inner[i+1])
					i += 2
					continue
				}
				if c == '"' {
					i++
					break
				}
				b.WriteByte(c)
				i++
			}
		} else {
			for i < n && inner[i] != ',' {
				b.WriteByte(inner[i])
				i++
			}
		}
		ev, err := convertPGElem(strings.TrimSpace(b.String()), quoted, elemType)
		if err != nil {
			return reflect.Value{}, err
		}
		out = reflect.Append(out, ev)
		if i < n && inner[i] == ',' {
			i++
		}
	}
	return out, nil
}

// convertPGElem converts one parsed array element token into a value of
// elemType. An unquoted NULL maps to the zero value.
func convertPGElem(tok string, quoted bool, elemType reflect.Type) (reflect.Value, error) {
	if !quoted && tok == "NULL" {
		return reflect.Zero(elemType), nil
	}
	v := reflect.New(elemType).Elem()
	switch elemType.Kind() {
	case reflect.String:
		v.SetString(tok)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("raizel: parse int array element %q: %w", tok, err)
		}
		v.SetInt(n)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("raizel: parse float array element %q: %w", tok, err)
		}
		v.SetFloat(f)
	case reflect.Bool:
		// PostgreSQL renders booleans as t / f in array literals.
		v.SetBool(tok == "t" || tok == "true")
	default:
		return reflect.Value{}, fmt.Errorf("raizel: unsupported array element type %s", elemType)
	}
	return v, nil
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
