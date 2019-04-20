package spanner

import (
	"cloud.google.com/go/spanner"
	sppb "google.golang.org/genproto/googleapis/spanner/v1"
)

type Key = spanner.Key

type KeySet = spanner.KeySet

type KeyRange = spanner.KeyRange

type ApplyOption = spanner.ApplyOption

type Mutation = spanner.Mutation

type Statement = spanner.Statement

type TimestampBound = spanner.TimestampBound

type ReadWriteTransaction = spanner.ReadWriteTransaction

type BatchReadOnlyTransactionID = spanner.BatchReadOnlyTransactionID

type ReadOptions = spanner.ReadOptions

type QueryPlan = sppb.QueryPlan

type Partition = spanner.Partition

type PartitionOptions = spanner.PartitionOptions

type Row interface {
	Column(int, interface{}) error
	ColumnByName(string, interface{}) error
	ColumnIndex(string) (int, error)
	ColumnName(int) string
	ColumnNames() []string
	Columns(...interface{}) error
	Size() int
	ToStruct(interface{}) error
}

type RowIterator interface {
	Do(func(Row) error) error
	Next() (Row, error)
	Stop()
}

type row struct {
	*spanner.Row
}

func newRow(r *spanner.Row) Row {
	return &row{Row: r}
}

type rowIterator struct {
	*spanner.RowIterator
}

func newRowIterator(r *spanner.RowIterator) RowIterator {
	return &rowIterator{RowIterator: r}
}

func (i *rowIterator) Do(f func(Row) error) error {
	return i.RowIterator.Do(
		func(r *spanner.Row) error {
			return f(newRow(r))
		},
	)
}

func (i *rowIterator) Next() (Row, error) {
	row, err := i.RowIterator.Next()
	if err != nil {
		return nil, err
	}
	return newRow(row), nil
}

func (i *rowIterator) Stop() {
	i.RowIterator.Stop()
}
