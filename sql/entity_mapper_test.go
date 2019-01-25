package sql

import (
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/require"
)

func TestEntityMapper(test *testing.T) {
	var (
		entities = map[raizel.EntityKey]raizel.Entity{
			entityKeyMock{
				table: "mock_table_1",
				name:  "pk_mock_1",
				value: 111,
			}: entityKeyMock{},
			entityKeyMock{
				table: "mock_table_2",
				name:  "pk_mock_2",
				value: 222,
			}: entityKeyMock{},
			entityKeyMock{
				table: "mock_table_3",
				name:  "pk_mock_3",
				value: 333,
			}: entityKeyMock{},
			entityKeyMock{
				table: "mock_table_9",
				name:  "pk_mock_9",
				value: 999,
			}: entityKeyMock{},
		}
	)
	builder := NewMapperBuilder()
	for key, entity := range entities {
		builder.Set(key.EntityName(), sqlbuilder.NewStruct(entity))
	}
	mapper := builder.NewMapper()
	require.NotNil(test, mapper, "mapper invalid instance")

	nilEntityMapper := mapper.Get("invalid_entity_name")
	require.Nil(test, nilEntityMapper, "nilentitymapper invalid instance")

	for key, _ := range entities {
		entityMapper := mapper.Get(key.EntityName())
		require.NotNil(test, entityMapper, "entitymapper invalid instance")
	}
}
