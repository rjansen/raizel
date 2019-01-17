package sql

import (
	sqlbuilder "github.com/huandu/go-sqlbuilder"
)

type MapperBuilder struct {
	register map[string]*sqlbuilder.Struct
}

func (builder *MapperBuilder) Set(entityName string, structBuilder *sqlbuilder.Struct) *MapperBuilder {
	builder.register[entityName] = structBuilder
	return builder
}

type Mapper interface {
	Get(string) *sqlbuilder.Struct
}

type mapper struct {
	register map[string]*sqlbuilder.Struct
}

func (mapper mapper) Get(entityName string) *sqlbuilder.Struct {
	structBuilder, exists := mapper.register[entityName]
	if !exists {
		return nil
	}
	return structBuilder
}

func (builder *MapperBuilder) NewMapper() Mapper {
	register := make(map[string]*sqlbuilder.Struct)
	for key, structBuilder := range builder.register {
		register[key] = structBuilder
	}
	return mapper{
		register: register,
	}
}

func NewMapperBuilder() *MapperBuilder {
	return &MapperBuilder{
		register: make(map[string]*sqlbuilder.Struct),
	}
}
