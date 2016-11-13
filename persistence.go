package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
)

var (
	ErrInvalidClientPool = errors.New("The provided ClientPool is invalid")
	ErrInvalidContext    = errors.New("The provided Context is invalid")
	persistenceClientKey = 40
	pool                 ClientPool
)

//Configuration holds Cassandra connections parameters
type Configuration struct {
	Provider string `json:"provider" mapstructure:"provider"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("persistence.Configuration Provider=%v", c.Provider)
}

//Fetchable supply the gocql.Query.Scan interface for a struct
type Fetchable interface {
	Scan(dest ...interface{}) error
}

//Iterable supply the gocql.Query.Iter interface for a struct
type Iterable interface {
	Iter() *gocql.Iter
}

//Reader provides the interface for persistence read actions
type Reader interface {
	QueryOne(query string, fetchFunc func(Fetchable) error, params ...interface{}) error
	Query(query string, iterFunc func(Iterable) error, params ...interface{}) error
	Close() error
}

//Executor provides cassandra exec supports
type Executor interface {
	Exec(command string, params ...interface{}) error
}

//ClientPool is a contrant for a persistence Client pool
type ClientPool interface {
	Get() (Client, error)
}

//Client adds full persistence supports
type Client interface {
	Reader
	Executor
}

//Readable provides read actions for a struct
type Readable interface {
	Reader
	Fetch(fetchable Fetchable) error
	Iter(iterable Iterable) error
	Read() error
	//ReadExample() ([]Readable, error)
}

//Writable provides persistence actions for a struct
type Writable interface {
	Executor
	Create() error
	Update() error
	Delete() error
}

func Setup(p ClientPool) error {
	if p == nil {
		return ErrInvalidClientPool
	}
	pool = p
	return nil
}

//GetClient reads a dbClient from a context
func GetClient(c context.Context) (Client, error) {
	if c == nil {
		return nil, ErrInvalidContext
	}
	client, ok := c.Value(persistenceClientKey).(Client)
	if !ok {
		return nil, fmt.Errorf("persistence.ErrInvalidClient client=%+v", client)
	}
	return client, nil
}

//SetClient preapres and set a dbClient into context
func SetClient(c context.Context, persistenceClient Client) (context.Context, error) {
	if c == nil {
		return nil, ErrInvalidContext
	}
	return context.WithValue(c, persistenceClientKey, persistenceClient), nil
}

//ContextFunc is a functions with context olny parameter
type ContextFunc func(context.Context) error

//Execute preapres a dbClient and set it inside context to call the provided function
func Execute(ctxFunc ContextFunc) error {
	var err error
	client, err := pool.Get()
	defer client.Close()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := SetClient(ctx, client)
	if err != nil {
		return err
	}
	return ctxFunc(c)
}
