package sql

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRepository(test *testing.T) {
	repository := NewRepository(nil, NewMapperBuilder().NewMapper())
	require.NotNil(test, repository, "invalid repository instance")
}

type testRepositoryGet struct {
	name   string
	ctx    context.Context
	row    *rowMock
	db     *dbMock
	mapper Mapper
	key    raizel.EntityKey
	result raizel.Entity
	err    error
}

func (scenario *testRepositoryGet) setup(t *testing.T) {
	var (
		row = newRowMock()
		db  = newDBMock()
	)
	require.NotNil(t, row, "mock row instance")
	require.NotNil(t, db, "mock db instance")

	row.On("Scan", mock.Anything).Return(scenario.err)
	db.On("QueryRow", mock.AnythingOfType("string"), mock.Anything).Return(row)
	db.On("Close").Return(nil)

	scenario.row = row
	scenario.db = db
	scenario.ctx = context.Background()
}

func TestRepositoryGet(test *testing.T) {
	scenarios := []testRepositoryGet{
		{
			name: "Get entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
		{
			name: "Error when try to Get an entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.db, scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.ctx, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
				err = repository.Close(scenario.ctx)
				require.Nil(t, err, "close error")
				scenario.db.AssertExpectations(t)
				scenario.row.AssertExpectations(t)
			},
		)
	}
}

type testRepositorySet struct {
	name   string
	ctx    context.Context
	result *resultMock
	db     *dbMock
	mapper Mapper
	key    raizel.EntityKey
	data   raizel.Entity
	err    error
}

func (scenario *testRepositorySet) setup(t *testing.T) {
	var (
		result = newResultMock()
		db     = newDBMock()
	)
	require.NotNil(t, result, "mock result instance")
	require.NotNil(t, db, "mock db instance")

	// result.On("LastInsertId").Return(1, nil)
	// result.On("AffectedRows").Return(1, nil)
	db.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(result, scenario.err)
	db.On("Close").Return(nil)

	scenario.db = db
	scenario.ctx = context.Background()
}

func TestRepositorySet(test *testing.T) {
	scenarios := []testRepositorySet{
		{
			name: "Set entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			data: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
		{
			name: "Error when try to Set an entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			data: &entityMock{},
			err:  errors.New("errMock"),
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.db, scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.ctx, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
				err = repository.Close(scenario.ctx)
				require.Nil(t, err, "close error")
				scenario.db.AssertExpectations(t)
				// scenario.result.AssertExpectations(t)
			},
		)
	}
}

type testRepositoryDelete struct {
	name   string
	ctx    context.Context
	result *resultMock
	db     *dbMock
	mapper Mapper
	key    raizel.EntityKey
	data   raizel.Entity
	err    error
}

func (scenario *testRepositoryDelete) setup(t *testing.T) {
	var (
		result = newResultMock()
		db     = newDBMock()
	)
	require.NotNil(t, result, "mock result instance")
	require.NotNil(t, db, "mock db instance")

	// result.On("LastInsertId").Return(1, nil)
	// result.On("AffectedRows").Return(1, nil)
	db.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(result, scenario.err)
	db.On("Close").Return(nil)

	scenario.db = db
	scenario.ctx = context.Background()
}

func TestRepositoryDelete(test *testing.T) {
	scenarios := []testRepositoryDelete{
		{
			name: "Delete entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
		{
			name: "Error when try to Delete an entity",
			key: entityKeyMock{
				table: "entity_table",
				name:  "id",
				value: "identifier",
			},
			err: errors.New("errMock"),
			mapper: NewMapperBuilder().
				Set("entity_table", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.db, scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.ctx, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
				repository.Close(scenario.ctx)
				scenario.db.AssertExpectations(t)
				// scenario.ref.AssertExpectations(t)
			},
		)
	}
}
