package spanner

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
)

type Client interface {
	Apply(context.Context, []*Mutation, ...ApplyOption) (time.Time, error)
	BatchReadOnlyTransaction(context.Context, TimestampBound) (BatchReadOnlyTransaction, error)
	BatchReadOnlyTransactionFromID(BatchReadOnlyTransactionID) BatchReadOnlyTransaction
	Close()
	PartitionedUpdate(context.Context, Statement) (int64, error)
	ReadOnlyTransaction() ReadOnlyTransaction
	ReadWriteTransaction(context.Context, func(context.Context, *ReadWriteTransaction) error) (time.Time, error)
	Single() ReadOnlyTransaction
}

type client struct {
	*spanner.Client
}

func NewClient(c *spanner.Client) Client {
	return &client{Client: c}
}

func (c *client) BatchReadOnlyTransaction(ctx context.Context, tb TimestampBound) (BatchReadOnlyTransaction, error) {
	transaction, err := c.Client.BatchReadOnlyTransaction(ctx, tb)
	if err != nil {
		return nil, err
	}
	return newBatchReadOnlyTransaction(transaction), nil
}

func (c *client) BatchReadOnlyTransactionFromID(tid BatchReadOnlyTransactionID) BatchReadOnlyTransaction {
	return newBatchReadOnlyTransaction(c.Client.BatchReadOnlyTransactionFromID(tid))
}

func (c *client) ReadOnlyTransaction() ReadOnlyTransaction {
	return newReadOnlyTransaction(c.Client.ReadOnlyTransaction())
}

func (c *client) Single() ReadOnlyTransaction {
	return newReadOnlyTransaction(c.Client.Single())
}
