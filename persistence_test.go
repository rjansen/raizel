package raizel

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitGetPoolErr(t *testing.T) {
	pool, err := GetPool()
	assert.NotNil(t, err)
	assert.Nil(t, pool)
}

func TestUnitSetupErr(t *testing.T) {
	err := Setup(nil)
	assert.NotNil(t, err)
}

func TestUnitSetupSuccess(t *testing.T) {
	mockClient := NewClientMock()
	mockClient.On("Close").Return(nil)
	mockPool := NewClientPoolMock()
	mockPool.On("Get").Return(mockClient, nil)
	mockPool.On("Close").Return(nil)
	err := Setup(mockPool)
	assert.Nil(t, err)
}

func TestUnitGetPool(t *testing.T) {
	pool, err := GetPool()
	assert.Nil(t, err)
	assert.NotNil(t, pool)
}

func TestUnitConfiguration(t *testing.T) {
	provider := "mockProvider"
	cfg := &Configuration{
		Provider: provider,
	}
	cfgStr := cfg.String()
	assert.Contains(t, cfgStr, provider)
}

func TestUnitSetGetClientOnContext(t *testing.T) {
	c := context.Background()
	c, err := SetClient(c)
	assert.Nil(t, err)
	assert.NotZero(t, c)

	client, err := GetClient(c)
	assert.Nil(t, err)
	assert.NotZero(t, client)
}

func TestUnitSetGetClientOnContextErr(t *testing.T) {
	c, err := SetClient(nil)
	assert.NotNil(t, err)
	assert.Zero(t, c)

	client, err := GetClient(c)
	assert.NotNil(t, err)
	assert.Zero(t, client)

	c = context.Background()
	client, err = GetClient(c)
	assert.NotNil(t, err)
	assert.Zero(t, client)
}

func TestUnitExecuteContext(t *testing.T) {
	err := ExecuteContext(
		func(c context.Context) error {
			assert.NotNil(t, c)
			client, err := GetClient(c)
			assert.Nil(t, err)
			assert.NotNil(t, client)
			return nil
		},
	)
	assert.Nil(t, err)
}

func TestUnitExecuteContextErr(t *testing.T) {
	err := ExecuteContext(
		func(c context.Context) error {
			assert.NotNil(t, c)
			client, err := GetClient(c)
			assert.Nil(t, err)
			assert.NotNil(t, client)
			return errors.New("MockExecuteContextErr")
		},
	)
	assert.NotNil(t, err)
}

func TestUnitExecute(t *testing.T) {
	err := Execute(
		func(c Client) error {
			assert.NotNil(t, c)
			return nil
		},
	)
	assert.Nil(t, err)
}

func TestUnitExecuteClientErr(t *testing.T) {
	err := Execute(
		func(c Client) error {
			assert.NotNil(t, c)
			return errors.New("MockExecuteErr")
		},
	)
	assert.NotNil(t, err)
}

func TestUnitExecuteWith(t *testing.T) {
	err := ExecuteWith(
		func(c Client, args ...interface{}) error {
			assert.NotNil(t, c)
			assert.NotNil(t, args)
			assert.Len(t, args, 3, "The size os arguments is not 3")
			return nil
		},
		1,
		"arg2",
		3,
	)
	assert.Nil(t, err)
}

func TestUnitExecuteWithWithoutArgs(t *testing.T) {
	err := ExecuteWith(
		func(c Client, args ...interface{}) error {
			assert.NotNil(t, c)
			assert.Nil(t, args)
			assert.Len(t, args, 0, "The size os arguments is not 0")
			return nil
		},
	)
	assert.Nil(t, err)
}

func TestUnitExecuteWithClientErr(t *testing.T) {
	err := ExecuteWith(
		func(c Client, args ...interface{}) error {
			assert.NotNil(t, c)
			assert.NotNil(t, args)
			assert.Len(t, args, 2, "The size os arguments is not 2")
			return errors.New("MockExecuteErr")
		},
		"arg1",
		2,
	)
	assert.NotNil(t, err)
}
