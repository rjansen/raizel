package mock

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/rjansen/raizel"
	"github.com/stretchr/testify/mock"
)

type Object map[string]interface{}

func (o Object) Value() (driver.Value, error) {
	j, err := json.Marshal(o)
	return j, err
}

func (o *Object) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("err_invalid_dbtype: != []byte")
	}

	err := json.Unmarshal(source, o)
	if err != nil {
		return err
	}
	return nil
}

type MockEntity struct {
	ID       string    `json:"id" firestore:"id" db:"id"`
	String   string    `json:"string" firestore:"string" db:"string"`
	Integer  int32     `json:"integer" firestore:"integer" db:"integer"`
	Float    float32   `json:"float" firestore:"float" db:"float"`
	DateTime time.Time `json:"date_time" firestore:"date_time" db:"date_time"`
	Boolean  bool      `json:"boolean" firestore:"boolean" db:"boolean"`
	Object   Object    `json:"object" firestore:"object" db:"object"`
}

func NewMockEntity() *MockEntity {
	return new(MockEntity)
}

type MockEntityKey struct {
	mock.Mock
}

func NewMockEntityKey() *MockEntityKey {
	return new(MockEntityKey)
}

func (mock *MockEntityKey) EntityName() string {
	args := mock.Called()
	return args.String(0)
}

func (mock *MockEntityKey) Value() interface{} {
	args := mock.Called()
	return args.Get(0)
}

func (mock *MockEntityKey) Name() string {
	args := mock.Called()
	return args.String(0)
}

type MockRepository struct {
	mock.Mock
}

func NewMockRepository() *MockRepository {
	return new(MockRepository)
}

func (mock *MockRepository) Get(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	args := mock.Called(ctx, key, entity)
	return args.Error(0)
}

func (mock *MockRepository) Set(ctx context.Context, key raizel.EntityKey, entity raizel.Entity) error {
	args := mock.Called(ctx, key, entity)
	return args.Error(0)
}

func (mock *MockRepository) Delete(ctx context.Context, key raizel.EntityKey) error {
	args := mock.Called(ctx, key)
	return args.Error(0)
}

func (mock *MockRepository) Close(context.Context) error {
	args := mock.Called()
	return args.Error(0)
}
