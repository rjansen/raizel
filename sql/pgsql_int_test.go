package sql

import (
	"farm.e-pedion.com/repo/persistence"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	pgsqlPool   persistence.ClientPool
	pgsqlClient persistence.Client
)

func TestIntPgSqlSetup(t *testing.T) {
	var err error
	pgsqlConfig := &Configuration{
		Driver:   "postgres",
		URL:      "127.0.0.1:5432/fivecolors?sslmode=disable",
		Database: "fivecolors",
		Username: "fivecolors",
		Password: "fivecolors",
		NumConns: 10,
	}
	err = Setup(pgsqlConfig)
	assert.Nil(t, err)
	pgsqlPool, err = persistence.GetPool()
	assert.Nil(t, err)
	assert.NotNil(t, pgsqlPool)
	pgsqlClient, err = pgsqlPool.Get()
	assert.Nil(t, err)
	assert.NotNil(t, pgsqlClient)
}

func TestIntPgSqlQueryOne(t *testing.T) {
	assert.NotNil(t, pgsqlClient)
	err := pgsqlClient.QueryOne("select * from rarity where id = $1",
		func(f persistence.Fetchable) error {
			assert.NotNil(t, f)
			var id int
			var name string
			var alias string
			err := f.Scan(&id, &name, &alias)
			assert.Nil(t, err)
			assert.NotZero(t, id)
			assert.NotZero(t, name)
			assert.NotZero(t, alias)
			return err
		}, 1)
	assert.Nil(t, err)
}

func TestIntPgSqlQueryOneErr(t *testing.T) {
	assert.NotNil(t, pgsqlClient)
	err := pgsqlClient.QueryOne("select * from cql.mock m where m.mockField = ?",
		func(f persistence.Fetchable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, "mockValue")
	assert.NotNil(t, err)
}

func TestIntPgSqlQuery(t *testing.T) {
	assert.NotNil(t, pgsqlClient)
	err := pgsqlClient.Query("select * from rarity where id > $1",
		func(f persistence.Iterable) error {
			assert.NotNil(t, f)
			var id int
			var name string
			var alias string
			for f.Next() {
				err := f.Scan(&id, &name, &alias)
				assert.Nil(t, err)
				assert.NotZero(t, id)
				assert.NotZero(t, name)
				assert.NotZero(t, alias)
			}
			return nil
		}, 0)
	assert.Nil(t, err)
}

func TestIntPgSqlQueryErr(t *testing.T) {
	assert.NotNil(t, pgsqlClient)
	err := pgsqlClient.Query("select * from cql.mock m where m.mockField > ?",
		func(f persistence.Iterable) error {
			assert.NotNil(t, f)
			return f.Scan()
		}, 0)
	assert.NotNil(t, err)
}
