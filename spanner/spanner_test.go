package spanner

import (
	"context"
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/require"
)

func TestRow(t *testing.T) {
	row := newRow(new(spanner.Row))
	require.NotNil(t, row, "invalid row instance")
}

func TestRowIterator(t *testing.T) {
	iterator := newRowIterator(new(spanner.RowIterator))
	require.NotNil(t, iterator, "invalid iterator instance")
}

func TestRowIteratorDo(t *testing.T) {
	mockServer, mockClient := newSpannerClientMock(t)
	defer mockServer.Stop()

	client := NewClient(mockClient)
	iterator := client.ReadOnlyTransaction().Query(context.Background(), Statement{SQL: "UPDATE t SET x = 2 WHERE x = 1"})
	err := iterator.Do(func(Row) error { return nil })
	require.Nil(t, err, "iterator do error")
	iterator.Stop()
}

func TestRowIteratorNext(t *testing.T) {
	mockServer, mockClient := newSpannerClientMock(t)
	defer mockServer.Stop()

	client := NewClient(mockClient)
	iterator := client.ReadOnlyTransaction().Query(context.Background(), Statement{SQL: "SELECT column1, column2, columnN from available"})
	_, err := iterator.Next()
	require.Nil(t, err, "iterator next error")
	iterator.Stop()
}
