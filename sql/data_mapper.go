package sql

import (
	"reflect"
)

type DataMapper struct {
	Name   string
	Type   reflect.Type
	Fields Fields
}

type Fields struct {
	Scalars []ScalarField
	// nested?
	Nesteds []NestedField
	Slices  []SliceField
}

type Field struct {
	Name string
	// Type T
	Type reflect.Type
}

type ScalarField struct {
	Field
}

type NestedField struct {
	Field
}

type SliceField struct {
	Field
}

func NewDataMapper[T any]() (DataMapper, error) {
	var (
		dataMapper DataMapper
		data       T
		dataType   = reflect.TypeOf(data)
	)
	// for k := 0; k < dataType.NumField(); k++ {
	for _, field := range reflect.VisibleFields(dataType) {
		if !field.IsExported() {
			continue
		}
		switch fieldType := field.Type; fieldType.Kind() {
		case reflect.Invalid:
			return dataMapper, nil
		case reflect.Chan, reflect.Func:
			// do nothing
		case reflect.Interface, reflect.Pointer, reflect.Struct:
			switch typeName := fieldType.Name(); typeName {
			// case "time.Time":
			case "Time":
				dataMapper.Fields.Scalars = append(
					dataMapper.Fields.Scalars,
					ScalarField{Field{Name: fieldType.Name(), Type: fieldType}},
				)
				// panic(typeName)
			default:
				dataMapper.Fields.Nesteds = append(
					dataMapper.Fields.Nesteds,
					NestedField{Field{Name: fieldType.Name(), Type: fieldType}},
				)

			}
		case reflect.Array, reflect.Slice:
			dataMapper.Fields.Slices = append(
				dataMapper.Fields.Slices,
				SliceField{Field{Name: fieldType.Name(), Type: fieldType}},
			)

			// slice fields
		default:
			dataMapper.Fields.Scalars = append(
				dataMapper.Fields.Scalars,
				ScalarField{Field{Name: fieldType.Name(), Type: fieldType}},
			)
		}
	}

	return dataMapper, nil
	//	errors.New("not_implemented")
}
