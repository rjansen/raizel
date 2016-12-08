package cassandra

import (
	"database/sql"
	"farm.e-pedion.com/repo/persistence"
	testify "github.com/stretchr/testify/mock"
)

func NewClientPoolMock() *ClientPoolMock {
	return new(ClientPoolMock)
}

//ClientPoolMock is a mock for a cache client pool
type ClientPoolMock struct {
	testify.Mock
}

//Get returns a cache Client instance
func (m *ClientPoolMock) Get() (persistence.Client, error) {
	args := m.Called()
	result := args.Get(0)
	if result != nil {
		return result.(persistence.Client), args.Error(1)
	}
	return nil, args.Error(1)
}

//Close finalizes the pool instance
func (m *ClientPoolMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

//NewDBMock creates a new mock instance for the sql DB component
func NewDBMock() *DBMock {
	return new(DBMock)
}

//DBMock is a mock client for database/sql contract
type DBMock struct {
	testify.Mock
}

func (m *DBMock) QueryRow(sql string, params ...interface{}) Row {
	args := m.Called(sql, params)
	result := args.Get(0)
	if result != nil {
		return result.(Row)
	}
	return nil
}

func (m *DBMock) Query(sql string, params ...interface{}) (Rows, error) {
	args := m.Called(sql, params)
	result := args.Get(0)
	if result != nil {
		return result.(Rows), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *DBMock) Exec(query string, params ...interface{}) (sql.Result, error) {
	args := m.Called(query, params)
	result := args.Get(0)
	if result != nil {
		return result.(sql.Result), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *DBMock) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DBMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func NewRowMock() *RowMock {
	return new(RowMock)
}

type RowMock struct {
	testify.Mock
}

func (m *RowMock) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

func NewResultMock() *ResultMock {
	return new(ResultMock)
}

type ResultMock struct {
	testify.Mock
}

func (m *ResultMock) LastInsertId() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *ResultMock) RowsAffected() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
