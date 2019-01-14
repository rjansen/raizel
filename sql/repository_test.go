package sql

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

type testEntityKey struct {
	pkField string
	id      string
}

func (k testEntityKey) GetKeyValue() interface{} {
	return k.id
}

func (k testEntityKey) GetEntityName() string {
	return k.pkField
}

func TestNewRepository(test *testing.T) {
	repository := NewRepository()
	require.NotNil(test, repository, "invalid repository instance")
}

type testRepositoryGet struct {
	name   string
	tree   yggdrasil.Tree
	row    *rowMock
	db     *dbMock
	key    raizel.EntityKey
	result raizel.Entity
	err    error
}

func (scenario *testRepositoryGet) setup(t *testing.T) {
	var (
		row   = newRowMock()
		db    = newDBMock()
		roots = yggdrasil.NewRoots()
		err   = Register(&roots, db)
	)
	require.NotNil(t, db, "mock db instance")
	require.NotNil(t, roots, "roots instance")
	require.Nil(t, err, "register db err")

	row.On("Scan", mock.Anything).Return(scenario.err)
	db.On("QueryRow", mock.AnythingOfType("string"), mock.Anything).Return(row)
	db.On("Close").Return(nil)

	scenario.row = row
	scenario.db = db
	scenario.tree = roots.NewTreeDefault()
}

func TestRepositoryGet(test *testing.T) {
	scenarios := []testRepositoryGet{
		{
			name: "Get entity",
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
			},
			result: &testEntity{},
		},
		{
			name: "Error when try to Get an entity",
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
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

				repository := NewRepository()
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.tree, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
				err = repository.Close(scenario.tree)
				require.Nil(t, err, "close error")
				scenario.db.AssertExpectations(t)
				scenario.row.AssertExpectations(t)
			},
		)
	}
}

type testRepositorySet struct {
	name   string
	tree   yggdrasil.Tree
	result *resultMock
	db     *dbMock
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
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
			},
			data: &testEntity{},
		},
		{
			name: "Error when try to Set an entity",
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
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

				repository := NewRepository()
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
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
			},
		},
		{
			name: "Error when try to Delete an entity",
			key: testEntityKey{
				pkField: "id",
				id:      "identifier",
			},
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := NewRepository()
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
