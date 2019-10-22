package firestore_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/raizel"
	"github.com/rjansen/raizel/firestore"
	fmock "github.com/rjansen/raizel/firestore/mock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID        string    `firestore:"id,omitempty"`
	Name      string    `firestore:"name,omitempty"`
	Age       int       `firestore:"age,omitempty"`
	CreatedAt time.Time `firestore:"createdAt,omitempty"`
	UpdatedAt time.Time `firestore:"updatedAt,omitempty"`
}

type testEntityKey struct {
	collection string
	name       string
	value      interface{}
}

func (k testEntityKey) Name() string {
	return k.name
}

func (k testEntityKey) Value() interface{} {
	return k.value
}

func (k testEntityKey) EntityName() string {
	return k.collection
}

func TestNewRepository(test *testing.T) {
	repository := firestore.NewRepository(nil)
	require.NotNil(test, repository, "invalid repository instance")
}

type testRepositoryGet struct {
	name   string
	ctx    context.Context
	ref    *fmock.DocumentRefMock
	doc    *fmock.DocumentSnapshotMock
	client *fmock.ClientMock
	key    raizel.EntityKey
	result raizel.Entity
	err    error
}

func (scenario *testRepositoryGet) setup(t *testing.T) {
	var (
		ref = fmock.NewDocumentRefMock()
		doc = fmock.NewDocumentSnapshotMock()
		cli = fmock.NewClientMock()
	)
	require.NotNil(t, ref, "mock docref instance")
	require.NotNil(t, doc, "mock doc instance")
	require.NotNil(t, cli, "mock client instance")

	doc.On("DataTo", mock.AnythingOfType("firestore_test.testEntity")).Return(nil)
	ref.On("Get", mock.Anything).Return(doc, scenario.err)
	cli.On("Doc", mock.AnythingOfType("string")).Return(ref)
	cli.On("Close", mock.Anything).Return(nil)

	scenario.ref = ref
	scenario.doc = doc
	scenario.client = cli
	scenario.ctx = context.Background()
}

func TestRepositoryGet(test *testing.T) {
	scenarios := []testRepositoryGet{
		{
			name: "Get entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
			result: testEntity{},
		},
		{
			name: "Error when try to Get an entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := firestore.NewRepository(scenario.client)
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.ctx, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
				repository.Close(scenario.ctx)
				scenario.client.AssertExpectations(t)
				scenario.ref.AssertExpectations(t)
				if scenario.err == nil {
					scenario.doc.AssertExpectations(t)
				} else {
					scenario.doc.AssertNotCalled(t, "DataTo", mock.Anything)
				}
			},
		)
	}
}

type testRepositorySet struct {
	name   string
	ctx    context.Context
	ref    *fmock.DocumentRefMock
	client *fmock.ClientMock
	key    raizel.EntityKey
	data   raizel.Entity
	err    error
}

func (scenario *testRepositorySet) setup(t *testing.T) {
	var (
		ref = fmock.NewDocumentRefMock()
		cli = fmock.NewClientMock()
	)
	require.NotNil(t, ref, "mock docref instance")
	require.NotNil(t, cli, "mock client instance")

	ref.On("Set", mock.Anything, mock.Anything, mock.AnythingOfType("[]firestore.SetOption")).Return(scenario.err)
	cli.On("Doc", mock.AnythingOfType("string")).Return(ref)
	cli.On("Close", mock.Anything).Return(nil)

	scenario.ref = ref
	scenario.client = cli
	scenario.ctx = context.Background()
}

func TestRepositorySet(test *testing.T) {
	scenarios := []testRepositorySet{
		{
			name: "Set entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
			data: testEntity{},
		},
		{
			name: "Error when try to Set an entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := firestore.NewRepository(scenario.client)
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.ctx, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
				repository.Close(scenario.ctx)
				scenario.client.AssertExpectations(t)
				scenario.ref.AssertExpectations(t)
			},
		)
	}
}

type testRepositoryDelete struct {
	name   string
	ctx    context.Context
	ref    *fmock.DocumentRefMock
	client *fmock.ClientMock
	key    raizel.EntityKey
	err    error
}

func (scenario *testRepositoryDelete) setup(t *testing.T) {
	var (
		ref = fmock.NewDocumentRefMock()
		cli = fmock.NewClientMock()
	)
	require.NotNil(t, ref, "mock docref instance")
	require.NotNil(t, cli, "mock client instance")

	ref.On("Delete", mock.Anything).Return(scenario.err)
	cli.On("Doc", mock.AnythingOfType("string")).Return(ref)
	cli.On("Close").Return(nil)

	scenario.ref = ref
	scenario.client = cli
	scenario.ctx = context.Background()
}

func TestRepositoryDelete(test *testing.T) {
	scenarios := []testRepositoryDelete{
		{
			name: "Delete entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
		},
		{
			name: "Error when try to Delete an entity",
			key: testEntityKey{
				collection: "mymockcollection",
				name:       "id",
				value:      "identifier",
			},
			err: errors.New("errMock"),
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				repository := firestore.NewRepository(scenario.client)
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.ctx, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
				repository.Close(scenario.ctx)
				scenario.client.AssertExpectations(t)
				scenario.ref.AssertExpectations(t)
			},
		)
	}
}
