package spanner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteMutation(t *testing.T) {
	mutations := Delete("MockTable", Key{1, 2})
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestInsertMutation(t *testing.T) {
	mutations := Insert(
		"MockTable",
		[]string{"Column1", "Column2", "ColumnsN"},
		[]interface{}{1, 2, "N"},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestInsertMapMutation(t *testing.T) {
	mutations := InsertMap(
		"MockTable",
		map[string]interface{}{
			"Column1": 1,
			"Column2": 2,
			"ColumnN": "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestInsertOrUpdateMutation(t *testing.T) {
	mutations := InsertOrUpdate(
		"MockTable",
		[]string{"Column1", "Column2", "ColumnsN"},
		[]interface{}{1, 2, "N"},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestInsertOrUpdateMapMutation(t *testing.T) {
	mutations := InsertOrUpdateMap(
		"MockTable",
		map[string]interface{}{
			"Column1": 1,
			"Column2": 2,
			"ColumnN": "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

type data struct {
	Column1, Column2 int
	ColumnN          string
}

func TestInsertOrUpdateStructMutation(t *testing.T) {
	mutations, err := InsertOrUpdateStruct(
		"MockTable",
		data{
			Column1: 1,
			Column2: 2,
			ColumnN: "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
	require.Nil(t, err, "mutations error")
}

func TestInsertStructMutation(t *testing.T) {
	mutations, err := InsertStruct(
		"MockTable",
		data{
			Column1: 1,
			Column2: 2,
			ColumnN: "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
	require.Nil(t, err, "mutations error")
}

func TestReplaceMutation(t *testing.T) {
	mutations := Replace(
		"MockTable",
		[]string{"Column1", "Column2", "ColumnsN"},
		[]interface{}{1, 2, "N"},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestReplaceMapMutation(t *testing.T) {
	mutations := ReplaceMap(
		"MockTable",
		map[string]interface{}{
			"Column1": 1,
			"Column2": 2,
			"ColumnN": "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestReplaceStructMutation(t *testing.T) {
	mutations, err := ReplaceStruct(
		"MockTable",
		data{
			Column1: 1,
			Column2: 2,
			ColumnN: "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
	require.Nil(t, err, "mutations error")
}

func TestUpdateMutation(t *testing.T) {
	mutations := Update(
		"MockTable",
		[]string{"Column1", "Column2", "ColumnsN"},
		[]interface{}{1, 2, "N"},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestUpdateMapMutation(t *testing.T) {
	mutations := UpdateMap(
		"MockTable",
		map[string]interface{}{
			"Column1": 1,
			"Column2": 2,
			"ColumnN": "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
}

func TestUpdateStructMutation(t *testing.T) {
	mutations, err := UpdateStruct(
		"MockTable",
		data{
			Column1: 1,
			Column2: 2,
			ColumnN: "N",
		},
	)
	require.NotNil(t, mutations, "invalid mutations instance")
	require.Nil(t, err, "mutations error")
}
