// +build integration

package firestore

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rjansen/raizel"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/require"
)

const (
	testProjectID  = "e-pedion"
	testCollection = "environments/test/entity_mock"
)

type testRepositoryFirestoreGet struct {
	name     string
	tree     yggdrasil.Tree
	mockData *entityMock
	key      entityKeyMock
	result   raizel.Entity
	err      error
}

func (scenario *testRepositoryFirestoreGet) setup(t *testing.T) {
	var (
		fclient, errFclient = newFirestoreClient(testProjectID)
		client, errClient   = newClient(fclient)
		roots               = yggdrasil.NewRoots()
		err                 = Register(&roots, client)
	)
	require.Nil(t, errFclient, "new firestore error")
	require.Nil(t, errClient, "new client error")
	require.Nil(t, err, "register client err")
	require.NotNil(t, roots, "roots instance")
	require.NotNil(t, client, "client instance")

	if scenario.err == nil {
		currentTime := time.Now().UTC()
		scenario.mockData = &entityMock{
			ID:        newUUID(),
			Name:      scenario.name,
			Age:       777,
			Data:      dynamicData{"key": "value"},
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
		}
		var (
			doc = client.Doc(entityMockRef(testCollection, scenario.mockData.ID))
			err = doc.Set(context.Background(), scenario.mockData)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryFirestoreGet) tearDown(t *testing.T) {
	if scenario.mockData != nil {
		var (
			client = MustReference(scenario.tree)
			err    = client.Doc(
				entityMockRef(testCollection, scenario.mockData.ID),
			).Delete(context.Background())
		)
		require.Nil(t, err, "teardown error")
	}
	scenario.tree.Close()
}

func TestRepositoryFirestoreGet(test *testing.T) {
	scenarios := []testRepositoryFirestoreGet{
		{
			name: "When try to Get an entity on database returns it successfully",
			key: entityKeyMock{
				collection: testCollection,
				name:       "id",
			},
			result: &entityMock{},
		},
		{
			name: "When try to Get an entity on database but an error is returned",
			key: entityKeyMock{
				collection: testCollection,
				name:       "id",
				value:      newUUID(),
			},
			result: &entityMock{},
			err:    raizel.ErrNotFound,
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository()
				require.NotNil(t, repository, "repository instance")
				err := repository.Get(scenario.tree, scenario.key, scenario.result)
				require.Equal(t, scenario.err, err, "get error")
			},
		)
	}
}

type testRepositoryFirestoreSet struct {
	name     string
	tree     yggdrasil.Tree
	mockData *entityMock
	key      entityKeyMock
	data     *entityMock
	err      error
}

func (scenario *testRepositoryFirestoreSet) setup(t *testing.T) {
	var (
		fclient, errFclient = newFirestoreClient(testProjectID)
		client, errClient   = newClient(fclient)
		roots               = yggdrasil.NewRoots()
		err                 = Register(&roots, client)
	)
	require.Nil(t, errFclient, "new firestore error")
	require.Nil(t, errClient, "new client error")
	require.Nil(t, err, "register client err")
	require.NotNil(t, roots, "roots instance")
	require.NotNil(t, client, "client instance")

	if scenario.mockData != nil {
		var (
			doc = client.Doc(entityMockRef(testCollection, scenario.mockData.ID))
			err = doc.Set(context.Background(), scenario.mockData)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
		scenario.data.ID = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryFirestoreSet) tearDown(t *testing.T) {
	client := MustReference(scenario.tree)
	if scenario.mockData != nil {
		err := client.Doc(
			entityMockRef(testCollection, scenario.mockData.ID),
		).Delete(context.Background())
		require.Nil(t, err, "teardown mock error")
	}
	if scenario.data != nil {
		err := client.Doc(
			entityMockRef(testCollection, scenario.data.ID),
		).Delete(context.Background())
		require.Nil(t, err, "teardown data error")
	}
	scenario.tree.Close()
}

func TestRepositoryFirestoreSet(test *testing.T) {
	currentTime := time.Now().UTC()
	uuid := newUUID()
	scenarios := []testRepositoryFirestoreSet{
		{
			name: "When try to Set a new entity on database successfully",
			key: entityKeyMock{
				collection: testCollection,
				name:       "id",
				value:      uuid,
			},
			data: &entityMock{
				ID:   uuid,
				Name: "When try to Set a new entity on database successfully",
				Age:  10,
				Data: map[string]interface{}{
					"golangkey": "golangvalue",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		},
		{
			name: "When try to Set an existent entity on database successfully",
			key: entityKeyMock{
				collection: testCollection,
				name:       "id",
				value:      uuid,
			},
			mockData: &entityMock{
				ID:   uuid,
				Name: "When try to Set an existent entity on database successfully",
				Age:  66,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
			data: &entityMock{
				ID:   uuid,
				Name: "Changed to: 'When try to Set an existent entity on database successfully'",
				Age:  333,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
					"changedkey":  "1",
				},
				Deleted:   false,
				CreatedAt: currentTime,
				UpdatedAt: currentTime,
			},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository()
				require.NotNil(t, repository, "repository instance")
				err := repository.Set(scenario.tree, scenario.key, scenario.data)
				require.Equal(t, scenario.err, err, "set error")
			},
		)
	}
}

type testRepositoryFirestoreDelete struct {
	name     string
	tree     yggdrasil.Tree
	mockData *entityMock
	key      entityKeyMock
	err      error
}

func (scenario *testRepositoryFirestoreDelete) setup(t *testing.T) {
	var (
		fclient, errFclient = newFirestoreClient(testProjectID)
		client, errClient   = newClient(fclient)
		roots               = yggdrasil.NewRoots()
		err                 = Register(&roots, client)
	)
	require.Nil(t, errFclient, "new firestore error")
	require.Nil(t, errClient, "new client error")
	require.Nil(t, err, "register client err")
	require.NotNil(t, roots, "roots instance")
	require.NotNil(t, client, "client instance")

	if scenario.mockData != nil {
		var (
			doc = client.Doc(entityMockRef(testCollection, scenario.mockData.ID))
			err = doc.Set(context.Background(), scenario.mockData)
		)
		require.Nil(t, err, "setup data error")
		scenario.key.value = scenario.mockData.ID
	}

	scenario.tree = roots.NewTreeDefault()
}

func (scenario *testRepositoryFirestoreDelete) tearDown(t *testing.T) {
	if scenario.mockData != nil {
		var (
			client = MustReference(scenario.tree)
			err    = client.Doc(
				entityMockRef(testCollection, scenario.mockData.ID),
			).Delete(context.Background())
		)
		require.Nil(t, err, "teardown error")
	}
	scenario.tree.Close()
}

func TestRepositoryFirestoreDelete(test *testing.T) {
	scenarios := []testRepositoryFirestoreDelete{
		{
			name: "When try to Delete an entity on database successfully",
			key: entityKeyMock{
				collection: testCollection,
				name:       "id",
			},
			mockData: &entityMock{
				ID:   newUUID(),
				Name: "When try to Delete an entity on database successfully",
				Age:  96,
				Data: map[string]interface{}{
					"golangkey":   "golangvalue",
					"existentkey": "existentvalue",
					"deletedkey":  "deletedvalue",
				},
			},
		},
	}
	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)
				defer scenario.tearDown(t)

				repository := NewRepository()
				require.NotNil(t, repository, "repository instance")
				err := repository.Delete(scenario.tree, scenario.key)
				require.Equal(t, scenario.err, err, "set error")
			},
		)
	}
}
