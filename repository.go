package raizel

import "github.com/rjansen/yggdrasil"

type EntityKey interface {
	GetKeyValue() interface{}
	GetEntityName() string
}

type Entity interface{}

type Repository interface {
	Get(yggdrasil.Tree, EntityKey, Entity) error
	Set(yggdrasil.Tree, EntityKey, Entity) error
	Delete(yggdrasil.Tree, EntityKey) error
	Close(yggdrasil.Tree) error
}
