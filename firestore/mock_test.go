package firestore

import (
	"context"
	"crypto/sha1"
	"fmt"
	"hash/fnv"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

func entityMockRef(collection string, id interface{}) string {
	return fmt.Sprintf("%s/%s", collection, id)
}

func newUUID() string {
	return uuid.New().String()
}

func newID() int {
	hash32 := fnv.New32()
	hash32.Write([]byte(newUUID()))
	return int(hash32.Sum32())
}

func Sha1(v string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(v)))
}

func Sha1f(f string, a ...interface{}) string {
	return Sha1(fmt.Sprintf(f, a...))
}

type dynamicData map[string]interface{}

type entityMock struct {
	ID        string      `firestore:"id"`
	Name      string      `firestore:"name"`
	Age       int         `firestore:"age"`
	Data      dynamicData `firestore:"data"`
	Deleted   bool        `firestore:"deleted"`
	CreatedAt time.Time   `firestore:"created_at"`
	UpdatedAt time.Time   `firestore:"updated_at"`
}

type entityKeyMock struct {
	collection string
	name       string
	value      interface{}
}

func (k entityKeyMock) Name() string {
	return k.name
}

func (k entityKeyMock) Value() interface{} {
	return k.value
}

func (k entityKeyMock) EntityName() string {
	return k.collection
}

type documentRefMock struct {
	mock.Mock
}

func newDocumentRefMock() *documentRefMock {
	return new(documentRefMock)
}

func (mock *documentRefMock) Get(ctx context.Context) (DocumentSnapshot, error) {
	var (
		args   = mock.Called(ctx)
		result = args.Get(0)
		err    = args.Error(1)
	)
	if result == nil {
		return nil, err
	}
	return result.(DocumentSnapshot), err
}

func (mock *documentRefMock) Set(ctx context.Context, data interface{}, opts ...SetOption) error {
	args := mock.Called(ctx, data, opts)
	return args.Error(0)
}

func (mock *documentRefMock) Delete(ctx context.Context) error {
	args := mock.Called(ctx)
	return args.Error(0)
}

func (mock *documentRefMock) delegate() *firestore.DocumentRef {
	var (
		args   = mock.Called()
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(*firestore.DocumentRef)
}

type documentSnapshotMock struct {
	mock.Mock
}

func newDocumentSnapshotMock() *documentSnapshotMock {
	return new(documentSnapshotMock)
}

func (mock *documentSnapshotMock) DataTo(data interface{}) error {
	args := mock.Called(data)
	return args.Error(0)
}

func (mock *documentSnapshotMock) Exists() bool {
	args := mock.Called()
	return args.Bool(0)
}

type clientMock struct {
	mock.Mock
}

func newClientMock() *clientMock {
	return new(clientMock)
}

func (mock *clientMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *clientMock) Doc(path string) DocumentRef {
	var (
		args   = mock.Called(path)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(DocumentRef)
}

func (mock *clientMock) Collection(path string) CollectionRef {
	var (
		args   = mock.Called(path)
		result = args.Get(0)
	)
	if result == nil {
		return nil
	}
	return result.(CollectionRef)
}

func (mock *clientMock) GetAll(ctx context.Context, refs ...DocumentRef) ([]DocumentSnapshot, error) {
	var (
		args   = mock.Called(ctx, refs)
		result = args.Get(0)
		err    = args.Error(1)
	)
	if result == nil {
		return nil, err
	}
	return result.([]DocumentSnapshot), err
}
