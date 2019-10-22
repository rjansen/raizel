package cassandra

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type testEntityKey struct {
	entityName string
	name       string
	value      interface{}
}

func (k testEntityKey) EntityName() string {
	return k.entityName
}

func (k testEntityKey) Value() interface{} {
	return k.value
}

func (k testEntityKey) Name() string {
	return k.name
}

func TestNewRepository(test *testing.T) {
	repository := NewRepository(nil)
	require.NotNil(test, repository, "invalid repository instance")
}

type testRepositoryGet struct {
	name    string
	ctx     context.Context
	query   *queryMock
	session *sessionMock
	key     raizel.EntityKey
	result  raizel.Entity
	err     error
}

func (scenario *testRepositoryGet) setup(t *testing.T) {
	var (
		query   = newQueryMock()
		session = newSessionMock()
	)
	require.NotNil(t, query, "mock query instance")
	require.NotNil(t, session, "mock session instance")

	query.On("Scan", mock.Anything).Return(scenario.err)
	// query.On("delegate").Return(new(gocql.Query))
	session.On("Query", mock.AnythingOfType("string"), mock.Anything).Return(query)
	session.On("Close")

	scenario.query = query
	scenario.session = session
	scenario.ctx = context.Background()
}

func TestRepositoryGet(test *testing.T) {
	scenarios := []testRepositoryGet{
		{
			name: "Get entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
			result: &testEntity{},
		},
		{
			name: "Error when try to Get an entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
			result: &testEntity{},
			err:    errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.session)
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.ctx, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
				err = repository.Close(scenario.ctx)
				require.Nil(t, err, "close error")
				scenario.session.AssertExpectations(t)
				scenario.query.AssertExpectations(t)
			},
		)
	}
}

type testRepositorySet struct {
	name    string
	ctx     context.Context
	query   *queryMock
	session *sessionMock
	key     raizel.EntityKey
	data    raizel.Entity
	err     error
}

func (scenario *testRepositorySet) setup(t *testing.T) {
	var (
		query   = newQueryMock()
		session = newSessionMock()
	)
	require.NotNil(t, query, "mock query instance")
	require.NotNil(t, session, "mock session instance")

	query.On("Exec").Return(scenario.err)
	// query.On("delegate").Return(new(gocql.Query))
	session.On("Query", mock.AnythingOfType("string"), mock.Anything).Return(query)
	session.On("Close")

	scenario.query = query
	scenario.session = session
	scenario.ctx = context.Background()
}

func TestRepositorySet(test *testing.T) {
	scenarios := []testRepositorySet{
		{
			name: "Set entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
			data: &testEntity{},
		},
		{
			name: "Error when try to Set an entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
			data: &testEntity{},
			err:  errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.session)
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.ctx, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
				err = repository.Close(scenario.ctx)
				require.Nil(t, err, "close error")
				scenario.session.AssertExpectations(t)
				scenario.query.AssertExpectations(t)
			},
		)
	}
}

type testRepositoryDelete struct {
	name    string
	ctx     context.Context
	query   *queryMock
	session *sessionMock
	key     raizel.EntityKey
	data    raizel.Entity
	err     error
}

func (scenario *testRepositoryDelete) setup(t *testing.T) {
	var (
		query   = newQueryMock()
		session = newSessionMock()
	)
	require.NotNil(t, query, "mock query instance")
	require.NotNil(t, session, "mock session instance")

	query.On("Exec").Return(scenario.err)
	// query.On("delegate").Return(new(gocql.Query))
	session.On("Query", mock.AnythingOfType("string"), mock.Anything).Return(query)
	session.On("Close")

	scenario.query = query
	scenario.session = session
	scenario.ctx = context.Background()
}

func TestRepositoryDelete(test *testing.T) {
	scenarios := []testRepositoryDelete{
		{
			name: "Delete entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
		},
		{
			name: "Error when try to Delete an entity",
			key: testEntityKey{
				entityName: "testEntityKey",
				name:       "id",
				value:      "identifier",
			},
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository(scenario.session)
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.ctx, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
				repository.Close(scenario.ctx)
				scenario.session.AssertExpectations(t)
				scenario.query.AssertExpectations(t)
			},
		)
	}
}
