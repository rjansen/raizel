package raizel

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("err_notfound")
)

type EntityKey interface {
	EntityName() string
	Value() interface{}
	Name() string
}

type Entity interface{}

type Repository interface {
	Get(ctx context.Context, key EntityKey, entity Entity) error
	Set(ctx context.Context, key EntityKey, entity Entity) error
	Delete(ctx context.Context, key EntityKey) error
	Close(ctx context.Context) error
}

type dynamicEntityKey struct {
	entityName string
	name       string
	value      interface{}
}

func (key dynamicEntityKey) EntityName() string {
	return key.entityName
}

func (key dynamicEntityKey) Name() string {
	return key.name
}

func (key dynamicEntityKey) Value() interface{} {
	return key.value
}

func NewDynamicKey(entityName, keyName string, keyValue interface{}) EntityKey {
	return dynamicEntityKey{
		entityName: entityName,
		name:       keyName,
		value:      keyValue,
	}
}
