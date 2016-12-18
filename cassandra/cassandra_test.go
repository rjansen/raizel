package cassandra

import (
	"errors"
	"github.com/gocql/gocql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func init() {
	if err := l.Setup(&l.Configuration{}); err != nil {
		panic(err)
	}
}

func TestUnitDelegateSession(t *testing.T) {
	assert.NotPanics(t,
		func() {
			session := NewDelegateSession(new(gocql.Session))
			assert.NotNil(t, session)
			query := session.Query("mockCQL", 1, "param2", 0.12)
			assert.NotNil(t, query)
			assert.False(t, session.Closed())
			session.Close()
		},
	)
}

func TestUnitDelegateQuery(t *testing.T) {
	assert.NotPanics(t,
		func() {
			query := NewDelegateQuery(new(gocql.Query))
			assert.NotNil(t, query)
			pageSizeQuery := query.PageSize(10)
			assert.NotNil(t, pageSizeQuery)
			assert.Equal(t, query, pageSizeQuery)

			consistencyQuery := query.Consistency(gocql.One)
			assert.NotNil(t, consistencyQuery)
			assert.Equal(t, pageSizeQuery, consistencyQuery)
			assert.Equal(t, query, consistencyQuery)

			assert.Panics(t,
				func() {
					query.Exec()
				},
			)
			assert.Panics(t,
				func() {
					query.Scan(nil)
				},
			)

			query.Release()
		},
	)
}

func TestUnitClient(t *testing.T) {
	assert.NotPanics(t,
		func() {
			client := NewClient(NewSessionMock())
			assert.Nil(t, client.Close())
		},
	)
}

func TestUnitClientPool(t *testing.T) {
	sessionMock := NewSessionMock()
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Close").Return(nil)
	pool = &Pool{
		cluster: new(gocql.ClusterConfig),
		session: sessionMock,
	}
	assert.Nil(t, raizel.Setup(pool))
	Config = &Configuration{}
	pool, err := raizel.GetPool()
	assert.Nil(t, err)
	assert.NotNil(t, pool)
	client, err := pool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, client)
	err = client.Close()
	assert.Nil(t, err)
	err = pool.Close()
	assert.Nil(t, err)
}

func TestUnitQueryOneExec(t *testing.T) {
	mockQuery := NewQueryMock()
	mockQuery.On("Consistency", mock.Anything).Return(mockQuery)
	mockQuery.On("Scan", mock.Anything).Return(nil)
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan(nil)
		}, "mockValue")
	assert.Nil(t, err)
}

func TestUnitQueryOneExecErr(t *testing.T) {
	mockQuery := NewQueryMock()
	mockQuery.On("Consistency", mock.Anything).Return(mockQuery)
	mockErr := errors.New("FetchMockErr")
	mockQuery.On("Scan", mock.Anything).Return(mockErr)
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			scanErr := f.Scan(nil)
			assert.Equal(t, mockErr, scanErr)
			return scanErr
		}, "mockValue")
	assert.NotNil(t, err)
}

func TestUnitQueryExec(t *testing.T) {
	mockIterable := NewIterableMock()
	records := 5
	nextRecordIdx := 0
	nextCall := mockIterable.On("Next")
	nextCall.Run(
		func(args mock.Arguments) {
			nextRecordIdx++
			if nextRecordIdx <= records {
				nextCall.ReturnArguments = []interface{}{true}
			} else {
				nextCall.ReturnArguments = []interface{}{false}
			}
		},
	)
	mockIterable.On("Scan", mock.Anything).Return(nil)
	mockIter := NewIterMock()
	mockIter.On("Close").Return(nil)
	mockIter.On("Scanner").Return(mockIterable)
	mockQuery := NewQueryMock()
	mockQuery.On("Consistency", mock.Anything).Return(mockQuery)
	mockQuery.On("Iter").Return(mockIter)
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.Query("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Iterable) error {
			assert.NotNil(t, f)
			fetchedRecords := 0
			for f.Next() {
				fetchErr := f.Scan(nil)
				assert.Nil(t, fetchErr)
				fetchedRecords++
			}
			assert.Equal(t, records, fetchedRecords)
			return nil
		}, "mockValue")
	assert.Nil(t, err)
}

func TestUnitQueryExecErr(t *testing.T) {
	mockIterable := NewIterableMock()
	mockIterable.On("Next").Return(true)
	mockErr := errors.New("mockErr")
	mockIterable.On("Scan", mock.Anything).Return(mockErr)
	mockIter := NewIterMock()
	mockIter.On("Close").Return(nil)
	mockIter.On("Scanner").Return(mockIterable)
	mockQuery := NewQueryMock()
	mockQuery.On("Consistency", mock.Anything).Return(mockQuery)
	mockQuery.On("Iter").Return(mockIter)
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.Query("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Iterable) error {
			assert.NotNil(t, f)
			return f.Scan(nil)
		}, "mockValue")
	assert.NotNil(t, err)
}

func TestUnitExec(t *testing.T) {
	mockQuery := NewQueryMock()
	mockQuery.On("Exec").Return(nil, nil)
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	_, err := persistenceClient.Exec("insert into cql.mock values (?)", "mockValue1", "mockValue2")
	assert.Nil(t, err)
}

func TestUnitExecErr(t *testing.T) {
	mockQuery := NewQueryMock()
	mockQuery.On("Exec").Return(nil, errors.New("ExecMockErr"))
	sessionMock := NewSessionMock()
	sessionMock.On("Close")
	sessionMock.On("Closed").Return(false)
	sessionMock.On("Query", mock.Anything, mock.Anything).Return(mockQuery)
	persistenceClient := NewClient(sessionMock)
	assert.NotNil(t, persistenceClient)
	_, err := persistenceClient.Exec("insert into cql.mock values (?, ?)", "mockValue", "anotherMockValue")
	assert.NotNil(t, err)
}
