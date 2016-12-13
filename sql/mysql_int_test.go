package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mysqlPool   raizel.ClientPool
	mysqlClient raizel.Client
)

func TestIntMySqlSetup(t *testing.T) {
	var err error
	mysqlConfig := &Configuration{
		Driver:   "mysql",
		URL:      "fivecolors:fivecolors@tcp(127.0.0.1:3306)/fivecolors",
		NumConns: 10,
	}
	err = Setup(mysqlConfig)
	assert.Nil(t, err)
	mysqlPool, err = raizel.GetPool()
	assert.Nil(t, err)
	assert.NotNil(t, mysqlPool)
	mysqlClient, err = mysqlPool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, mysqlClient)
}

func TestIntMySqlQueryOne(t *testing.T) {
	assert.NotNil(t, mysqlClient)
	err := mysqlClient.QueryOne("select * from deck where id = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			var id int
			var name string
			var idPlayer int
			err := f.Scan(&id, &name, &idPlayer)
			assert.Nil(t, err)
			assert.NotZero(t, id)
			assert.NotZero(t, name)
			assert.NotZero(t, idPlayer)
			return err
		}, 1)
	assert.Nil(t, err)
}

func TestIntMySqlQueryOneErr(t *testing.T) {
	assert.NotNil(t, mysqlClient)
	err := mysqlClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f raizel.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}

func TestIntMySqlQuery(t *testing.T) {
	assert.NotNil(t, mysqlClient)
	err := mysqlClient.Query("select * from deck where id > ?",
		func(f raizel.Iterable) error {
			assert.NotNil(t, f)
			var err error
			for f.Next() {
				var id int
				var name string
				var idPlayer int
				err = f.Scan(&id, &name, &idPlayer)
				assert.Nil(t, err)
				assert.NotZero(t, id)
				assert.NotZero(t, name)
				assert.NotZero(t, idPlayer)
			}
			return err
		}, 0)
	assert.Nil(t, err)
}

func TestIntMySqlQueryErr(t *testing.T) {
	assert.NotNil(t, mysqlClient)
	err := mysqlClient.Query("select * from cql.mock m",
		func(f raizel.Iterable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}
