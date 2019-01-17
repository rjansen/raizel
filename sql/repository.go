package sql

import (
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
)

type repository struct {
	mapper Mapper
}

func NewRepository(mapper Mapper) repository {
	return repository{mapper: mapper}
}

func (repository repository) Get(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		db        = MustReference(tree)
		sqlStruct = repository.mapper.Get(key.EntityName())
		builder   = sqlStruct.SelectFrom(key.EntityName())
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		row = db.QueryRow(sql, args...)
	)
	return row.Scan(sqlStruct.Addr(entity))
}

func (repository repository) Set(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		db        = MustReference(tree)
		sqlStruct = repository.mapper.Get(key.EntityName())
		builder   = sqlStruct.Update(key.EntityName(), entity)
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		_, err = db.Exec(sql, args...)
	)
	if err != nil {
		return err
	}
	return nil
}

func (repository repository) Delete(tree yggdrasil.Tree, key raizel.EntityKey) error {
	var (
		db        = MustReference(tree)
		sqlStruct = repository.mapper.Get(key.EntityName())
		builder   = sqlStruct.DeleteFrom(key.EntityName())
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		_, err = db.Exec(sql, args...)
	)
	if err != nil {
		return err
	}
	return nil
}

func (repository) Close(tree yggdrasil.Tree) error {
	client := MustReference(tree)
	return client.Close()
}
