package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDataMapper(t *testing.T) {
	t.Run(
		"new_entity_mapper",
		func(t *testing.T) {
			// dataMapper, err := NewDataMapper((*dataMock)(nil))
			dataMapper, err := NewDataMapper[dataMock]()

			assert.Nil(t, err)
			assert.NotNil(t, dataMapper)

			assert.Len(t, dataMapper.Fields.Scalars, 7)
			assert.Len(t, dataMapper.Fields.Nesteds, 1)
			assert.Len(t, dataMapper.Fields.Slices, 1)
		},
	)
}
