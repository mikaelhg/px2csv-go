package main

import (
	pf "github.com/apache/arrow/go/v9/parquet/file"
)

func FooTest() {
	pf.OpenParquetFile("foo.parquet", false)
}
