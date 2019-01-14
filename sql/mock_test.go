package sql

import (
	"github.com/stretchr/testify/mock"
)

func newDBMock() *dbMock {
	return new(dbMock)
}

type dbMock struct {
	mock.Mock
}

func (mock *dbMock) QueryRow(sql string, params ...interface{}) Row {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Row)
	}
	return nil
}

func (mock *dbMock) Query(sql string, params ...interface{}) (Rows, error) {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Rows), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *dbMock) Exec(query string, params ...interface{}) (Result, error) {
	var (
		args   = mock.Called(query, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Result), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *dbMock) Ping() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *dbMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func newRowMock() *rowMock {
	return new(rowMock)
}

type rowMock struct {
	mock.Mock
}

func (mock *rowMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func newRowsMock() *rowsMock {
	return new(rowsMock)
}

type rowsMock struct {
	mock.Mock
}

func (mock *rowsMock) Next() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *rowsMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func newResultMock() *resultMock {
	return new(resultMock)
}

type resultMock struct {
	mock.Mock
}

func (mock *resultMock) LastInsertId() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (mock *resultMock) RowsAffected() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}
