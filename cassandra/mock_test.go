package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/stretchr/testify/mock"
)

type sessionMock struct {
	mock.Mock
}

func newSessionMock() *sessionMock {
	return new(sessionMock)
}
func (mock *sessionMock) Close() {
	_ = mock.Called()
}
func (mock *sessionMock) Query(cql string, arguments ...interface{}) Query {
	var (
		args   = mock.Called(cql, arguments)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(Query)
}

func (mock *sessionMock) Closed() bool {
	args := mock.Called()
	return args.Bool(0)
}

type queryMock struct {
	mock.Mock
}

func newQueryMock() *queryMock {
	return new(queryMock)
}

func (mock *queryMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func (mock *queryMock) Exec() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *queryMock) Iter() Iter {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(Iter)
}

func (mock *queryMock) Consistency(consistency gocql.Consistency) Query {
	var (
		args   = mock.Called(consistency)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(Query)
}

func (mock *queryMock) PageSize(size int) Query {
	var (
		args   = mock.Called(size)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(Query)
}

func (mock *queryMock) Release() {
	mock.Called()
}

func (mock *queryMock) String() string {
	args := mock.Called()
	return args.String(0)
}

type iterMock struct {
	mock.Mock
}

func newIterMock() *iterMock {
	return new(iterMock)
}

func (mock *iterMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *iterMock) NumRows() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *iterMock) Scanner() gocql.Scanner {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(gocql.Scanner)
}

type scannerMock struct {
	mock.Mock
}

func newScannerMock() *scannerMock {
	return new(scannerMock)
}

func (mock *scannerMock) Next() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *scannerMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func (mock *scannerMock) Err() error {
	args := mock.Called()
	return args.Error(0)
}
