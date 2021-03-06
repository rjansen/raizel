package sql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/stretchr/testify/mock"
)

type dynamicData map[string]interface{}

func (d dynamicData) Value() (driver.Value, error) {
	j, err := json.Marshal(d)
	return j, err
}

func (d *dynamicData) Scan(src interface{}) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("err_invalid_dbtype: != []byte")
	}

	err := json.Unmarshal(source, d)
	if err != nil {
		return err
	}
	return nil
}

type entityMock struct {
	ID        int         `db:"id"`
	Name      string      `db:"name"`
	Age       int         `db:"age"`
	Data      dynamicData `db:"data"`
	Deleted   bool        `db:"deleted"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}

type entityKeyMock struct {
	table string
	name  string
	value interface{}
}

func (k entityKeyMock) Name() string {
	return k.name
}

func (k entityKeyMock) Value() interface{} {
	return k.value
}

func (k entityKeyMock) EntityName() string {
	return k.table
}

func newDBMock() *dbMock {
	return new(dbMock)
}

type dbMock struct {
	mock.Mock
}

func (mock *dbMock) QueryRow(sql string, params ...interface{}) Row {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Row)
	}
	return nil
}

func (mock *dbMock) Query(sql string, params ...interface{}) (Rows, error) {
	var (
		args   = mock.Called(sql, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Rows), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *dbMock) Exec(query string, params ...interface{}) (Result, error) {
	var (
		args   = mock.Called(query, params)
		result = args.Get(0)
	)
	if result != nil {
		return result.(Result), args.Error(1)
	}
	return nil, args.Error(1)
}

func (mock *dbMock) Ping() error {
	args := mock.Called()
	return args.Error(0)
}

func (mock *dbMock) Close() error {
	args := mock.Called()
	return args.Error(0)
}

func newRowMock() *rowMock {
	return new(rowMock)
}

type rowMock struct {
	mock.Mock
}

func (mock *rowMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func newRowsMock() *rowsMock {
	return new(rowsMock)
}

type rowsMock struct {
	mock.Mock
}

func (mock *rowsMock) Next() bool {
	args := mock.Called()
	return args.Bool(0)
}

func (mock *rowsMock) Scan(dest ...interface{}) error {
	args := mock.Called(dest)
	return args.Error(0)
}

func newResultMock() *resultMock {
	return new(resultMock)
}

type resultMock struct {
	mock.Mock
}

func (mock *resultMock) LastInsertId() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (mock *resultMock) RowsAffected() (int64, error) {
	args := mock.Called()
	return args.Get(0).(int64), args.Error(1)
}
