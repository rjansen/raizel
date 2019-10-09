package spanner

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/rjansen/raizel/spanner/internal/testutil"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func newSpannerClientMock(t *testing.T) (*testutil.MockCloudSpanner, *spanner.Client) {
	mockServer := testutil.NewMockCloudSpanner(t, time.Now())
	require.NotNil(t, mockServer, "invalid mock server instance")
	mockServer.Serve()
	t.Logf("spannerMockServer.Addr=%s", mockServer.Addr())
	require.NotNil(t, mockServer.Addr(), "invalid mock server addr instance")
	conn, err := grpc.Dial(mockServer.Addr(), grpc.WithInsecure())
	require.Nil(t, err, "grpc connection dial error")
	require.NotNil(t, conn, "invalid grpc connection instance")
	client, err := spanner.NewClient(
		context.Background(),
		"projects/mockproject/instances/mockinstance/databases/mockdb",
		option.WithGRPCConn(conn),
	)
	require.Nil(t, err, "new spanner client error")
	require.NotNil(t, client, "invalid spanner client instance")
	return mockServer, client
}

func TestClient(t *testing.T) {
	client := NewClient(new(spanner.Client))
	require.NotNil(t, client, "invalid client instance")
	require.Implements(t, (*Client)(nil), client, "invalid client type")
}

func TestClientReadOnlyTransaction(t *testing.T) {
	client := NewClient(new(spanner.Client))
	require.NotNil(t, client, "invalid client instance")
	transaction := client.ReadOnlyTransaction()
	require.NotNil(t, transaction, "invalid transaction instance")
}

func TestClientSingle(t *testing.T) {
	client := NewClient(new(spanner.Client))
	require.NotNil(t, client, "invalid client instance")
	transaction := client.Single()
	require.NotNil(t, transaction, "invalid transaction instance")
}

func TestClientBatchReadOnlyTransaction(t *testing.T) {
	mockServer, mockClient := newSpannerClientMock(t)
	defer mockServer.Stop()

	client := NewClient(mockClient)
	require.NotNil(t, client, "invalid client instance")
	transaction, err := client.BatchReadOnlyTransaction(context.Background(), spanner.StrongRead())
	require.Nil(t, err, "transaction error")
	require.NotNil(t, transaction, "invalid transaction instance")

	transaction = client.BatchReadOnlyTransactionFromID(transaction.GetID())
}
