// Package named tokenises `:name` placeholders in a SQL string and binds
// them to struct fields tagged with `db:"name"`. It is dialect-agnostic:
// the caller passes in a placeholder-rendering callback, so the package
// has no dependency on the public Dialect type.
package named

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var timeType = reflect.TypeFor[time.Time]()

// Rewrite walks query, replacing every `:name` token with placeholder(n)
// (1-indexed) and collecting the corresponding values from model's `db`
// tags. Postgres-style `::` type casts are passed through verbatim.
func Rewrite(query string, placeholder func(n int) string, model any) (string, []any, error) {
	values, err := tagValues(model)
	if err != nil {
		return "", nil, err
	}

	var out strings.Builder
	out.Grow(len(query) + 8)
	var params []any
	n := 1

	for i := 0; i < len(query); {
		c := query[i]
		if c != ':' {
			out.WriteByte(c)
			i++
			continue
		}
		// `::` is a Postgres type cast — keep it literal.
		if i+1 < len(query) && query[i+1] == ':' {
			out.WriteString("::")
			i += 2
			continue
		}
		// Tokenize the name following the colon.
		j := i + 1
		for j < len(query) && isNameChar(query[j]) {
			j++
		}
		if j == i+1 {
			// Lone ':' (e.g. inside a string literal we are not parsing).
			out.WriteByte(':')
			i++
			continue
		}
		name := query[i+1 : j]
		v, ok := values[name]
		if !ok {
			return "", nil, fmt.Errorf("named param :%s not found in db tags of %T", name, model)
		}
		out.WriteString(placeholder(n))
		n++
		params = append(params, v)
		i = j
	}
	return out.String(), params, nil
}

// tagValues flattens model's exported `db`-tagged fields into a name→value
// map. Pointer fields contribute the dereferenced value or nil. Nested
// structs are skipped — named params take the flat tag form only.
func tagValues(model any) (map[string]any, error) {
	rv := reflect.ValueOf(model)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, errors.New("raizel.ExecNamed: nil model")
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return nil, fmt.Errorf("raizel.ExecNamed: expected struct, got %s", rv.Kind())
	}
	rt := rv.Type()
	out := make(map[string]any, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
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
		fv := rv.Field(i)
		ft := f.Type
		if ft.Kind() == reflect.Pointer {
			if fv.IsNil() {
				out[tag] = nil
			} else {
				out[tag] = fv.Elem().Interface()
			}
			continue
		}
		if ft.Kind() == reflect.Struct && ft != timeType {
			continue
		}
		out[tag] = fv.Interface()
	}
	return out, nil
}

func isNameChar(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}
