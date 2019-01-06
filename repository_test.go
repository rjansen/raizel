package raizel

import (
	"testing"

	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

type repositoryTest struct{}

func (repositoryTest) Get(yggdrasil.Tree, EntityKey, Entity) error { return nil }
func (repositoryTest) Set(yggdrasil.Tree, EntityKey, Entity) error { return nil }
func (repositoryTest) Delete(yggdrasil.Tree, EntityKey) error      { return nil }
func (repositoryTest) Close(yggdrasil.Tree) error                  { return nil }

type repositoryScenarioTest struct {
	tree   yggdrasil.Tree
	key    EntityKey
	result Entity
	entity Entity
}

func TestRepository(test *testing.T) {
	var (
		repository Repository = repositoryTest{}
		scenario              = repositoryScenarioTest{}
	)
	require.Implements(test, (*Repository)(nil), repository, "invalid repository type")
	_ = repository.Get(scenario.tree, scenario.key, &scenario.result)
	_ = repository.Set(scenario.tree, scenario.key, &scenario.entity)
	_ = repository.Delete(scenario.tree, scenario.key)
	_ = repository.Close(scenario.tree)
}
