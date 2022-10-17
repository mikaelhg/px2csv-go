package internal

import (
	"bufio"
	"strings"

	"github.com/apache/arrow/go/v9/parquet"
	"github.com/apache/arrow/go/v9/parquet/file"
	"github.com/apache/arrow/go/v9/parquet/schema"
)

type StatCubeWriter interface {
	WriteHeading(stub, headingCsv []string)

	// Yes, the signature is a bit funny, but it's like this to optimize
	// the way data is laid out in a tight loop that has to avoid allocs.
	WriteRow(stubs *[]*string, values *[][]byte,
		valueLengths *[]int, stubWidth, headingWidth int)
}

type StatCubeCsvWriter struct {
	Writer *bufio.Writer
}

type StatCubeParquetWriter struct {
	Writer *bufio.Writer
}

func (w *StatCubeCsvWriter) WriteHeading(stub, headingCsv []string) {
	w.Writer.WriteString("\"")
	w.Writer.WriteString(strings.Join(stub, "\";\""))
	w.Writer.WriteString("\";\"")
	w.Writer.WriteString(strings.Join(headingCsv, "\";\""))
	w.Writer.WriteString("\"\n")
}

func (w *StatCubeCsvWriter) WriteRow(stubs *[]*string, values *[][]byte,
	valueLengths *[]int, stubWidth, headingWidth int) {
	w.Writer.WriteByte('"')
	for i, s := range *stubs {
		w.Writer.WriteString(*s)
		if i < stubWidth-1 {
			w.Writer.WriteByte('"')
			w.Writer.WriteByte(';')
			w.Writer.WriteByte('"')
		}
	}
	w.Writer.WriteByte('"')
	w.Writer.WriteByte(';')
	for i, s := range *values {
		w.Writer.Write(s[0:(*valueLengths)[i]])
		if i < headingWidth-1 {
			w.Writer.WriteByte(';')
		}
	}
	w.Writer.WriteByte('\n')
}

func (w *StatCubeParquetWriter) WriteHeading(stub, headingCsv []string) {
	const (
		valueCount = 10000
		pageSize   = 16384
	)
	var (
		props = parquet.NewWriterProperties(
			parquet.WithDictionaryDefault(false), parquet.WithDataPageSize(pageSize))
		fieldList = schema.FieldList{
			schema.NewFloat32Node("col", parquet.Repetitions.Required, -1),
		}
		sc, _ = schema.NewGroupNode("schema", parquet.Repetitions.Required, fieldList, -1)
	)

	writer := file.NewParquetWriter(w.Writer, sc, file.WithWriterProps(props))
	rgWriter := writer.AppendBufferedRowGroup()
	cwr, _ := rgWriter.Column(0)
	cw := cwr.(*file.Float32ColumnChunkWriter)
	valuesIn := make([]float32, 0, valueCount)
	for i := 0; i < valueCount; i++ {
		valuesIn = append(valuesIn, float32((i%100)+1))
	}
	cw.WriteBatch(valuesIn, nil, nil)
	rgWriter.Close()
	writer.Close()
}

func (w *StatCubeParquetWriter) WriteRow(stubs *[]*string, values *[][]byte,
	valueLengths *[]int, stubWidth, headingWidth int) {
}
