package cassandra

import (
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/scylladb/gocqlx/qb"
)

type repository struct{}

func NewRepository() *repository {
	return new(repository)
}

func (*repository) Get(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		session = MustReference(tree)
		cql, _  = qb.Select(key.GetEntityName()).Columns(
			"id", "name", "age", "created_at", "updated_at",
		).Where(
			qb.Eq(key.GetEntityName()),
		).ToCql()
		query = session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Scan(entity)
}

func (*repository) Set(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		session = MustReference(tree)
		cql, _  = qb.Insert(key.GetEntityName()).Columns(
			"id", "name", "age", "created_at", "updated_at",
		).ToCql()
		query = session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Exec()
}

func (*repository) Delete(tree yggdrasil.Tree, key raizel.EntityKey) error {
	var (
		session = MustReference(tree)
		cql, _  = qb.Delete(key.GetEntityName()).Where(
			qb.Eq(key.GetEntityName()),
		).ToCql()
		query = session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Exec()
}

func (*repository) Close(tree yggdrasil.Tree) error {
	session := MustReference(tree)
	session.Close()
	return nil
}
