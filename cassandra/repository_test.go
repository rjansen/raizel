package cassandra

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testEntity struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Age       int       `db:"age"`
	CreatedAt time.Time `db:"createdAt"`
	UpdatedAt time.Time `db:"updatedAt"`
}

type testEntityKey struct {
	pkField string
	id      string
}

func (k testEntityKey) GetKeyValue() interface{} {
	return k.id
}

func (k testEntityKey) GetEntityName() string {
	return k.pkField
}

func TestNewRepository(test *testing.T) {
	repository := NewRepository()
	require.NotNil(test, repository, "invalid repository instance")
}
