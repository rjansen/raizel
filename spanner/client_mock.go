package spanner

import (
	"context"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func (mock *ClientMock) Apply(ctx context.Context, ms []*Mutation, opts ...ApplyOption) (time.Time, error) {
	args := mock.Called(ctx, ms, opts)
	return args.Get(0).(time.Time), args.Error(1)
}

func (mock *ClientMock) BatchReadOnlyTransaction(ctx context.Context, tb TimestampBound) (BatchReadOnlyTransaction, error) {
	args := mock.Called(ctx, tb)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(BatchReadOnlyTransaction), args.Error(1)

}

func (mock *ClientMock) BatchReadOnlyTransactionFromID(tid BatchReadOnlyTransactionID) BatchReadOnlyTransaction {
	args := mock.Called(tid)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(BatchReadOnlyTransaction)
}

func (mock *ClientMock) Close() {
	mock.Called()
}

func (mock *ClientMock) PartitionedUpdate(ctx context.Context, statement Statement) (int64, error) {
	args := mock.Called(ctx, statement)
	return args.Get(0).(int64), args.Error(1)
}

func (mock *ClientMock) ReadOnlyTransaction() ReadOnlyTransaction {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(ReadOnlyTransaction)
}

func (mock *ClientMock) ReadWriteTransaction(
	ctx context.Context, f func(context.Context, *spanner.ReadWriteTransaction) error,
) (time.Time, error) {
	args := mock.Called(ctx, f)
	return args.Get(0).(time.Time), args.Error(1)
}

func (mock *ClientMock) Single() ReadOnlyTransaction {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(ReadOnlyTransaction)
}
