package firestore

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid Client Reference")
	clientPath          = yggdrasil.NewPath("/raizel/firestore/client")
)

func Register(roots *yggdrasil.Roots, client Client) error {
	return roots.Register(clientPath, client)
}

func Reference(tree yggdrasil.Tree) (Client, error) {
	reference, err := tree.Reference(clientPath)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, nil
	}
	client, is := reference.(Client)
	if !is {
		return nil, ErrInvalidReference
	}
	return client, nil
}

func MustReference(tree yggdrasil.Tree) Client {
	client, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return client
}
