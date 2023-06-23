package internal

import (
	"bufio"
	"bytes"

	parquet "github.com/segmentio/parquet-go"
)

type StatCubeParquetWriter struct {
	counter int32
	Values  [][]float32
	Writer  *bufio.Writer
	Storage bytes.Buffer
}

func ParquetWriter() *StatCubeParquetWriter {
	return &StatCubeParquetWriter{
		counter: 0,
	}
}

func (w *StatCubeParquetWriter) WriteHeading(stub []string, headingFlattened [][]string) {
	node := parquet.String()
	node = parquet.Encoded(node, &parquet.RLEDictionary)

	g := parquet.Group{}
	g["mystring"] = node

	schema := parquet.NewSchema("test", g)

	rows := []parquet.Row{[]parquet.Value{parquet.ValueOf("hello").Level(0, 0, 0)}}

	_, _ = schema, rows
}

func (w *StatCubeParquetWriter) writeBatch() {
}

func (w *StatCubeParquetWriter) WriteFooting() {
}

func (w *StatCubeParquetWriter) WriteRow(stubs *[]*string, buffer *[]byte,
	valueLengths *[]int, stubWidth, headingWidth int) {

	w.counter += 1
	if w.counter > 10000 {
		w.counter = 1
		w.writeBatch()
	}
}
