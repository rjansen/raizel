package spanner

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/mock"
)

type ReadOnlyTransactionMock struct {
	mock.Mock
}

func (mock *ReadOnlyTransactionMock) AnalyzeQuery(ctx context.Context, statement spanner.Statement) (*QueryPlan, error) {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*QueryPlan), args.Error(1)
}

func (mock *ReadOnlyTransactionMock) Close() {
	mock.Called()
}

func (mock *ReadOnlyTransactionMock) Query(ctx context.Context, statement spanner.Statement) RowIterator {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *ReadOnlyTransactionMock) QueryWithStats(ctx context.Context, statement Statement) RowIterator {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *ReadOnlyTransactionMock) Read(
	ctx context.Context, table string, keys KeySet, columns []string,
) RowIterator {
	args := mock.Called(ctx, table, keys, columns)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *ReadOnlyTransactionMock) ReadRow(
	ctx context.Context, table string, key spanner.Key, columns []string,
) (Row, error) {
	args := mock.Called(ctx, table, key, columns)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(Row), args.Error(1)
}

func (mock *ReadOnlyTransactionMock) ReadUsingIndex(
	ctx context.Context, table, index string, keys KeySet, columns []string,
) RowIterator {
	args := mock.Called(ctx, table, index, keys, columns)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *ReadOnlyTransactionMock) ReadWithOptions(
	ctx context.Context, table string, keys KeySet, columns []string, opts *ReadOptions,
) RowIterator {
	args := mock.Called(ctx, table, keys, columns, opts)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *ReadOnlyTransactionMock) Timestamp() (time.Time, error) {
	args := mock.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (mock *ReadOnlyTransactionMock) WithTimestampBound(tb TimestampBound) ReadOnlyTransaction {
	args := mock.Called(tb)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(ReadOnlyTransaction)
}

type BatchReadOnlyTransactionMock struct {
	mock.Mock
}

func (mock *BatchReadOnlyTransactionMock) GetID() BatchReadOnlyTransactionID {
	args := mock.Called()
	return args.Get(0).(BatchReadOnlyTransactionID)
}

func (mock *BatchReadOnlyTransactionMock) AnalyzeQuery(ctx context.Context, statement Statement) (*QueryPlan, error) {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*QueryPlan), args.Error(1)
}

func (mock *BatchReadOnlyTransactionMock) Cleanup(ctx context.Context) {
	mock.Called(ctx)
}

func (mock *BatchReadOnlyTransactionMock) Close() {
	mock.Called()
}

func (mock *BatchReadOnlyTransactionMock) Execute(ctx context.Context, p *Partition) RowIterator {
	args := mock.Called(ctx, p)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *BatchReadOnlyTransactionMock) PartitionQuery(
	ctx context.Context, statement Statement, opts PartitionOptions,
) ([]*Partition, error) {
	args := mock.Called(ctx, statement, opts)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*Partition), args.Error(1)
}

func (mock *BatchReadOnlyTransactionMock) PartitionRead(
	ctx context.Context, table string, keys KeySet, columns []string, opts PartitionOptions,
) ([]*spanner.Partition, error) {
	args := mock.Called(ctx, table, keys, columns, opts)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*Partition), args.Error(1)
}

func (mock *BatchReadOnlyTransactionMock) PartitionReadUsingIndex(
	ctx context.Context, table, index string, keys KeySet, columns []string, opts PartitionOptions,
) ([]*spanner.Partition, error) {
	args := mock.Called(ctx, table, index, keys, columns, opts)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]*Partition), args.Error(1)
}

func (mock *BatchReadOnlyTransactionMock) Query(ctx context.Context, statement Statement) RowIterator {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *BatchReadOnlyTransactionMock) QueryWithStats(ctx context.Context, statement Statement) RowIterator {
	args := mock.Called(ctx, statement)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *BatchReadOnlyTransactionMock) Read(
	ctx context.Context, table string, keys KeySet, columns []string,
) RowIterator {
	args := mock.Called(ctx, table, keys, columns)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *BatchReadOnlyTransactionMock) ReadRow(
	ctx context.Context, table string, key Key, columns []string,
) (Row, error) {
	args := mock.Called(ctx, table, key, columns)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(Row), args.Error(1)
}

func (mock *BatchReadOnlyTransactionMock) ReadUsingIndex(
	ctx context.Context, table, index string, keys KeySet, columns []string,
) RowIterator {
	args := mock.Called(ctx, table, index, keys, columns)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}

func (mock *BatchReadOnlyTransactionMock) ReadWithOptions(
	ctx context.Context, table string, keys KeySet, columns []string, opts *ReadOptions) RowIterator {
	args := mock.Called(ctx, table, keys, columns, opts)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(RowIterator)
}
