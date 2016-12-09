package persistence

import (
	"context"
	"errors"
	"fmt"
)

var (
	persistenceClientKey = 40
	ErrInvalidState      = errors.New("The persistence current state is invalid. Setup never called")
	ErrInvalidClientPool = errors.New("The provided ClientPool is invalid")
	ErrInvalidContext    = errors.New("The provided Context is invalid")
	ErrInvalidConfig     = errors.New("The provided Configuration is invalid")
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

//Iterable supply the Iter interface for a struct
type Iterable interface {
	Fetchable
	Next() bool
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
	Close() error
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
	// Iter(iterable Iterable) error
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

//Setup initializes the persistence package
func Setup(p ClientPool) error {
	if p == nil {
		return ErrInvalidClientPool
	}
	pool = p
	return nil
}

//GetPool returns the pool instance
func GetPool() (ClientPool, error) {
	if pool == nil {
		return nil, ErrInvalidState
	}
	return pool, nil
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

//SetClient preapres and set a persistence Client into context
func SetClient(c context.Context) (context.Context, error) {
	if c == nil {
		return nil, ErrInvalidContext
	}
	persistenceClient, err := pool.Get()
	if err != nil {
		return nil, err
	}
	return context.WithValue(c, persistenceClientKey, persistenceClient), nil
}

//ContextFunc is a functions with context olny parameter
type ContextFunc func(context.Context) error

//ExecuteContext preapres a Client and set it inside context to call the provided function
func ExecuteContext(ctxFunc ContextFunc) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx, err := SetClient(ctx)
	if err != nil {
		return err
	}
	return ctxFunc(ctx)
}

//ClientFunc is a functions with context olny parameter
type ClientFunc func(Client) error

//Execute gets a Client from the ClientPool and calls the provided function with the Client instance
func Execute(cliFunc ClientFunc) error {
	persistenceClient, err := pool.Get()
	defer persistenceClient.Close()
	if err != nil {
		return err
	}
	return cliFunc(persistenceClient)
}
