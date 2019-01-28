package raizel

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Repository Reference")
	repositoryPath      = yggdrasil.NewPath("/raizel/repository")
)

func Register(roots *yggdrasil.Roots, repository Repository) error {
	return roots.Register(repositoryPath, repository)
}

func Reference(tree yggdrasil.Tree) (Repository, error) {
	reference, err := tree.Reference(repositoryPath)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, nil
	}
	client, is := reference.(Repository)
	if !is {
		return nil, ErrInvalidReference
	}
	return client, nil
}

func MustReference(tree yggdrasil.Tree) Repository {
	client, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return client
}
