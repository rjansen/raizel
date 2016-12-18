package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"
	"strings"
)

var (
	NotFoundErr       = sql.ErrNoRows
	SqlClientKey      = "sql.Client"
	ClientNotFoundErr = errors.New("sql.ClientNotFoundErr message='Sql client does not found at context'")
)

func NewDelegateDB(d *sql.DB) SqlDB {
	return &DelegateDB{
		db: d,
	}
}

type DelegateDB struct {
	db *sql.DB
}

func (d *DelegateDB) Exec(sql string, params ...interface{}) (sql.Result, error) {
	return d.db.Exec(sql, params...)
}

func (d *DelegateDB) QueryRow(sql string, params ...interface{}) Row {
	row := d.db.QueryRow(sql, params...)
	return NewDelegateRow(row)
}

func (d *DelegateDB) Query(sql string, params ...interface{}) (Rows, error) {
	rows, err := d.db.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	return NewDelegateRows(rows), nil
}

func (d *DelegateDB) Ping() error {
	return d.db.Ping()
}

func (d *DelegateDB) Close() error {
	return d.db.Close()
}

func NewDelegateRow(row *sql.Row) Row {
	return &DelegateRow{
		row: row,
	}
}

type DelegateRow struct {
	row *sql.Row
}

func (d *DelegateRow) Scan(dest ...interface{}) error {
	return d.row.Scan(dest...)
}

func NewDelegateRows(rows *sql.Rows) Rows {
	return &DelegateRows{
		rows: rows,
	}
}

type DelegateRows struct {
	rows *sql.Rows
}

func (d *DelegateRows) Next() bool {
	return d.rows.Next()
}

func (d *DelegateRows) Scan(dest ...interface{}) error {
	return d.rows.Scan(dest...)
}

type dbObject struct {
	//db is a transient pointer to database connection
	db SqlDB
}

//SetSession attachs a database connection to Card
func (d *dbObject) SetDB(db SqlDB) error {
	if db == nil {
		return errors.New("NullSqlDBReferenceErr: Message='The db parameter is required'")
	}
	d.db = db
	return nil
}

//GetDB returns the attached database connection
func (d dbObject) GetDB() (SqlDB, error) {
	if d.db == nil {
		return nil, errors.New("NotAttachedErr: Message='The sql connection is null'")
	}
	return d.db, nil
}

func (d *dbObject) Close() error {
	// return d.db.Close()
	return nil
}

//QuerySupport adds query capability to the struct
type QuerySupport struct {
	dbObject
}

//QueryOne executes the sql query with the provided parameters and send the single result to the provided fetch function
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
	row := q.db.QueryRow(query, params...)
	err := fetchFunc(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return raizel.ErrNotFound
		}
		return err
	}
	return nil
}

//Query executes the sql query with the provided parameters and send the results to the provided iter function
func (q *QuerySupport) Query(query string, iterFunc func(raizel.Iterable) error, params ...interface{}) error {
	if strings.TrimSpace(query) == "" {
		return errors.New("QueryError[Messages='EmptyCQLQuery']")
	}
	if iterFunc == nil {
		return errors.New("QueryError[Messages='NilIterFunc']")
	}
	rows, err := q.db.Query(query, params...)
	if err != nil {
		return err
	}
	return iterFunc(rows)
}

//ExecSupport adds cql exec capability to the struct
type ExecSupport struct {
	dbObject
}

//Exec exeutes the sql command with the provided parameters
func (i *ExecSupport) Exec(sql string, params ...interface{}) (raizel.Result, error) {
	if strings.TrimSpace(sql) == "" {
		return nil, errors.New("sql.ExecError[Messages='NilSQLQuery']")
	}
	if params == nil || len(params) <= 0 {
		return nil, errors.New("ExecParametersLenInvalid[Messages='EmptyExecParameters']")
	}
	result, err := i.db.Exec(sql, params...)
	if err != nil {
		l.Error("sql.ExecutionFalied",
			l.String("SQL", sql),
			l.Struct("Parameters", params),
		)
		return nil, err
	}
	l.Debug("SQLExecutedSuccessfully",
		l.String("SQL", sql),
		l.Struct("Parameters", params),
		l.Struct("Result", result),
	)
	return result, nil
}

//NewClient creates a new instance of the SQLClient
func NewClient(db SqlDB) *Client {
	client := new(Client)
	client.QuerySupport.SetDB(db)
	client.ExecSupport.SetDB(db)
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
		return fmt.Errorf("sql.Client.CloseErr msg='%v'", errs)
	}
	return nil
}
