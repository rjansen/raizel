package cassandra

import (
	"farm.e-pedion.com/repo/config"
	"farm.e-pedion.com/repo/persistence"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	pool              persistence.ClientPool
	persistenceClient persistence.Client
	setted            = false
	key1              = "8b06603b-9b0d-4e8c-8aae-10f988639fe6"
	expires           = 60
	testConfig        *Configuration
)

func init() {
	os.Args = append(os.Args, "-ecf", "etc/persistence/persistence.yaml")
	if err := config.UnmarshalKey("persistence.cassandra", &testConfig); err != nil {
		panic(err)
	}
}

func setup() error {
	if err := Setup(testConfig); err != nil {
		return err
	}
	var err error
	pool, err = persistence.GetPool()
	return err
}

func before() error {
	if !setted {
		if err := setup(); err != nil {
			return err
		}
	}
	var err error
	persistenceClient, err = pool.Get()
	return err
}

func TestIntQuery(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from login where username = ?",
		func(f persistence.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan(nil, nil, nil, nil)
		}, "darkside")
	assert.Nil(t, err)
}

func TestIntQueryErr(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f persistence.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}
