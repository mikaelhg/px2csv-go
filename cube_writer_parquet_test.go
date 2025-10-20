package internal

import (
	"testing"

	"github.com/parquet-go/parquet-go"
)

func TestGroups(t *testing.T) {

	// https://github.com/parquet-go/parquet-go/blob/main/dictionary_test.go#L202

	node := parquet.String()

	g := parquet.Group{}
	g["test_string"] = node

	schema := parquet.NewSchema("test_schema", g)

	t.Logf("schema: %v\n", schema)

}
