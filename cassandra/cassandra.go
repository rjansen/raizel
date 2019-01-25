package cassandra

import (
	"errors"

	"github.com/gocql/gocql"
)

var (
	ErrBlankSession = errors.New("err_blanksession")
	// 	NotFoundErr        = gocql.ErrNotFound
)

type Session interface {
	Close()
	Query(string, ...interface{}) Query
	Closed() bool
}

type Query interface {
	Scan(...interface{}) error
	Exec() error
	Iter() Iter
	Consistency(gocql.Consistency) Query
	PageSize(int) Query
	Release()
	String() string
	delegate() *gocql.Query
}

type Iter interface {
	Close() error
	NumRows() int
	Scanner() gocql.Scanner
}

type session struct {
	*gocql.Session
}

func newSession(cqlSession *gocql.Session) (Session, error) {
	if cqlSession == nil {
		return nil, ErrBlankSession
	}
	return &session{Session: cqlSession}, nil
}

func (session *session) Query(cql string, arguments ...interface{}) Query {
	return &query{
		Query: session.Session.Query(cql, arguments...),
	}
}

type query struct {
	*gocql.Query
}

func (delegate *query) Consistency(consistency gocql.Consistency) Query {
	return &query{
		Query: delegate.Query.Consistency(consistency),
	}
}

func (delegate *query) Iter() Iter {
	return delegate.Query.Iter()
}

func (delegate *query) PageSize(size int) Query {
	return &query{
		Query: delegate.Query.PageSize(size),
	}
}

func (delegate *query) delegate() *gocql.Query {
	return delegate.Query
}
