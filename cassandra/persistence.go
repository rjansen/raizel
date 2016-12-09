package cassandra

import (
	"errors"
	"farm.e-pedion.com/repo/logger"
	"farm.e-pedion.com/repo/persistence"
	"fmt"
	"github.com/gocql/gocql"

	// "github.com/matryer/resync"
	"time"
)

//pool is a variable to hold the Cassandra Pool
var (
	// once          resync.Once
	Config *Configuration
)

//Configuration holds Cassandra connections parameters
type Configuration struct {
	URL       string        `json:"url" mapstructure:"url"`
	Keyspace  string        `json:"keyspace" mapstructure:"keyspace"`
	Username  string        `json:"username" mapstructure:"username"`
	Password  string        `json:"password" mapstructure:"password"`
	NumConns  int           `json:"numConns" mapstructure:"numConns"`
	KeepAlive time.Duration `json:"keepAliveDuration" mapstructure:"keepAliveDuration"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("cassandra.Configuration URL=%v Keyspace=%v Username=%v Password=%v NumConns=%d KeepAlive=%s",
		c.URL, c.Keyspace, c.Username, c.Password, c.NumConns, c.KeepAlive,
	)
}

//Session is the interface interact with the cassandra
type Session interface {
	Query(string, ...interface{}) Query
	Closed() bool
	Close()
}

//Query is an interface to execute cql commands
type Query interface {
	// func (q *Query) Bind(v ...interface{}) *Query
	Consistency(gocql.Consistency) Query
	Exec() error
	Iter() Iter
	PageSize(n int) Query
	Release()
	Scan(dest ...interface{}) error
	//String() string
	//WithContext(ctx context.Context) Query
}

//Iter is an interface to read data sets from cassandra
type Iter interface {
	Close() error
	NumRows() int
	Scanner() persistence.Iterable
}

//Setup configures a poll for database connections
func Setup(cfg *Configuration) error {
	logger.Info("cassandra.ConfigCluster",
		logger.String("configuration", cfg.String()),
	)
	cluster := gocql.NewCluster(cfg.URL)
	cluster.NumConns = cfg.NumConns
	cluster.SocketKeepalive = cfg.KeepAlive
	cluster.ProtoVersion = 4
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return fmt.Errorf("cassandra.CreateSessionErr err=%v", err.Error())
	}
	pool := &Pool{
		cluster: cluster,
		session: NewDelegateSession(session),
	}
	if err = persistence.Setup(pool); err != nil {
		return fmt.Errorf("cassandra.SetupPersistenceErr err=%v", err.Error())
	}
	logger.Info("cassandra.DriverConfigured",
		logger.String("config", cfg.String()),
	)
	Config = cfg
	return nil
}

//Pool controls how new gocql.Session will create and maintained
type Pool struct {
	cluster *gocql.ClusterConfig
	session Session
}

func (c Pool) String() string {
	return fmt.Sprintf("CassandraPool Configuration=%s ClusterIsNil=%t SessionIsNil=%t",
		Config.String(),
		c.cluster == nil,
		c.session == nil,
	)
}

//Get creates and returns a Client reference
func (c *Pool) Get() (persistence.Client, error) {
	if c == nil || c.session == nil {
		return nil, errors.New("SetupMustCalled: Message='You must call Setup with a CassandraConfig before get a Cassandrapool reference')")
	}
	if c.session.Closed() {
		return nil, fmt.Errorf("cassandra.SessionIsClosedErr")
	}
	logger.Debug("cassandra.GetSession",
		logger.String("Pool", c.String()),
		logger.Bool("SessionIsNil", c.session == nil),
		logger.Bool("SessionIsClosed", c.session.Closed()),
	)
	return NewClient(c.session), nil
}

//Close close the database pool
func (c *Pool) Close() error {
	if c == nil || c.cluster == nil {
		return errors.New("SetupMustCalled: Message='You must call Setup with a CassandraBConfig before get a Cassandrapool reference')")
	}
	logger.Info("CloseCassandraSession",
		logger.String("CassandraPool", c.String()),
	)
	c.session.Close()
	return nil
}
