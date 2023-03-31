package internal

import (
	"bufio"
	"strings"

	"github.com/samber/lo"
)

type StatCubeWriter interface {
	WriteHeading(stub []string, headingFlattened [][]string)

	// Yes, the signature is a bit funny, but it's like this to optimize
	// the way data is laid out in a tight loop that has to avoid allocs.
	WriteRow(stubs *[]*string, buffer *[]byte,
		valueLengths *[]int, stubWidth, headingWidth int)

	WriteFooting()
}

type StatCubeCsvWriter struct {
	Writer *bufio.Writer
}

func (w *StatCubeCsvWriter) WriteHeading(stub []string, headingFlattened [][]string) {
	headingCsv := lo.Map(headingFlattened, joinStringSlice)
	w.Writer.WriteString("\"")
	w.Writer.WriteString(strings.Join(stub, "\";\""))
	w.Writer.WriteString("\";\"")
	w.Writer.WriteString(strings.Join(headingCsv, "\";\""))
	w.Writer.WriteString("\"\n")
}

func (w *StatCubeCsvWriter) WriteFooting() {
	// NOP
}

func (w *StatCubeCsvWriter) WriteRow(stubs *[]*string, buffer *[]byte,
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
	for i := 0; i < headingWidth; i++ {
		offset := DataValueWidth * i
		w.Writer.Write((*buffer)[offset : offset+(*valueLengths)[i]])
		if i < headingWidth-1 {
			w.Writer.WriteByte(';')
		}
	}
	w.Writer.WriteByte('\n')
}
