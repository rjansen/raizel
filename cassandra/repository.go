package cassandra

import (
	"context"

	"github.com/rjansen/raizel"
	"github.com/scylladb/gocqlx/qb"
)

type repository struct {
	session Session
}

func NewRepository(session Session) *repository {
	return &repository{session: session}
}

func (r *repository) Get(tree context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		cql, _ = qb.Select(key.EntityName()).Columns(
			"id", "name", "age", "created_at", "updated_at",
		).Where(
			qb.Eq(key.Name()),
		).ToCql()
		query = r.session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Scan(entity)
}

func (r *repository) Set(tree context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		cql, _ = qb.Insert(key.EntityName()).Columns(
			"id", "name", "age", "created_at", "updated_at",
		).ToCql()
		query = r.session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Exec()
}

func (r *repository) Delete(tree context.Context, key raizel.EntityKey) error {
	var (
		cql, _ = qb.Delete(key.EntityName()).Where(
			qb.Eq(key.Name()),
		).ToCql()
		query = r.session.Query(cql)
		// queryx = gocqlx.Query(query.delegate(), args).BindStruct(entity)
	)
	return query.Exec()
}

func (r *repository) Close(tree context.Context) error {
	r.session.Close()
	return nil
}
