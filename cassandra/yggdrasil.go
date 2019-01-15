package cassandra

import (
	"errors"

	"github.com/rjansen/yggdrasil"
)

var (
	ErrInvalidReference = errors.New("Invalid DB Reference")
	sessionPath         = yggdrasil.NewPath("/raizel/cassandra/session")
)

func Register(roots *yggdrasil.Roots, session Session) error {
	return roots.Register(sessionPath, session)
}

func Reference(tree yggdrasil.Tree) (Session, error) {
	reference, err := tree.Reference(sessionPath)
	if err != nil {
		return nil, err
	}
	if reference == nil {
		return nil, nil
	}
	session, is := reference.(Session)
	if !is {
		return nil, ErrInvalidReference
	}
	return session, nil
}

func MustReference(tree yggdrasil.Tree) Session {
	session, err := Reference(tree)
	if err != nil {
		panic(err)
	}
	return session
}
