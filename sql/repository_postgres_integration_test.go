// +build integration

package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

type testRepositoryPostgresGet struct {
	name   string
	tree   yggdrasil.Tree
	mapper Mapper
	key    raizel.EntityKey
	result raizel.Entity
	err    error
}

func (scenario *testRepositoryPostgresGet) setup(t *testing.T) {
	var (
		driver          = "postgres"
		dsn             = "postgres://postgres:@127.0.0.1:5432/postgres?sslmode=disable"
		sqlDB, errSqlDB = sql.Open(driver, dsn)
		db, errDB       = newDB(sqlDB)
		roots           = yggdrasil.NewRoots()
		err             = Register(&roots, db)
	)
	require.Nil(t, errSqlDB, "sqlopen error")
	require.Nil(t, errDB, "newdb error")
	require.Nil(t, err, "register db err")
	require.NotNil(t, roots, "roots instance")
	require.NotNil(t, db, "db instance")

	scenario.tree = roots.NewTreeDefault()
}

func TestRepositoryPostgresGet(test *testing.T) {
	scenarios := []testRepositoryPostgresGet{
		{
			name: "Get entity",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
				value: 111,
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_mock", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
		},
		{
			name: "Error when try to Get an entity",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
				value: 222,
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_mock", sqlbuilder.NewStruct(new(entityMock))).
				NewMapper(),
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.tree, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
				err = repository.Close(scenario.tree)
				require.Nil(t, err, "close error")
			},
		)
	}
}

/*
type testRepositorySet struct {
	name   string
	tree   yggdrasil.Tree
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
		roots  = yggdrasil.NewRoots()
		err    = Register(&roots, db)
	)
	require.NotNil(t, result, "mock result instance")
	require.NotNil(t, db, "mock db instance")
	require.NotNil(t, roots, "roots instance")
	require.Nil(t, err, "register db err")

	// result.On("LastInsertId").Return(1, nil)
	// result.On("AffectedRows").Return(1, nil)
	db.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(result, scenario.err)
	db.On("Close").Return(nil)

	scenario.db = db
	scenario.tree = roots.NewTreeDefault()
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

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.tree, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
				err = repository.Close(scenario.tree)
				require.Nil(t, err, "close error")
				scenario.db.AssertExpectations(t)
				// scenario.result.AssertExpectations(t)
			},
		)
	}
}

type testRepositoryDelete struct {
	name   string
	tree   yggdrasil.Tree
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
		roots  = yggdrasil.NewRoots()
		err    = Register(&roots, db)
	)
	require.NotNil(t, result, "mock result instance")
	require.NotNil(t, db, "mock db instance")
	require.NotNil(t, roots, "roots instance")
	require.Nil(t, err, "register db err")

	// result.On("LastInsertId").Return(1, nil)
	// result.On("AffectedRows").Return(1, nil)
	db.On("Exec", mock.AnythingOfType("string"), mock.Anything).Return(result, scenario.err)
	db.On("Close").Return(nil)

	scenario.db = db
	scenario.tree = roots.NewTreeDefault()
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

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.tree, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
				repository.Close(scenario.tree)
				scenario.db.AssertExpectations(t)
				// scenario.ref.AssertExpectations(t)
			},
		)
	}
}
*/
