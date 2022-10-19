package internal

import (
	"bufio"
	"strings"

	"github.com/apache/arrow/go/v9/parquet"
	"github.com/apache/arrow/go/v9/parquet/file"
	"github.com/apache/arrow/go/v9/parquet/schema"
)

type StatCubeWriter interface {
	WriteHeading(stub []string, headingFlattened [][]string)

	// Yes, the signature is a bit funny, but it's like this to optimize
	// the way data is laid out in a tight loop that has to avoid allocs.
	WriteRow(stubs *[]*string, values *[][]byte,
		valueLengths *[]int, stubWidth, headingWidth int)

	WriteFooting()
}

type StatCubeCsvWriter struct {
	Writer *bufio.Writer
}

func joinStringSlice(ss []string) string {
	return strings.Join(ss, " ")
}

func (w *StatCubeCsvWriter) WriteHeading(stub []string, headingFlattened [][]string) {
	headingCsv := MapXtoY(headingFlattened, joinStringSlice)
	w.Writer.WriteString("\"")
	w.Writer.WriteString(strings.Join(stub, "\";\""))
	w.Writer.WriteString("\";\"")
	w.Writer.WriteString(strings.Join(headingCsv, "\";\""))
	w.Writer.WriteString("\"\n")
}

func (w *StatCubeCsvWriter) WriteFooting() {
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
