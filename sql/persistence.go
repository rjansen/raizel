package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/rjansen/l"
	"github.com/rjansen/raizel"

	"time"
)

var (
	Config *Configuration
)

//Configuration holds Cassandra connections parameters
type Configuration struct {
	Driver    string        `json:"driver" mapstructure:"driver"`
	URL       string        `json:"url" mapstructure:"url"`
	Database  string        `json:"database" mapstructure:"database"`
	Username  string        `json:"username" mapstructure:"username"`
	Password  string        `json:"password" mapstructure:"password"`
	NumConns  int           `json:"numConns" mapstructure:"numConns"`
	KeepAlive time.Duration `json:"keepAliveDuration" mapstructure:"keepAliveDuration"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("sql.Configuration Driver=%s URL=%s Database=%s Username=%s Password=%s NumConns=%d KeepAlive=%s",
		c.Driver, c.URL, c.Database, c.Username, c.Password, c.NumConns, c.KeepAlive,
	)
}

type SqlDB interface {
	//func (db *DB) Begin() (*Tx, error)
	Close() error
	//func (db *DB) Driver() driver.Driver
	Exec(string, ...interface{}) (sql.Result, error)
	Ping() error
	//func (db *DB) Prepare(query string) (*Stmt, error)
	Query(string, ...interface{}) (Rows, error)
	QueryRow(string, ...interface{}) Row
	//func (db *DB) SetConnMaxLifetime(d time.Duration)
	//func (db *DB) SetMaxIdleConns(n int)
	//func (db *DB) SetMaxOpenConns(n int)
	//func (db *DB) Stats() DBStats
}

type Rows interface {
	raizel.Iterable
}

type Row interface {
	raizel.Fetchable
}

//Setup configures a poll for database connections
func Setup(cfg *Configuration) error {
	l.Info("sql.ConfigDatasource",
		l.String("configuration", cfg.String()),
	)
	datasource := &Datasource{
		Driver:   cfg.Driver,
		Username: cfg.Username,
		Password: cfg.Password,
		URL:      cfg.URL,
	}
	db, err := sql.Open(datasource.Driver, datasource.DSN())
	if err != nil {
		return fmt.Errorf("sql.OpenErr err=%v", err.Error())
	}
	pool := &Pool{
		datasource: datasource,
		db:         NewDelegateDB(db),
	}
	if err = raizel.Setup(pool); err != nil {
		return fmt.Errorf("sql.SetupPersistenceErr err=%v", err.Error())
	}
	l.Info("sql.DriverConfigured",
		l.String("driver", cfg.URL),
		l.String("url", cfg.Driver),
	)
	Config = cfg
	return nil
}

//Pool controls how new sql.DB will create and maintained
type Pool struct {
	datasource *Datasource
	db         SqlDB
}

func (c Pool) String() string {
	return fmt.Sprintf("sql.Pool Configuration=%s DatasourceIsNil=%t DBIsNil=%t",
		Config.String(),
		c.datasource == nil,
		c.db == nil,
	)
}

//Get creates and returns a Client reference
func (c *Pool) Get() (raizel.Client, error) {
	if c == nil || c.db == nil {
		return nil, errors.New("SetupMustCalled: Message='You must call Setup with a SqlConfig before get a SqlPool reference')")
	}
	if err := c.db.Ping(); err != nil {
		return nil, err
	}
	l.Debug("sql.Get",
		l.String("Pool", c.String()),
	)
	return NewClient(c.db), nil
}

//Close close the database pool
func (c *Pool) Close() error {
	if c == nil || c.db == nil {
		return errors.New("SetupMustCalled: Message='You must call Setup with a SqlConfig before get a SqlPool reference')")
	}
	l.Info("sql.CloseDB",
		l.String("DBPool", c.String()),
	)
	return c.db.Close()
}

//Datasource holds parameterts to create new sql.DB connections
type Datasource struct {
	Driver   string
	URL      string
	Username string
	Password string
}

//DSN retuns a DSN representation of Datasource struct
//DSN format: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (d *Datasource) DSN() string {
	//TODO: Fix this
	if d.Driver == "postgres" {
		return fmt.Sprintf("postgres://%s:%s@%s", d.Username, d.Password, d.URL)
	}
	return fmt.Sprintf("%s:%s@%s", d.Username, d.Password, d.URL)

}

//FromDSN fills the connection parameters of this Datasource instance
// func (d *Datasource) FromDSN(DSN string) error {
// 	regex := "(()?(:())@)?()?/()?"
// 	return fmt.Errorf("NotImplemented: Regex='%v'", regex)
// }
