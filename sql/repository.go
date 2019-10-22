package sql

import (
	"context"
	database "database/sql"

	"github.com/lib/pq"
	"github.com/rjansen/raizel"
)

type repository struct {
	db     DB
	mapper Mapper
}

func NewRepository(db DB, mapper Mapper) repository {
	return repository{db: db, mapper: mapper}
}

func (repository repository) Get(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		sqlStruct = repository.mapper.Get(key.EntityName())
		builder   = sqlStruct.SelectFrom(key.EntityName())
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		row = repository.db.QueryRow(sql, args...)
		err = row.Scan(sqlStruct.Addr(entity)...)
	)
	if err != nil {
		if err == database.ErrNoRows {
			return raizel.ErrNotFound
		}
		return err
	}
	return nil
}

func (repository repository) Set(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		sqlStruct = repository.mapper.Get(key.EntityName())
		sql, args = sqlStruct.InsertInto(key.EntityName(), entity).Build()
		_, err    = repository.db.Exec(sql, args...)
	)
	if err != nil {
		pgerr, ispgerr := err.(*pq.Error)
		if !ispgerr {
			return err
		}

		// unique contraint validation code
		/*
			// Class 23 - Integrity Constraint Violation
			"23000": "integrity_constraint_violation",
			"23001": "restrict_violation",
			"23502": "not_null_violation",
			"23503": "foreign_key_violation",
			"23505": "unique_violation",
			"23514": "check_violation",
			"23P01": "exclusion_violation",
		*/
		if pgerr.Code != "23505" {
			return err
		}

		builder := sqlStruct.Update(key.EntityName(), entity)
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		_, err = repository.db.Exec(sql, args...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repository repository) Delete(ctx context.Context, key raizel.EntityKey) error {
	var (
		sqlStruct = repository.mapper.Get(key.EntityName())
		builder   = sqlStruct.DeleteFrom(key.EntityName())
		sql, args = builder.Where(
			builder.E(key.Name(), key.Value()),
		).Build()
		_, err = repository.db.Exec(sql, args...)
	)
	if err != nil {
		return err
	}
	return nil
}

func (r repository) Close(ctx context.Context) error {
	return r.db.Close()
}
