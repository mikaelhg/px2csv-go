package internal

import (
	"bufio"

	"github.com/apache/arrow/go/v9/parquet"
	"github.com/apache/arrow/go/v9/parquet/file"
	"github.com/apache/arrow/go/v9/parquet/schema"
)

type StatCubeParquetWriter struct {
	counter       int32
	Values        [][]float32
	Writer        *bufio.Writer
	ParquetWriter *file.Writer
}

func (w *StatCubeParquetWriter) WriteHeading(stub []string, headingFlattened [][]string) {
	const (
		pageSize = 16384
	)
	var (
		props = parquet.NewWriterProperties(
			parquet.WithDictionaryDefault(false), parquet.WithDataPageSize(pageSize))
		fieldList = schema.FieldList{
			schema.NewFloat32Node("col", parquet.Repetitions.Required, -1),
		}
		sc, _ = schema.NewGroupNode("schema", parquet.Repetitions.Required, fieldList, -1)
	)

	w.ParquetWriter = file.NewParquetWriter(w.Writer, sc, file.WithWriterProps(props))
}

func (w *StatCubeParquetWriter) writeBatch() {
	const (
		valueCount = 10000
	)
	rgWriter := w.ParquetWriter.AppendBufferedRowGroup()
	cwr, _ := rgWriter.Column(0)
	cw := cwr.(*file.Float32ColumnChunkWriter)
	valuesIn := make([]float32, 0, valueCount)
	for i := 0; i < valueCount; i++ {
		valuesIn = append(valuesIn, float32((i%100)+1))
	}
	cw.WriteBatch(valuesIn, nil, nil)
	rgWriter.Close()
}

func (w *StatCubeParquetWriter) WriteFooting() {
	w.ParquetWriter.Close()
}

func (w *StatCubeParquetWriter) WriteRow(stubs *[]*string, values *[][]byte,
	valueLengths *[]int, stubWidth, headingWidth int) {

	w.counter += 1
	if w.counter > 10000 {
		w.counter = 1
		w.writeBatch()
	}
}
