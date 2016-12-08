package cassandra

import (
	"farm.e-pedion.com/repo/persistence"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mysqlPool   persistence.ClientPool
	mysqlClient persistence.Client
)

func TestIntMySqlSetup(t *testing.T) {
	var err error
	mysqlConfig := &Configuration{
		Driver:   "mysql",
		URL:      "tcp(127.0.0.1:3306)/fivecolors",
		Database: "fivecolors",
		Username: "fivecolors",
		Password: "fivecolors",
		NumConns: 10,
	}
	err = Setup(mysqlConfig)
	assert.Nil(t, err)
	mysqlPool, err = persistence.GetPool()
	assert.Nil(t, err)
	assert.NotNil(t, mysqlPool)
	mysqlClient, err = mysqlPool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, mysqlClient)
}

func TestIntMySqlQuery(t *testing.T) {
	err := mysqlClient.QueryOne("select * from deck where id = ?",
		func(f persistence.Fetchable) error {
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

func TestIntMySqlQueryErr(t *testing.T) {
	err := mysqlClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f persistence.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}
