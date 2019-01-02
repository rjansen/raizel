// +build integration

package cassandra

import (
	"os"
	"testing"

	"github.com/rjansen/migi"
	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/assert"
)

var (
	pool              raizel.ClientPool
	persistenceClient raizel.Client
	setted            = false
	key1              = "8b06603b-9b0d-4e8c-8aae-10f988639fe6"
	expires           = 60
	testConfig        *Configuration
)

func init() {
	os.Args = append(os.Args, "-ecf", "etc/raizel/cassandra.yaml")
	if err := migi.UnmarshalKey("raizel.cassandra", &testConfig); err != nil {
		panic(err)
	}
}

func setup() error {
	if err := Setup(testConfig); err != nil {
		return err
	}
	var err error
	pool, err = raizel.GetPool()
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

func TestIntQueryOne(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from login where username = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan(nil, nil, nil, nil)
		}, "darkside")
	assert.Nil(t, err)
}

func TestIntQueryOneErr(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}

func TestIntQuery(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.Query("select * from login where username in (?, ?)",
		func(f raizel.Iterable) error {
			assert.NotNil(t, f)
			records := 0
			for f.Next() {
				fetchErr := f.Scan(nil, nil, nil, nil)
				assert.Nil(t, fetchErr)
				records++
			}
			assert.Equal(t, 2, records)
			return nil
		}, "rjansen", "darkside")
	assert.Nil(t, err)
}

func TestIntQueryErr(t *testing.T) {
	if beforeErr := before(); beforeErr != nil {
		assert.Fail(t, beforeErr.Error())
	}
	assert.NotNil(t, persistenceClient)
	err := persistenceClient.Query("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Iterable) error {
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)

}
