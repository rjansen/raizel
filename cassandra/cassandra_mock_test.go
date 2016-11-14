package cassandra

import (
	"farm.e-pedion.com/repo/persistence"
	"github.com/gocql/gocql"
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

//NewSessionMock creates a new mock instance for the gocql session component
func NewSessionMock() *SessionMock {
	return new(SessionMock)
}

//SessionMock is a mock client for cassandra
type SessionMock struct {
	testify.Mock
}

//Get reads the value associated with the provided key
func (m *SessionMock) Query(cql string, params ...interface{}) Query {
	args := m.Called(cql, params)
	result := args.Get(0)
	if result != nil {
		return result.(Query)
	}
	return nil
}

//Closed is a flag to check the state of the session
func (m *SessionMock) Closed() bool {
	args := m.Called()
	return args.Bool(0)
}

//Close finalizes the cassandra session
func (m *SessionMock) Close() {
	m.Called()
}

//NewQueryMock returns and initializes a new query mock instance
func NewQueryMock() *QueryMock {
	return new(QueryMock)
}

//QueryMock is mock for Query interface
type QueryMock struct {
	testify.Mock
}

func (m *QueryMock) Consistency(c gocql.Consistency) Query {
	args := m.Called(c)
	result := args.Get(0)
	if result != nil {
		return args.Get(0).(Query)
	}
	return nil
}

func (m *QueryMock) Exec() error {
	args := m.Called()
	return args.Error(0)
}

// func (m *QueryMock) Iter() Iter {
// 	args := m.Called()
// 	result := args.Get(0)
// 	if result != nil {
// 		return result.(Iter)
// 	}
// 	return nil
// }

func (m *QueryMock) PageSize(n int) Query {
	args := m.Called(n)
	result := args.Get(0)
	if result != nil {
		return result.(Query)
	}
	return nil
}

func (m *QueryMock) Release() {
	m.Called()
}

func (m *QueryMock) Scan(dest ...interface{}) error {
	args := m.Called(dest)
	return args.Error(0)
}

//func (m *QueryMock) String() string {}
//func (m *QueryMock) WithContext(ctx context.Context) Query {}

//NewIterMock returns and initializes a new iter mock instance
func NewIterMock() *IterMock {
	return new(IterMock)
}

//IterMock is a mock for the Iter interface
type IterMock struct {
	testify.Mock
}

func (m *IterMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *IterMock) NumRows() int {
	args := m.Called()
	result := args.Get(0)
	if result != nil {
		return result.(int)
	}
	return 0
}

func (m *IterMock) Scan(dest ...interface{}) bool {
	args := m.Called(dest)
	result := args.Get(0)
	if result != nil {
		return result.(bool)
	}
	return false
}
