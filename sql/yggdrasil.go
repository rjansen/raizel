package sql

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid DB Reference")
	dbPath              = yggdrasil.NewPath("/raizel/sql/db")
)

func Register(roots *yggdrasil.Roots, db DB) error {
	return roots.Register(dbPath, db)
}

func Reference(tree yggdrasil.Tree) (DB, error) {
	reference, err := tree.Reference(dbPath)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, nil
	}
	db, is := reference.(DB)
	if !is {
		return nil, ErrInvalidReference
	}
	return db, nil
}

func MustReference(tree yggdrasil.Tree) DB {
	db, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return db
}
