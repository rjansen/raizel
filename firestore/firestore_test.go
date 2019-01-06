package firestore

import (
	"context"
	"fmt"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/require"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testClient struct {
	name   string
	client *firestore.Client
	server *mockServer
	path   string
	err    error
}

func (scenario *testClient) setup(t *testing.T) {
	if scenario.err == nil {
		c, srv := newMock(t)
		scenario.client = c
		scenario.server = srv
	}
}

func TestClient(test *testing.T) {
	scenarios := []testClient{
		{
			name: "Creates a new client",
			path: "mockcol1/mockref1",
			err:  nil,
		},
		{
			name: "Returns error because client is blank",
			err:  ErrBlankFirestoreClient,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				cli, err := newClient(scenario.client)
				require.Equal(t, scenario.err, err, "newClient error")
				if scenario.err == nil {
					require.NotNil(t, cli, "invalid client")
					reference := cli.Doc(scenario.path)
					require.NotNil(t, reference, "invalid reference")
				}
			},
		)
	}
}

type testSetDocumentRef struct {
	name    string
	client  *firestore.Client
	server  *mockServer
	path    string
	err     error
	data    interface{}
	options []SetOption
}

func (scenario *testSetDocumentRef) setup(t *testing.T) {
	c, srv := newMock(t)
	scenario.client = c
	scenario.server = srv

	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
	} else {
		srv.addRPC(nil,
			&pb.CommitResponse{
				WriteResults: []*pb.WriteResult{
					{UpdateTime: aTimestamp},
				},
			},
		)
	}
}

func TestSetDocumentRef(test *testing.T) {
	scenarios := []testSetDocumentRef{
		{
			name: "Merges document ref",
			path: "mockcoll1/mockref1",
			data: map[string]interface{}{
				"id": "#mockref1", "mame": "Mock One", "age": 25,
			},
			err: nil,
			options: []SetOption{
				MergeAll,
			},
		},
		{
			name: "Creates new document ref",
			path: "mockcoll1/newref1",
			data: map[string]interface{}{
				"id": "#newref1", "mame": "Mock New", "age": 18,
			},
			err: nil,
		},
		{
			name: "Returns error when bad set data is submitted",
			path: "mockcoll1/mockref1",
			data: map[string]interface{}{},
			err:  status.Error(codes.InvalidArgument, "mockScenarioError"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				ref := documentRef{
					scenario.client.Doc(scenario.path),
				}
				err := ref.Set(context.Background(), scenario.data, scenario.options...)
				require.Equal(t, grpc.Code(scenario.err), grpc.Code(err), "set invalid grpccode")
				require.Equal(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "set invalid grpcdesc")
			},
		)
	}
}

type testDeleteDocumentRef struct {
	name   string
	client Client
	server *mockServer
	path   string
	err    error
}

func (scenario *testDeleteDocumentRef) setup(t *testing.T) {
	c, srv := newMock(t)
	cli, err := newClient(c)
	require.Nil(t, err, "delete setup client err")
	scenario.client = cli
	scenario.server = srv

	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
	} else {
		srv.addRPC(nil,
			&pb.CommitResponse{
				WriteResults: []*pb.WriteResult{
					{UpdateTime: aTimestamp},
				},
			},
		)
	}
}

func TestDeleteDocumentRef(test *testing.T) {
	scenarios := []testDeleteDocumentRef{
		{
			name: "Deletes document ref",
			path: "mockcoll1/mockref1",
			err:  nil,
		},
		{
			name: "Returns the delete error cause",
			path: "mockcoll1/mockref1",
			err:  status.Error(codes.NotFound, "mockScenarioError"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				ref := scenario.client.Doc(scenario.path)
				require.NotNil(t, ref, "invalid reference")
				err := ref.Delete(context.Background())
				require.Equal(t, grpc.Code(scenario.err), grpc.Code(err), "delete invalid grpccode")
				require.Equal(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "delete invalid grpcdesc")
			},
		)
	}
}

type testGetDocumentRef struct {
	name   string
	client *firestore.Client
	server *mockServer
	path   string
	err    error
	data   map[string]*pb.Value
}

func (scenario *testGetDocumentRef) setup(t *testing.T) {
	c, srv := newMock(t)
	scenario.client = c
	scenario.server = srv

	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
	} else {
		srv.addRPC(nil,
			[]interface{}{
				&pb.BatchGetDocumentsResponse{
					Result: &pb.BatchGetDocumentsResponse_Found{
						&pb.Document{
							Name:       scenario.path,
							CreateTime: aTimestamp,
							UpdateTime: aTimestamp,
							Fields:     scenario.data,
						},
					},
					ReadTime: aTimestamp2,
				},
			},
		)
	}
}

func TestGetDocumentRef(test *testing.T) {
	scenarios := []testGetDocumentRef{
		{
			name: "Gets document ref",
			path: "mockcoll1/mockref1",
			data: map[string]*pb.Value{
				"id":   strval("#mockref1"),
				"mame": strval("Mock One"),
				"age":  intval(25),
			},
			err: nil,
		},
		{
			name: "Returns a server error",
			path: "mockcoll1/badref1",
			data: nil,
			err:  status.Error(codes.Unknown, "mockBadGateway"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				ref := documentRef{
					scenario.client.Doc(scenario.path),
				}
				_, err := ref.Get(context.Background())
				require.Equal(t, grpc.Code(scenario.err), grpc.Code(err), "get invalid grpccode")
				require.Equal(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "get invalid grpcdesc")
			},
		)
	}
}

type testGetAllDocumentRef struct {
	name       string
	client     Client
	server     *mockServer
	err        error
	paths      []string
	data       []map[string]*pb.Value
	references []DocumentRef
}

func (scenario *testGetAllDocumentRef) setup(t *testing.T) {
	c, srv := newMock(t)
	cli, err := newClient(c)
	require.Nil(t, err, "getall setup client err")
	scenario.client = cli
	scenario.server = srv

	references := make([]DocumentRef, len(scenario.paths))
	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
		for index, path := range scenario.paths {
			references[index] = scenario.client.Doc(path)
		}
	} else {
		mockResults := make([]interface{}, len(scenario.paths))
		for index, path := range scenario.paths {
			mockResults[index] = &pb.BatchGetDocumentsResponse{
				Result: &pb.BatchGetDocumentsResponse_Found{
					&pb.Document{
						Name:       path,
						CreateTime: aTimestamp,
						UpdateTime: aTimestamp,
						Fields:     scenario.data[index],
					},
				},
				ReadTime: aTimestamp2,
			}
			references[index] = scenario.client.Doc(path)
		}
		srv.addRPC(nil, mockResults)
	}
	scenario.references = references
}

func TestGetAllDocumentRef(test *testing.T) {
	scenarios := []testGetAllDocumentRef{
		{
			name: "Gets all documents ref",
			paths: []string{
				"mockcoll1/mockref1",
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(25),
				},
			},
			err: nil,
		},
		{
			name: "Returns a server error",
			paths: []string{
				"mockcoll1/badref1",
			},
			data: nil,
			err:  status.Error(codes.Unknown, "mockBadGateway"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				docs, err := scenario.client.GetAll(
					context.Background(),
					scenario.references...,
				)
				require.Equal(t, grpc.Code(scenario.err), grpc.Code(err), "getall invalid grpccode")
				require.Equal(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "getall invalid grpcdesc")
				if scenario.err == nil {
					require.Len(t, docs, len(scenario.paths), "invalid result documents size")
				} else {
					require.Len(t, docs, 0, "not empty result documents")
				}
			},
		)
	}
}
