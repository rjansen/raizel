package firestore

import (
	"context"
	"fmt"

	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
)

type repository struct{}

func NewRepository() raizel.Repository {
	return new(repository)
}

func entityDocRef(key raizel.EntityKey) string {
	return fmt.Sprintf("%s/%s", key.EntityName(), key.Value())
}

func (*repository) Get(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		client   = MustReference(tree)
		ref      = client.Doc(entityDocRef(key))
		doc, err = ref.Get(context.Background())
	)
	if err != nil {
		return err
	}
	return doc.DataTo(entity)
}

func (*repository) Set(tree yggdrasil.Tree, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		client = MustReference(tree)
		ref    = client.Doc(entityDocRef(key))
	)
	return ref.Set(context.Background(), entity)
}

func (*repository) Delete(tree yggdrasil.Tree, key raizel.EntityKey) error {
	var (
		client = MustReference(tree)
		ref    = client.Doc(entityDocRef(key))
	)
	return ref.Delete(context.Background())
}

func (*repository) Close(tree yggdrasil.Tree) error {
	var client = MustReference(tree)
	return client.Close()
}
