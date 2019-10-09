package spanner

import (
	"testing"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/require"
)

func TestReadOnlyTransaction(t *testing.T) {
	transaction := newReadOnlyTransaction(new(spanner.ReadOnlyTransaction))
	require.NotNil(t, transaction, "invalid transaction instance")
}

func TestWithTimestampBound(t *testing.T) {
	transaction := newReadOnlyTransaction(new(spanner.ReadOnlyTransaction))
	require.NotNil(t, transaction, "invalid transaction instance")
	boundedTransaction := transaction.WithTimestampBound(TimestampBound{})
	require.NotNil(t, boundedTransaction, "invalid bounded transaction instance")
}

func TestBatchReadOnlyTransaction(t *testing.T) {
	transaction := newBatchReadOnlyTransaction(new(spanner.BatchReadOnlyTransaction))
	require.NotNil(t, transaction, "invalid transaction instance")
}
