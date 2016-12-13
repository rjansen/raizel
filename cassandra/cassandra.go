package cassandra

import (
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"strings"
)

var (
	NotFoundErr        = gocql.ErrNotFound
	CassandraClientKey = "cassandra.Client"
	ClientNotFoundErr  = errors.New("cassandra.ClientNotFoundErr message='Cassandra client does not found at context'")
)

func NewDelegateSession(d *gocql.Session) Session {
	return &DelegateSession{
		session: d,
	}
}

type DelegateSession struct {
	session *gocql.Session
}

func (d *DelegateSession) Query(cql string, params ...interface{}) Query {
	cqlQuery := d.session.Query(cql, params...)
	return NewDelegateQuery(cqlQuery)
}

func (d *DelegateSession) Closed() bool {
	return d.session.Closed()
}

func (d *DelegateSession) Close() {
	d.session.Close()
}

func NewDelegateQuery(d *gocql.Query) Query {
	return &DelegateQuery{
		query: d,
	}
}

type DelegateQuery struct {
	query *gocql.Query
}

func (d *DelegateQuery) Consistency(c gocql.Consistency) Query {
	d.query = d.query.Consistency(c)
	return d
}

func (d *DelegateQuery) Exec() error {
	return d.query.Exec()
}

func (d *DelegateQuery) Iter() Iter {
	cqlIter := d.query.Iter()
	return NewDelegateIter(cqlIter)
}

func (d *DelegateQuery) PageSize(n int) Query {
	d.query = d.query.PageSize(n)
	return d
}

func (d *DelegateQuery) Release() {
	d.query.Release()
}

func (d *DelegateQuery) Scan(dest ...interface{}) error {
	return d.query.Scan(dest...)
}

func NewDelegateIter(d *gocql.Iter) Iter {
	return &DelegateIter{
		iter: d,
	}
}

type DelegateIter struct {
	iter *gocql.Iter
}

func (d *DelegateIter) Close() error {
	return d.iter.Close()
}

func (d *DelegateIter) NumRows() int {
	return d.iter.NumRows()
}

func (d *DelegateIter) Scanner() raizel.Iterable {
	return d.iter.Scanner()
}

type sessionObject struct {
	//session is a transient pointer to database connection
	session Session
}

//SetSession attachs a database connection to Card
func (d *sessionObject) SetSession(session Session) error {
	if session == nil {
		return errors.New("NullSessionReferenceErr: Message='The db parameter is required'")
	}
	d.session = session
	return nil
}

//GetSession returns the Card attached connection
func (d sessionObject) GetSession() (Session, error) {
	if d.session == nil {
		return nil, errors.New("NotAttachedErr: Message='The cassandra session is null'")
	}
	return d.session, nil
}

func (d *sessionObject) Close() error {
	// return d.SetSession(nil)
	return nil
}

//QuerySupport adds query capability to the struct
type QuerySupport struct {
	sessionObject
}

//QueryOne executes the single result cql query with the provided parameters and fetch the result
func (q *QuerySupport) QueryOne(query string, fetchFunc func(raizel.Fetchable) error, params ...interface{}) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("identity.QuerySupport.QueryError: Messages='NilReadQuery")
	}
	if params == nil || len(params) <= 0 {
		return errors.New("identity.QuerySupport.QueryError: Messages='EmptyReadParameters")
	}
	if fetchFunc == nil {
		return errors.New("identity.QuerySupport.QueryError: Messages='NilFetchFunction")
	}
	cqlQuery := q.session.Query(query, params...).Consistency(gocql.One)
	return fetchFunc(cqlQuery)
}

//Query executes the cql query with the provided parameters and process the results
func (q *QuerySupport) Query(query string, iterFunc func(raizel.Iterable) error, params ...interface{}) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("QueryError[Messages='EmptyCQLQuery']")
	}
	if params == nil || len(params) <= 0 {
		return errors.New("QueryError[Messages='EmptyQueryParameters']")
	}
	if iterFunc == nil {
		return errors.New("QueryError[Messages='NilIterFunc']")
	}
	queryIter := q.session.Query(query, params...).Consistency(gocql.All).Iter()
	defer queryIter.Close()
	return iterFunc(queryIter.Scanner())
}

//ExecSupport adds cql exec capability to the struct
type ExecSupport struct {
	sessionObject
}

//Exec exeutes the command with the provided parameters
func (i *ExecSupport) Exec(cql string, params ...interface{}) error {
	if strings.TrimSpace(cql) == "" {
		return errors.New("ExecError[Messages='NilCQLQuery']")
	}
	if params == nil || len(params) <= 0 {
		return errors.New("ExecParametersLenInvalid[Messages='EmptyExecParameters']")
	}
	err := i.session.Query(cql, params...).Exec()
	if err != nil {
		l.Error("CQLExecutionFalied",
			l.String("CQL", cql),
			l.Struct("Parameters", params),
		)
		return err
	}
	l.Debug("CQLExecutedSuccessfully",
		l.String("CQL", cql),
		l.Struct("Parameters", params),
	)
	return nil
}

//NewClient creates a new instance of the CQLClient
func NewClient(session Session) *Client {
	client := new(Client)
	client.QuerySupport.SetSession(session)
	client.ExecSupport.SetSession(session)
	return client
}

//Client adds full query and exec support fot the struct
type Client struct {
	QuerySupport
	ExecSupport
}

func (c *Client) Close() error {
	var errs []error
	if err := c.QuerySupport.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := c.ExecSupport.Close(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("cassandra.Client.CloseErr msg='%v'", errs)
	}
	return nil
}
