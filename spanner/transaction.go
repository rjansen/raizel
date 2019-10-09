package spanner

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
)

type ReadOnlyTransaction interface {
	AnalyzeQuery(context.Context, Statement) (*QueryPlan, error)
	Close()
	Query(context.Context, Statement) RowIterator
	QueryWithStats(context.Context, Statement) RowIterator
	Read(context.Context, string, KeySet, []string) RowIterator
	ReadRow(context.Context, string, Key, []string) (Row, error)
	ReadUsingIndex(context.Context, string, string, KeySet, []string) RowIterator
	ReadWithOptions(context.Context, string, KeySet, []string, *ReadOptions) RowIterator
	Timestamp() (time.Time, error)
	WithTimestampBound(TimestampBound) ReadOnlyTransaction
}

type BatchReadOnlyTransaction interface {
	GetID() BatchReadOnlyTransactionID
	AnalyzeQuery(context.Context, Statement) (*QueryPlan, error)
	Cleanup(context.Context)
	Close()
	Execute(context.Context, *Partition) RowIterator
	PartitionQuery(context.Context, Statement, PartitionOptions) ([]*Partition, error)
	PartitionRead(context.Context, string, KeySet, []string, PartitionOptions) ([]*Partition, error)
	PartitionReadUsingIndex(context.Context, string, string, KeySet, []string, PartitionOptions) ([]*Partition, error)
	Query(context.Context, Statement) RowIterator
	QueryWithStats(context.Context, Statement) RowIterator
	Read(context.Context, string, KeySet, []string) RowIterator
	ReadRow(context.Context, string, Key, []string) (Row, error)
	ReadUsingIndex(context.Context, string, string, KeySet, []string) RowIterator
	ReadWithOptions(context.Context, string, KeySet, []string, *ReadOptions) RowIterator
}

type readOnlyTransaction struct {
	*spanner.ReadOnlyTransaction
}

func newReadOnlyTransaction(transaction *spanner.ReadOnlyTransaction) ReadOnlyTransaction {
	return &readOnlyTransaction{
		ReadOnlyTransaction: transaction,
	}
}

func (t *readOnlyTransaction) Query(ctx context.Context, stm Statement) RowIterator {
	return newRowIterator(t.ReadOnlyTransaction.Query(ctx, stm))
}

func (t *readOnlyTransaction) QueryWithStats(ctx context.Context, stm Statement) RowIterator {
	return newRowIterator(t.ReadOnlyTransaction.QueryWithStats(ctx, stm))
}

func (t *readOnlyTransaction) Read(ctx context.Context, table string, keys KeySet, cols []string) RowIterator {
	return newRowIterator(t.ReadOnlyTransaction.Read(ctx, table, keys, cols))
}

func (t *readOnlyTransaction) ReadRow(ctx context.Context, table string, key Key, cols []string) (Row, error) {
	row, err := t.ReadOnlyTransaction.ReadRow(ctx, table, key, cols)
	if err != nil {
		return nil, err
	}
	return newRow(row), nil
}

func (t *readOnlyTransaction) ReadUsingIndex(
	ctx context.Context, table string, index string, keys KeySet, cols []string,
) RowIterator {
	return newRowIterator(t.ReadOnlyTransaction.ReadUsingIndex(ctx, table, index, keys, cols))
}

func (t *readOnlyTransaction) ReadWithOptions(
	ctx context.Context, table string, keys KeySet, cols []string, options *ReadOptions,
) RowIterator {
	return newRowIterator(t.ReadOnlyTransaction.ReadWithOptions(ctx, table, keys, cols, options))
}

func (t *readOnlyTransaction) WithTimestampBound(tb TimestampBound) ReadOnlyTransaction {
	return newReadOnlyTransaction(t.ReadOnlyTransaction.WithTimestampBound(tb))
}

type batchReadOnlyTransaction struct {
	*spanner.BatchReadOnlyTransaction
}

func newBatchReadOnlyTransaction(transaction *spanner.BatchReadOnlyTransaction) BatchReadOnlyTransaction {
	return &batchReadOnlyTransaction{
		BatchReadOnlyTransaction: transaction,
	}
}

func (t *batchReadOnlyTransaction) GetID() BatchReadOnlyTransactionID {
	return t.BatchReadOnlyTransaction.ID
}

func (t *batchReadOnlyTransaction) Execute(ctx context.Context, partition *Partition) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.Execute(ctx, partition))
}

func (t *batchReadOnlyTransaction) Query(ctx context.Context, stm Statement) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.Query(ctx, stm))
}

func (t *batchReadOnlyTransaction) QueryWithStats(ctx context.Context, stm Statement) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.QueryWithStats(ctx, stm))
}

func (t *batchReadOnlyTransaction) Read(ctx context.Context, table string, keys KeySet, cols []string) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.Read(ctx, table, keys, cols))
}

func (t *batchReadOnlyTransaction) ReadRow(ctx context.Context, table string, key Key, cols []string) (Row, error) {
	row, err := t.BatchReadOnlyTransaction.ReadRow(ctx, table, key, cols)
	if err != nil {
		return nil, err
	}
	return newRow(row), nil
}

func (t *batchReadOnlyTransaction) ReadUsingIndex(
	ctx context.Context, table, index string, keys KeySet, cols []string,
) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.ReadUsingIndex(ctx, table, index, keys, cols))
}

func (t *batchReadOnlyTransaction) ReadWithOptions(
	ctx context.Context, table string, keys KeySet, cols []string, options *ReadOptions,
) RowIterator {
	return newRowIterator(t.BatchReadOnlyTransaction.ReadWithOptions(ctx, table, keys, cols, options))
}
