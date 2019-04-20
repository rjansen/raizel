package spanner

import "github.com/stretchr/testify/mock"

type RowMock struct {
	mock.Mock
}

func NewRowMock() *RowMock {
	return new(RowMock)
}

func (mock *RowMock) Column(i int, ptr interface{}) error {
	args := mock.Called(i, ptr)
	return args.Error(0)
}

func (mock *RowMock) ColumnByName(name string, ptr interface{}) error {
	args := mock.Called(name, ptr)
	return args.Error(0)
}

func (mock *RowMock) ColumnIndex(name string) (int, error) {
	args := mock.Called(name)
	return args.Int(0), args.Error(1)
}

func (mock *RowMock) ColumnName(i int) string {
	args := mock.Called(i)
	return args.String(0)
}

func (mock *RowMock) ColumnNames() []string {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.([]string)
}

func (mock *RowMock) Columns(ptrs ...interface{}) error {
	args := mock.Called(ptrs)
	return args.Error(0)
}

func (mock *RowMock) Size() int {
	args := mock.Called()
	return args.Int(0)
}

func (mock *RowMock) ToStruct(ptr interface{}) error {
	args := mock.Called(ptr)
	return args.Error(0)
}

type RowIteratorMock struct {
	mock.Mock
}

func NewRowIteratorMock() *RowIteratorMock {
	return new(RowIteratorMock)
}

func (mock *RowIteratorMock) Do(f func(Row) error) error {
	args := mock.Called(f)
	return args.Error(0)
}

func (mock *RowIteratorMock) Next() (Row, error) {
	args := mock.Called()
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(Row), args.Error(1)
}

func (mock *RowIteratorMock) Stop() {
	mock.Called()
}
