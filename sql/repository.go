package sql

import (
	sqlbuilder "github.com/huandu/go-sqlbuilder"
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
)

type repository struct{}

func NewRepository() *repository {
	return new(repository)
}

func (*repository) Get(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		db           = MustReference(tree)
		entityStruct = sqlbuilder.NewStruct(entity)
		builder      = entityStruct.SelectFrom("entity")
		sql, args    = builder.Where(
			builder.E(key.GetEntityName(), key.GetKeyValue()),
		).Build()
		row = db.QueryRow(sql, args...)
	)
	return row.Scan(entityStruct.Addr(entity))
}

func (*repository) Set(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		db           = MustReference(tree)
		entityStruct = sqlbuilder.NewStruct(entity)
		builder      = entityStruct.Update("entity", entity)
		sql, args    = builder.Where(
			builder.E(key.GetEntityName(), key.GetKeyValue()),
		).Build()
		_, err = db.Exec(sql, args...)
	)
	if err != nil {
		return err
	}
	return nil
}

func (*repository) Delete(tree yggdrasil.Tree, key raizel.EntityKey) error {
	var (
		db           = MustReference(tree)
		entityStruct = sqlbuilder.NewStruct(key)
		builder      = entityStruct.DeleteFrom("entity")
		sql, args    = builder.Where(
			builder.E(key.GetEntityName(), key.GetKeyValue()),
		).Build()
		_, err = db.Exec(sql, args...)
	)
	if err != nil {
		return err
	}
	return nil
}

func (*repository) Close(tree yggdrasil.Tree) error {
	var client = MustReference(tree)
	return client.Close()
}
