package raizel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type entityKeyMock struct{}

func (entityKeyMock) EntityName() string { return "" }
func (entityKeyMock) Name() string       { return "" }
func (entityKeyMock) Value() interface{} { return nil }

func TestEntityKey(test *testing.T) {
	var (
		key EntityKey = entityKeyMock{}
	)
	require.Implements(test, (*EntityKey)(nil), key, "invalid entitykey type")
	_ = key.EntityName()
	_ = key.Name()
	_ = key.Value()
}

func TestDynamicEntityKey(test *testing.T) {
	var (
		entityName = "entity_name"
		keyName    = "key_name"
		keyValue   = "key_value"
	)
	key := NewDynamicKey(entityName, keyName, keyValue)
	require.NotNil(test, key, "key invalid")
	require.Equal(test, entityName, key.EntityName(), "entityname invalid instance")
	require.Equal(test, keyName, key.Name(), "keyname invalid instance")
	require.Equal(test, keyValue, key.Value(), "keyvalue invalid instance")
}

type repositoryMock struct{}

func (repositoryMock) Get(context.Context, EntityKey, Entity) error { return nil }
func (repositoryMock) Set(context.Context, EntityKey, Entity) error { return nil }
func (repositoryMock) Delete(context.Context, EntityKey) error      { return nil }
func (repositoryMock) Close(context.Context) error                  { return nil }

type repositoryTest struct {
	ctx    context.Context
	key    EntityKey
	result Entity
	entity Entity
}

func TestRepository(test *testing.T) {
	var (
		repository Repository = repositoryMock{}
		scenario              = repositoryTest{}
	)
	require.Implements(test, (*Repository)(nil), repository, "invalid repository type")
	_ = repository.Get(scenario.ctx, scenario.key, &scenario.result)
	_ = repository.Set(scenario.ctx, scenario.key, &scenario.entity)
	_ = repository.Delete(scenario.ctx, scenario.key)
	_ = repository.Close(scenario.ctx)
}

func TestErrNoyFound(test *testing.T) {
	require.NotNil(test, ErrNotFound, "errnotfound invalid instance")
}
