package mock

import (
	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

func (mock *RepositoryMock) Get(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	args := mock.Called(tree, key, entity)
	return args.Error(0)
}

func (mock *RepositoryMock) Set(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	args := mock.Called(tree, key, entity)
	return args.Error(0)
}

func (mock *RepositoryMock) Delete(tree yggdrasil.Tree, key raizel.EntityKey) error {
	args := mock.Called(tree, key)
	return args.Error(0)
}

func (mock *RepositoryMock) Close(tree yggdrasil.Tree) error {
	args := mock.Called(tree)
	return args.Error(0)
}
