// +build integration

package sql

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	sqlbuilder "github.com/huandu/go-sqlbuilder"
	_ "github.com/lib/pq"
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

const (
	psqlInsertEntityMock = `
		insert into entity_mock(name, age, data) values ($1, $2, $3)
		returning id, name, age, data, deleted, created_at, updated_at
	`

	psqlDeleteEntityMock = `
		delete from entity_mock where id = $1
	`
)

type testRepositoryPostgresGet struct {
	name     string
	tree     yggdrasil.Tree
	mapper   Mapper
	mockData *entityMock
	key      entityKeyMock
	result   raizel.Entity
	err      error
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

	if scenario.err == nil {
		scenario.mockData = new(entityMock)
		var (
			sqlStruct = scenario.mapper.Get(scenario.key.EntityName())
			row       = db.QueryRow(
				psqlInsertEntityMock,
				scenario.name, 777, dynamicData{"key": "value"},
			)
			err = row.Scan(sqlStruct.Addr(scenario.mockData)...)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryPostgresGet) tearDown(t *testing.T) {
	if scenario.mockData != nil {
		var (
			db     = MustReference(scenario.tree)
			_, err = db.Exec(psqlDeleteEntityMock, scenario.mockData.ID)
		)
		require.Nil(t, err, "teardown error")
	}
	scenario.tree.Close()
}

func TestRepositoryPostgresGet(test *testing.T) {
	scenarios := []testRepositoryPostgresGet{
		{
			name: "When try to Get an entity on database returns it successfully",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_mock",
					sqlbuilder.NewStruct(new(entityMock)).For(sqlbuilder.PostgreSQL),
				).NewMapper(),
		},
		{
			name: "When try to Get an entity on database but an error is returned",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
				value: 133,
			},
			result: &entityMock{},
			mapper: NewMapperBuilder().
				Set("entity_mock",
					sqlbuilder.NewStruct(new(entityMock)).For(sqlbuilder.PostgreSQL),
				).NewMapper(),
			err: raizel.ErrNotFound,
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.tree, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
			},
		)
	}
}

type testRepositoryPostgresSet struct {
	name     string
	tree     yggdrasil.Tree
	mapper   Mapper
	mockData *entityMock
	key      entityKeyMock
	data     *entityMock
	err      error
}

func (scenario *testRepositoryPostgresSet) setup(t *testing.T) {
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

	if scenario.mockData != nil {
		var (
			sqlStruct = scenario.mapper.Get(scenario.key.EntityName())
			row       = db.QueryRow(
				psqlInsertEntityMock,
				scenario.mockData.Name, scenario.mockData.Age, scenario.mockData.Data,
			)
			err = row.Scan(sqlStruct.Addr(scenario.mockData)...)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
		scenario.data.ID = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryPostgresSet) tearDown(t *testing.T) {
	db := MustReference(scenario.tree)
	if scenario.mockData != nil {
		_, err := db.Exec(psqlDeleteEntityMock, scenario.mockData.ID)
		require.Nil(t, err, "teardown mock error")
	}
	if scenario.data != nil {
		_, err := db.Exec(psqlDeleteEntityMock, scenario.data.ID)
		require.Nil(t, err, "teardown data error")
	}
	scenario.tree.Close()
}

func TestRepositoryPostgresSet(test *testing.T) {
	currentTime := time.Now().UTC()
	scenarios := []testRepositoryPostgresSet{
		{
			name: "When try to Set a new entity on database successfully",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
			},
			data: &entityMock{
				Name: "When try to Set a new entity on database successfully",
				Age:  10,
				Data: map[string]interface{}{
					"golangkey": "golangvalue",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			mapper: NewMapperBuilder().
				Set("entity_mock",
					sqlbuilder.NewStruct(new(entityMock)).For(sqlbuilder.PostgreSQL),
				).NewMapper(),
		},
		{
			name: "When try to Set an existent entity on database successfully",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
			},
			mockData: &entityMock{
				Name: "When try to Set an existent entity on database successfully",
				Age:  66,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			data: &entityMock{
				Name: "Changed to: 'When try to Set an existent entity on database successfully'",
				Age:  333,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
					"changedkey":  "1",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			mapper: NewMapperBuilder().
				Set("entity_mock",
					sqlbuilder.NewStruct(new(entityMock)).For(sqlbuilder.PostgreSQL),
				).NewMapper(),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.tree, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
			},
		)
	}
}

type testRepositoryPostgresDelete struct {
	name     string
	tree     yggdrasil.Tree
	mapper   Mapper
	mockData *entityMock
	key      entityKeyMock
	err      error
}

func (scenario *testRepositoryPostgresDelete) setup(t *testing.T) {
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

	if scenario.err == nil {
		scenario.mockData = new(entityMock)
		var (
			sqlStruct = scenario.mapper.Get(scenario.key.EntityName())
			row       = db.QueryRow(
				psqlInsertEntityMock,
				scenario.mockData.Name, scenario.mockData.Age, scenario.mockData.Data,
			)
			err = row.Scan(sqlStruct.Addr(scenario.mockData)...)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryPostgresDelete) tearDown(t *testing.T) {
	scenario.tree.Close()
}

func TestRepositoryPostgresDelete(test *testing.T) {
	scenarios := []testRepositoryPostgresDelete{
		{
			name: "When try to Delete an entity on database successfully",
			key: entityKeyMock{
				table: "entity_mock",
				name:  "id",
			},
			mockData: &entityMock{
				Name: "When try to Delete an entity on database successfully",
				Age:  96,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
					"deletedkey":  "deletedvalue",
				},
			},
			mapper: NewMapperBuilder().
				Set("entity_mock",
					sqlbuilder.NewStruct(new(entityMock)).For(sqlbuilder.PostgreSQL),
				).NewMapper(),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository(scenario.mapper)
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.tree, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
			},
		)
	}
}
