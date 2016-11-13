package persistence

import (
	testify "github.com/stretchr/testify/mock"
)

//NewClientMock creates a new Cassandra Client mock
func NewClientMock() *ClientMock {
	return new(ClientMock)
}

type ReaderMock struct {
	testify.Mock
}

func (m *ReaderMock) QueryOne(query string, fetchFunc func(Fetchable) error, params ...interface{}) error {
	args := m.Called(query, fetchFunc, params)
	return args.Error(0)
}

func (m *ReaderMock) Query(query string, iterFunc func(Iterable) error, params ...interface{}) error {
	args := m.Called(query, iterFunc, params)
	return args.Error(0)
}

type ExecutorMock struct {
	testify.Mock
}

func (m *ExecutorMock) Exec(cql string, params ...interface{}) error {
	args := m.Called(cql, params)
	return args.Error(0)
}

type ClientMock struct {
	testify.Mock
}

func (m *ClientMock) QueryOne(query string, fetchFunc func(Fetchable) error, params ...interface{}) error {
	args := m.Called(query, fetchFunc, params)
	return args.Error(0)
}

func (m *ClientMock) Query(query string, iterFunc func(Iterable) error, params ...interface{}) error {
	args := m.Called(query, iterFunc, params)
	return args.Error(0)
}

func (m *ClientMock) Exec(cql string, params ...interface{}) error {
	args := m.Called(cql, params)
	return args.Error(0)
}
