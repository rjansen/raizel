package mock

import (
	"errors"
	"testing"

	"github.com/rjansen/raizel/firestore"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCollectionRefMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var ref *CollectionRefMock = NewCollectionRefMock()
			require.NotNil(t, ref, "invalid collection_ref instance")
			require.Implements(t, (*firestore.CollectionRef)(nil), ref, "invalid collection_ref type")
			require.Implements(t, (*firestore.Query)(nil), ref, "invalid query type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			ref := NewCollectionRefMock()
			ref.On("Documents", mock.Anything).Return(nil)
			ref.On("Where", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			ref.On("OrderBy", mock.Anything, mock.Anything).Return(nil)
			ref.On("Offset", mock.Anything).Return(nil)
			ref.On("Limit", mock.Anything).Return(nil)

			require.Nil(t, ref.Documents(nil), "invalid documents() response")
			require.Nil(t, ref.Where("", "", nil), "invalid where() response")
			require.Nil(t, ref.OrderBy("", 0), "invalid order_by() response")
			require.Nil(t, ref.Offset(0), "invalid offset() response")
			require.Nil(t, ref.Limit(0), "invalid limit() response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				ref      = NewCollectionRefMock()
				iterator = NewDocumentIteratorMock()
			)

			ref.On("Documents", mock.Anything).Return(iterator)
			ref.On("Where", mock.Anything, mock.Anything, mock.Anything).Return(ref)
			ref.On("OrderBy", mock.Anything, mock.Anything).Return(ref)
			ref.On("Offset", mock.Anything).Return(ref)
			ref.On("Limit", mock.Anything).Return(ref)

			require.Equal(t, iterator, ref.Documents(nil), "invalid documents() response")
			require.Equal(t, ref, ref.Where("", "", nil), "invalid where() response")
			require.Equal(t, ref, ref.OrderBy("", 0), "invalid order_by() response")
			require.Equal(t, ref, ref.Offset(0), "invalid offset() response")
			require.Equal(t, ref, ref.Limit(0), "invalid limit() response")
		},
	)

}

func TestDocumentSnapshotMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var ref *DocumentSnapshotMock = NewDocumentSnapshotMock()
			require.NotNil(t, ref, "invalid document_ref instance")
			require.Implements(t, (*firestore.DocumentSnapshot)(nil), ref, "invalid document_ref type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			snapshot := NewDocumentSnapshotMock()

			snapshot.On("DataTo", mock.Anything).Return(nil)
			snapshot.On("Exists").Return(false)

			require.Nil(t, snapshot.DataTo(nil), "invalid data_to() response")
			require.False(t, snapshot.Exists(), "invalid exists() response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				snapshot  = NewDocumentSnapshotMock()
				errDataTo = errors.New("err_mmock_data_to")
			)

			snapshot.On("DataTo", mock.Anything).Return(errDataTo)
			snapshot.On("Exists").Return(true)

			require.Equal(t, errDataTo, snapshot.DataTo(nil), "invalid data_to() response")
			require.True(t, snapshot.Exists(), "invalid exists() response")
		},
	)
}

func TestDocumentRefMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var ref *DocumentRefMock = NewDocumentRefMock()
			require.NotNil(t, ref, "invalid document_ref instance")
			require.Implements(t, (*firestore.DocumentRef)(nil), ref, "invalid document_ref type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			ref := NewDocumentRefMock()

			ref.On("Get", mock.Anything).Return(nil, nil)
			ref.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			ref.On("Delete", mock.Anything).Return(nil)

			snapshot, err := ref.Get(nil)
			require.Nil(t, err, "invalid get() error response")
			require.Nil(t, snapshot, "invalid get() snapshot response")
			require.Nil(t, ref.Set(nil, nil), "invalid set() response")
			require.Nil(t, ref.Delete(nil), "invalid delete() response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				ref       = NewDocumentRefMock()
				snapshot  = NewDocumentSnapshotMock()
				errGet    = errors.New("err_mock_get")
				errSet    = errors.New("err_mock_set")
				errDelete = errors.New("err_mock_delete")
			)

			ref.On("Get", mock.Anything).Return(snapshot, errGet)
			ref.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(errSet)
			ref.On("Delete", mock.Anything).Return(errDelete)

			document, err := ref.Get(nil)
			require.Equal(t, errGet, err, "invalid get() error response")
			require.Equal(t, snapshot, document, "invalid get() snapshot response")
			require.Equal(t, errSet, ref.Set(nil, nil), "invalid set() response")
			require.Equal(t, errDelete, ref.Delete(nil), "invalid delete() response")
		},
	)

}

func TestDocumentIteratorMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var iterator *DocumentIteratorMock = NewDocumentIteratorMock()
			require.NotNil(t, iterator, "invalid document_iterator instance")
			require.Implements(t, (*firestore.DocumentIterator)(nil), iterator, "invalid document_iterator type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			iterator := NewDocumentIteratorMock()

			iterator.On("GetAll").Return(nil, nil)

			documents, err := iterator.GetAll()
			require.Nil(t, err, "invalid get_all() error response")
			require.Nil(t, documents, "invalid get_all() documents response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				iterator  = NewDocumentIteratorMock()
				snapshots = []firestore.DocumentSnapshot{
					NewDocumentSnapshotMock(),
					NewDocumentSnapshotMock(),
					NewDocumentSnapshotMock(),
				}
				errGetAll = errors.New("err_mock_get_all")
			)

			iterator.On("GetAll").Return(snapshots, errGetAll)

			documents, err := iterator.GetAll()
			require.Equal(t, errGetAll, err, "invalid get_all() error response")
			require.Equal(t, snapshots, documents, "invalid get_all() documents response")
		},
	)
}

func TestWriteBatchMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var batch *WriteBatchMock = NewWriteBatchMock()
			require.NotNil(t, batch, "invalid write_batch instance")
			require.Implements(t, (*firestore.WriteBatch)(nil), batch, "invalid write_batch type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			batch := NewWriteBatchMock()

			batch.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			batch.On("Delete", mock.Anything).Return(nil)
			batch.On("Commit", mock.Anything).Return(nil)

			require.Nil(t, batch.Set(nil, nil), "invalid set() response")
			require.Nil(t, batch.Delete(nil), "invalid delete() response")
			require.Nil(t, batch.Commit(nil), "invalid commit() response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				batch     = NewWriteBatchMock()
				errCommit = errors.New("err_mock_commit")
			)

			batch.On("Set", mock.Anything, mock.Anything, mock.Anything).Return(batch)
			batch.On("Delete", mock.Anything).Return(batch)
			batch.On("Commit", mock.Anything).Return(errCommit)

			require.Equal(t, batch, batch.Set(nil, nil), "invalid set() response")
			require.Equal(t, batch, batch.Delete(nil), "invalid delete() response")
			require.Equal(t, errCommit, batch.Commit(nil), "invalid commit() response")
		},
	)

}

func TestClientMock(t *testing.T) {
	t.Run(
		"Validates mock interface",
		func(t *testing.T) {
			var client *ClientMock = NewClientMock()
			require.NotNil(t, client, "invalid client instance")
			require.Implements(t, (*firestore.Client)(nil), client, "invalid client type")
		},
	)

	t.Run(
		"Returns nil for function call",
		func(t *testing.T) {
			client := NewClientMock()
			client.On("Collection", mock.Anything).Return(nil)
			client.On("Doc", mock.Anything).Return(nil)
			client.On("Batch").Return(nil)
			client.On("GetAll", mock.Anything, mock.Anything).Return(nil, nil)
			client.On("Close").Return(nil)

			require.Nil(t, client.Collection("collection"), "invalid collection() response")
			require.Nil(t, client.Doc("collection/document"), "invalid doc() response")
			require.Nil(t, client.Batch(), "invalid batch() response")
			documents, err := client.GetAll(nil, nil)
			require.Nil(t, err, "invalid getall() error response")
			require.Nil(t, documents, "invalid getall() documents response")
			require.Nil(t, client.Close(), "invalid close() response")
		},
	)

	t.Run(
		"Returns a configured result for function call",
		func(t *testing.T) {
			var (
				client        = NewClientMock()
				collectionRef = NewCollectionRefMock()
				documentRef   = NewDocumentRefMock()
				snapshots     = []firestore.DocumentSnapshot{
					NewDocumentSnapshotMock(),
					NewDocumentSnapshotMock(),
					NewDocumentSnapshotMock(),
				}
				batch     = NewWriteBatchMock()
				errGetAll = errors.New("err_mock_get_all")
				errClose  = errors.New("err_close")
			)

			client.On("Collection", mock.Anything).Return(collectionRef)
			client.On("Doc", mock.Anything).Return(documentRef)
			client.On("GetAll", mock.Anything, mock.Anything).Return(snapshots, errGetAll)
			client.On("Batch").Return(batch)
			client.On("Close").Return(errClose)

			require.Equal(t, collectionRef, client.Collection("collection"), "invalid collection() response")
			require.Equal(t, documentRef, client.Doc("collection/document"), "invalid doc() response")
			require.Equal(t, batch, client.Batch(), "invalid batch() response")
			documents, err := client.GetAll(nil, nil)
			require.Equal(t, errGetAll, err, "invalid getall() error response")
			require.Equal(t, snapshots, documents, "invalid getall() documents response")
			require.Equal(t, errClose, client.Close(), "invalid close() response")
		},
	)
}
