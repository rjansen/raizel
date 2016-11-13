package cassandra

import (
	"errors"
	"farm.e-pedion.com/repo/logger"
	"farm.e-pedion.com/repo/persistence"
	"fmt"
	"github.com/gocql/gocql"
	"strings"
)

var (
	NotFoundErr        = gocql.ErrNotFound
	CassandraClientKey = "cassandra.Client"
	ClientNotFoundErr  = errors.New("cassandra.ClientNotFoundErr message='Cassandra client does not found at context'")
)

type sessionObject struct {
	//session is a transient pointer to database connection
	session *gocql.Session
}

//SetSession attachs a database connection to Card
func (d *sessionObject) SetSession(session *gocql.Session) error {
	if session == nil {
		return errors.New("NullSessionReferenceErr: Message='The db parameter is required'")
	}
	d.session = session
	return nil
}

//GetSession returns the Card attached connection
func (d sessionObject) GetSession() (*gocql.Session, error) {
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
func (q *QuerySupport) QueryOne(query string, fetchFunc func(persistence.Fetchable) error, params ...interface{}) error {
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
func (q *QuerySupport) Query(query string, iterFunc func(persistence.Iterable) error, params ...interface{}) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("QueryError[Messages='EmptyCQLQuery']")
	}
	if params == nil || len(params) <= 0 {
		return errors.New("QueryError[Messages='EmptyQueryParameters']")
	}
	if iterFunc == nil {
		return errors.New("QueryError[Messages='NilIterFunc']")
	}
	cqlQuery := q.session.Query(query, params...)
	return iterFunc(cqlQuery)
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
		logger.Error("CQLExecutionFalied",
			logger.String("CQL", cql),
			logger.Struct("Parameters", params),
		)
		return err
	}
	logger.Debug("CQLExecutedSuccessfully",
		logger.String("CQL", cql),
		logger.Struct("Parameters", params),
	)
	return nil
}

//NewClient creates a new instance of the CQLClient
func NewClient(session *gocql.Session) *CQLClient {
	client := new(CQLClient)
	client.QuerySupport.SetSession(session)
	client.ExecSupport.SetSession(session)
	return client
}

//CQLClient adds full query and exec support fot the struct
type CQLClient struct {
	QuerySupport
	ExecSupport
}

func (c *CQLClient) Close() error {
	var errs []error
	errs = append(errs, c.QuerySupport.Close())
	errs = append(errs, c.ExecSupport.Close())
	if len(errs) > 0 {
		return fmt.Errorf("cassandra.Client.CloseErr msg='%v'", errs)
	}
	return nil
}
