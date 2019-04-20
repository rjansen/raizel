package spanner

import (
	"cloud.google.com/go/spanner"
)

func Delete(table string, ks KeySet) *Mutation {
	return spanner.Delete(table, ks)
}

func Insert(table string, cols []string, vals []interface{}) *Mutation {
	return spanner.Insert(table, cols, vals)
}

func InsertMap(table string, in map[string]interface{}) *Mutation {
	return spanner.InsertMap(table, in)
}

func InsertOrUpdate(table string, cols []string, vals []interface{}) *Mutation {
	return spanner.InsertOrUpdate(table, cols, vals)
}

func InsertOrUpdateMap(table string, in map[string]interface{}) *Mutation {
	return spanner.InsertOrUpdateMap(table, in)
}

func InsertOrUpdateStruct(table string, in interface{}) (*Mutation, error) {
	return spanner.InsertOrUpdateStruct(table, in)
}

func InsertStruct(table string, in interface{}) (*Mutation, error) {
	return spanner.InsertStruct(table, in)
}

func Replace(table string, cols []string, vals []interface{}) *Mutation {
	return spanner.Replace(table, cols, vals)
}

func ReplaceMap(table string, in map[string]interface{}) *Mutation {
	return spanner.ReplaceMap(table, in)
}

func ReplaceStruct(table string, in interface{}) (*Mutation, error) {
	return spanner.ReplaceStruct(table, in)
}

func Update(table string, cols []string, vals []interface{}) *Mutation {
	return spanner.Update(table, cols, vals)
}

func UpdateMap(table string, in map[string]interface{}) *Mutation {
	return spanner.UpdateMap(table, in)
}

func UpdateStruct(table string, in interface{}) (*Mutation, error) {
	return spanner.UpdateStruct(table, in)
}
