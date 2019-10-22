package firestore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const rootDocumentsMock = "projects/projectID/databases/(default)/documents/"

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
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
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
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
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
						Found: &pb.Document{
							Name:       rootDocumentsMock + scenario.path,
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
				ref := documentRef{scenario.client.Doc(scenario.path)}
				_, err := ref.Get(context.Background())
				assert.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				assert.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
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
					Found: &pb.Document{
						Name:       rootDocumentsMock + path,
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
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
				if scenario.err == nil {
					require.Len(t, docs, len(scenario.paths), "invalid result documents size")
				} else {
					require.Len(t, docs, 0, "not empty result documents")
				}
			},
		)
	}
}

type testCollectionDocuments struct {
	name         string
	client       *firestore.Client
	raizelClient Client
	server       *mockServer
	collection   string
	err          error
	data         []map[string]*pb.Value
}

func (scenario *testCollectionDocuments) setup(t *testing.T) {
	c, srv := newMock(t)
	scenario.client = c
	scenario.server = srv
	raizelClient, err := newClient(scenario.client)
	require.Nil(t, err, "raizel client new error")
	scenario.raizelClient = raizelClient

	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
	} else {
		var (
			queryResponses = make([]interface{}, len(scenario.data))
			rootPath       = "projects/projectID/databases/(default)/documents"
		)
		for index, data := range scenario.data {
			readTime := aTimestamp
			if index%2 != 0 {
				readTime = aTimestamp2
			}
			queryResponses[index] = &pb.RunQueryResponse{
				Document: &pb.Document{
					Name:       fmt.Sprintf("%s/%s/mockref%d", rootPath, scenario.collection, index),
					CreateTime: aTimestamp,
					UpdateTime: aTimestamp,
					Fields:     data,
				},
				ReadTime: readTime,
			}
		}
		srv.addRPC(nil, queryResponses)
	}
}

func TestCollectionDocuments(test *testing.T) {
	scenarios := []testCollectionDocuments{
		{
			name:       "Get all collection documents",
			collection: "mockCol1",
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(35),
				},
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(55),
				},
			},
			err: nil,
		},
		{
			name:       "Returns a server error",
			collection: "mockEntity",
			data:       nil,
			err:        status.Error(codes.Unknown, "mockBadGateway"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				collection := scenario.raizelClient.Collection(scenario.collection)
				documents, err := collection.Documents(context.Background()).GetAll()
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
				if err == nil {
					require.NotZerof(t, documents, "documents response invalid: %+v", documents)
					require.NotEmpty(t, documents, "documents len invalid: %+v", documents)
				}
			},
		)
	}
}

type restriction struct {
	path  string
	op    string
	value interface{}
}

type queryOrder struct {
	path      string
	direction Direction
}

type queryScenario struct {
	collection   string
	restrictions []restriction
	order        queryOrder
	limit        int
	offset       int
}

type testQuery struct {
	name         string
	client       *firestore.Client
	raizelClient Client
	server       *mockServer
	query        queryScenario
	err          error
	data         []map[string]*pb.Value
}

func (scenario *testQuery) setup(t *testing.T) {
	c, srv := newMock(t)
	scenario.client = c
	scenario.server = srv
	raizelClient, err := newClient(scenario.client)
	require.Nil(t, err, "raizel client new error")
	scenario.raizelClient = raizelClient

	if scenario.err != nil {
		srv.addRPC(nil, scenario.err)
	} else {
		var (
			queryResponses = make([]interface{}, len(scenario.data))
			rootPath       = "projects/projectID/databases/(default)/documents"
		)
		for index, data := range scenario.data {
			readTime := aTimestamp
			if index%2 != 0 {
				readTime = aTimestamp2
			}
			queryResponses[index] = &pb.RunQueryResponse{
				Document: &pb.Document{
					Name:       fmt.Sprintf("%s/%s/mockref%d", rootPath, scenario.query.collection, index),
					CreateTime: aTimestamp,
					UpdateTime: aTimestamp,
					Fields:     data,
				},
				ReadTime: readTime,
			}
		}
		srv.addRPC(nil, queryResponses)
	}
}

func TestQuery(test *testing.T) {
	scenarios := []testQuery{
		{
			name: "Query documents",
			query: queryScenario{
				collection: "mockCol1",
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
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
			name: "Query documents with desc path order",
			query: queryScenario{
				collection: "mockCol1",
				order:      queryOrder{path: "name", direction: Desc},
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Query documents with asc path order",
			query: queryScenario{
				collection: "mockCol1",
				order:      queryOrder{path: "name", direction: Asc},
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Query documents with offset",
			query: queryScenario{
				collection: "mockCol1",
				offset:     100,
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Query documents with limit",
			query: queryScenario{
				collection: "mockCol1",
				limit:      100,
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Query documents with offset and limit",
			query: queryScenario{
				collection: "mockCol1",
				offset:     100,
				limit:      100,
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Query documents with orderby and offset and limit",
			query: queryScenario{
				collection: "mockCol1",
				offset:     100,
				limit:      100,
				order:      queryOrder{path: "name", direction: Asc},
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
			},
			data: []map[string]*pb.Value{
				{
					"id":   strval("#mockref3"),
					"mame": strval("Mock Three"),
					"age":  intval(25),
				},
				{
					"id":   strval("#mockref1"),
					"mame": strval("Mock One"),
					"age":  intval(5),
				},
				{
					"id":   strval("#mockref4"),
					"mame": strval("Mock Four"),
					"age":  intval(255),
				},
				{
					"id":   strval("#mockref2"),
					"mame": strval("Mock Two"),
					"age":  intval(20),
				},
			},
			err: nil,
		},
		{
			name: "Returns a server error",
			query: queryScenario{
				collection: "mockEntity",
				restrictions: []restriction{
					{path: "name", op: "==", value: "Mock One"},
					{path: "age", op: "==", value: int(25)},
				},
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
				var query Query = scenario.raizelClient.Collection(scenario.query.collection)
				for _, restriction := range scenario.query.restrictions {
					query = query.Where(restriction.path, restriction.op, restriction.value)
				}
				if scenario.query.order != (queryOrder{}) {
					query = query.OrderBy(scenario.query.order.path, scenario.query.order.direction)
				}
				if scenario.query.offset > 0 {
					query = query.Offset(scenario.query.offset)
				}
				if scenario.query.limit > 0 {
					query = query.Limit(scenario.query.limit)
				}
				documents, err := query.Documents(context.Background()).GetAll()
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
				if err == nil {
					require.NotZerof(t, documents, "documents response invalid: %+v", documents)
					require.NotEmpty(t, documents, "documents len invalid: %+v", documents)
				}
			},
		)
	}
}

type batchDocument struct {
	path string
	data map[string]interface{}
}

type batchCommand struct {
	method   string
	document batchDocument
}

type testBatch struct {
	name         string
	client       *firestore.Client
	raizelClient Client
	server       *mockServer
	commands     []batchCommand
	err          error
}

func (scenario *testBatch) setup(t *testing.T) {
	c, srv := newMock(t)
	scenario.client = c
	scenario.server = srv
	raizelClient, err := newClient(scenario.client)
	require.Nil(t, err, "client new error")
	scenario.raizelClient = raizelClient

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

func TestBatch(test *testing.T) {
	scenarios := []testBatch{
		{
			name: "Executes batch",
			commands: []batchCommand{
				{
					method: "Set",
					document: batchDocument{
						path: "mockentity/mock1",
						data: map[string]interface{}{
							"strfield":   "Mock One",
							"intfield":   23,
							"timefield":  time.Now().UTC(),
							"floatfield": 23.33,
							"boolfield":  true,
						},
					},
				},
				{
					method: "SetMerge",
					document: batchDocument{
						path: "mockentity/mock2",
						data: map[string]interface{}{
							"strfield":   "Mock Two",
							"intfield":   20,
							"timefield":  time.Now().UTC(),
							"floatfield": 20.33,
							"boolfield":  false,
						},
					},
				},
				{
					method: "Delete",
					document: batchDocument{
						path: "mockentity/mock1",
					},
				},
				{
					method: "Delete",
					document: batchDocument{
						path: "mockentity/mock2",
					},
				},
			},
			err: nil,
		},
		{
			name: "Returns server error",
			commands: []batchCommand{
				{
					method: "Delete",
					document: batchDocument{
						path: "mockentity/mock1",
					},
				},
				{
					method: "Delete",
					document: batchDocument{
						path: "mockentity/mock2",
					},
				},
			},
			err: status.Error(codes.InvalidArgument, "mockScenarioError"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				batch := scenario.raizelClient.Batch()
				for _, command := range scenario.commands {
					switch command.method {
					case "Set":
						newBatch := batch.Set(
							scenario.raizelClient.Doc(
								command.document.path,
							),
							command.document.data,
						)
						require.NotNil(t, newBatch, "invalid batch on set method result")
					case "SetMerge":
						newBatch := batch.Set(
							scenario.raizelClient.Doc(
								command.document.path,
							),
							command.document.data,
							MergeAll,
						)
						require.NotNil(t, newBatch, "invalid batch on setmerge method result")
					case "Delete":
						newBatch := batch.Delete(
							scenario.raizelClient.Doc(
								command.document.path,
							),
						)
						require.NotNil(t, newBatch, "invalid batch on delete method result")
					default:
						t.Errorf("batch invalid mehtod: method=%s", command.method)
					}
				}
				err := batch.Commit(context.Background())
				require.Equalf(t, grpc.Code(scenario.err), grpc.Code(err), "invalid grpccode: error=%+v", err)
				require.Equalf(t, grpc.ErrorDesc(scenario.err), grpc.ErrorDesc(err), "invalid grpcdesc: error=%v", err)
			},
		)
	}
}
