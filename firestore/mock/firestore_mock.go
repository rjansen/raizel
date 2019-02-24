package mock

import (
	"context"

	"github.com/rjansen/raizel/firestore"
	"github.com/stretchr/testify/mock"
)

type CollectionRefMock struct {
	firestore.CollectionRef
	mock.Mock
}

func NewCollectionRefMock() *CollectionRefMock {
	return new(CollectionRefMock)
}

func (mock *CollectionRefMock) Documents(ctx context.Context) firestore.DocumentIterator {
	var (
		args   = mock.Called(ctx)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.DocumentIterator)
}

func (mock *CollectionRefMock) Where(path string, op string, value interface{}) firestore.Query {
	var (
		args   = mock.Called(path, op, value)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.Query)
}

func (mock *CollectionRefMock) OrderBy(path string, direction firestore.Direction) firestore.Query {
	var (
		args   = mock.Called(path, direction)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.Query)
}

func (mock *CollectionRefMock) Offset(n int) firestore.Query {
	var (
		args   = mock.Called(n)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.Query)
}

func (mock *CollectionRefMock) Limit(n int) firestore.Query {
	var (
		args   = mock.Called(n)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.Query)
}

type DocumentSnapshotMock struct {
	mock.Mock
}

func NewDocumentSnapshotMock() *DocumentSnapshotMock {
	return new(DocumentSnapshotMock)
}

func (mock *DocumentSnapshotMock) DataTo(data interface{}) error {
	args := mock.Called(data)
	return args.Error(0)
}

func (mock *DocumentSnapshotMock) Exists() bool {
	args := mock.Called()
	return args.Bool(0)
}

type DocumentRefMock struct {
	firestore.DocumentRef
	mock.Mock
}

func NewDocumentRefMock() *DocumentRefMock {
	return new(DocumentRefMock)
}

func (mock *DocumentRefMock) Get(ctx context.Context) (firestore.DocumentSnapshot, error) {
	var (
		args   = mock.Called(ctx)
		result = args.Get(0)
		err    = args.Error(1)
	)
	if result == nil {
		return nil, err
	}
	return result.(firestore.DocumentSnapshot), err
}

func (mock *DocumentRefMock) Set(ctx context.Context, data interface{}, opts ...firestore.SetOption) error {
	args := mock.Called(ctx, data, opts)
	return args.Error(0)
}

func (mock *DocumentRefMock) Delete(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

type DocumentIteratorMock struct {
	mock.Mock
}

func NewDocumentIteratorMock() *DocumentIteratorMock {
	return new(DocumentIteratorMock)
}

func (mock *DocumentIteratorMock) GetAll() ([]firestore.DocumentSnapshot, error) {
	var (
		args   = mock.Called()
		result = args.Get(0)
		err    = args.Error(1)
	)
	if result == nil {
		return nil, err
	}
	return result.([]firestore.DocumentSnapshot), err
}

type WriteBatchMock struct {
	mock.Mock
}

func NewWriteBatchMock() *WriteBatchMock {
	return new(WriteBatchMock)
}

func (mock *WriteBatchMock) Set(ref firestore.DocumentRef, data interface{}, options ...firestore.SetOption) firestore.WriteBatch {
	var (
		args   = mock.Called(ref, data, options)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.WriteBatch)
}

func (mock *WriteBatchMock) Delete(ref firestore.DocumentRef) firestore.WriteBatch {
	var (
		args   = mock.Called(ref)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.WriteBatch)
}

func (mock *WriteBatchMock) Commit(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return new(ClientMock)
}

func (mock *ClientMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *ClientMock) Doc(path string) firestore.DocumentRef {
	var (
		args   = mock.Called(path)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.DocumentRef)
}

func (mock *ClientMock) Collection(path string) firestore.CollectionRef {
	var (
		args   = mock.Called(path)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.CollectionRef)
}

func (mock *ClientMock) GetAll(
	ctx context.Context, refs ...firestore.DocumentRef,
) ([]firestore.DocumentSnapshot, error) {
	var (
		args   = mock.Called(ctx, refs)
		result = args.Get(0)
		err    = args.Error(1)
	)
	if result == nil {
		return nil, err
	}
	return result.([]firestore.DocumentSnapshot), err
}

func (mock *ClientMock) Batch() firestore.WriteBatch {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(firestore.WriteBatch)
}
