package firestore

import (
	"context"
	"fmt"

	"github.com/rjansen/raizel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type repository struct {
	client Client
}

func NewRepository(client Client) raizel.Repository {
	return &repository{client: client}
}

func entityDocRef(key raizel.EntityKey) string {
	return fmt.Sprintf("%s/%s", key.EntityName(), key.Value())
}

func (r *repository) Get(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		ref      = r.client.Doc(entityDocRef(key))
		doc, err = ref.Get(context.Background())
	)
	if err != nil {
		if grpc.Code(err) == codes.NotFound {
			return raizel.ErrNotFound
		}
		return err
	}
	return doc.DataTo(entity)
}

func (r *repository) Set(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	var (
		ref = r.client.Doc(entityDocRef(key))
	)
	return ref.Set(context.Background(), entity)
}

func (r *repository) Delete(ctx context.Context, key raizel.EntityKey) error {
	var (
		ref = r.client.Doc(entityDocRef(key))
	)
	return ref.Delete(context.Background())
}

func (r *repository) Close(ctx context.Context) error {
	return r.client.Close()
}
