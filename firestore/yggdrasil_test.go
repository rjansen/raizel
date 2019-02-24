package firestore_test

import (
	"fmt"
	"testing"

	"github.com/rjansen/raizel/firestore"
	fmock "github.com/rjansen/raizel/firestore/mock"
	"github.com/rjansen/yggdrasil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRegister struct {
	name   string
	client firestore.Client
	err    error
}

func TestRegister(test *testing.T) {
	scenarios := []testRegister{
		{
			name:   "Register the Client reference",
			client: fmock.NewClientMock(),
		},
		{
			name:   "Register a nil Client reference",
			client: nil,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				roots := yggdrasil.NewRoots()
				err := firestore.Register(&roots, scenario.client)
				assert.Equal(t, scenario.err, err)

				tree := roots.NewTreeDefault()
				client, err := firestore.Reference(tree)

				require.Nil(t, err, "tree reference error")
				require.Exactly(t, scenario.client, client, "client reference")
			},
		)
	}
}

type testReference struct {
	name       string
	references map[yggdrasil.Path]yggdrasil.Reference
	tree       yggdrasil.Tree
	err        error
}

func (scenario *testReference) setup(t *testing.T) {
	roots := yggdrasil.NewRoots()
	for path, reference := range scenario.references {
		err := roots.Register(path, reference)
		assert.Nil(t, err, "register error")
	}
	scenario.tree = roots.NewTreeDefault()
}

func TestReference(test *testing.T) {
	scenarios := []testReference{
		{
			name: "Access the Client Reference",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				yggdrasil.NewPath("/raizel/firestore/client"): yggdrasil.NewReference(fmock.NewClientMock()),
			},
		},
		{
			name: "Access a nil Client Reference",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				yggdrasil.NewPath("/raizel/firestore/client"): yggdrasil.NewReference(nil),
			},
		},
		{
			name: "When Client was not register returns path not found",
			err:  yggdrasil.ErrPathNotFound,
		},
		{
			name: "When a invalid Client was register returns invalid reference error",
			references: map[yggdrasil.Path]yggdrasil.Reference{
				yggdrasil.NewPath("/raizel/firestore/client"): yggdrasil.NewReference(new(struct{})),
			},
			err: firestore.ErrInvalidReference,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				_, err := firestore.Reference(scenario.tree)
				assert.Equal(t, scenario.err, err, "reference error")
				if scenario.err != nil {
					assert.PanicsWithValue(t, scenario.err,
						func() {
							_ = firestore.MustReference(scenario.tree)
						},
					)
				} else {
					assert.NotPanics(t,
						func() {
							_ = firestore.MustReference(scenario.tree)
						},
					)
				}
			},
		)
	}
}
