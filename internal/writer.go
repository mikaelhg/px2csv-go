package internal

import (
	"bufio"
	"strings"
)

type StatCubeWriter interface {
	writeHeading(stub, headingCsv []string)
	writeRow(stubs *[]*string, values *[][]byte,
		valueLengths *[]int, stubWidth, headingWidth int)
}

type StatCubeCsvWriter struct {
	writer *bufio.Writer
}

func NewStatCubeCsvWriter(writer *bufio.Writer) StatCubeCsvWriter {
	return StatCubeCsvWriter{writer: writer}
}

func (w StatCubeCsvWriter) writeHeading(stub, headingCsv []string) {
	w.writer.WriteString("\"")
	w.writer.WriteString(strings.Join(stub, "\";\""))
	w.writer.WriteString("\";\"")
	w.writer.WriteString(strings.Join(headingCsv, "\";\""))
	w.writer.WriteString("\"\n")
}

func (w StatCubeCsvWriter) writeRow(stubs *[]*string, values *[][]byte,
	valueLengths *[]int, stubWidth, headingWidth int) {
	w.writer.WriteByte('"')
	for i, s := range *stubs {
		w.writer.WriteString(*s)
		if i < stubWidth-1 {
			w.writer.WriteByte('"')
			w.writer.WriteByte(';')
			w.writer.WriteByte('"')
		}
	}
	w.writer.WriteByte('"')
	w.writer.WriteByte(';')
	for i, s := range *values {
		w.writer.Write(s[0:(*valueLengths)[i]])
		if i < headingWidth-1 {
			w.writer.WriteByte(';')
		}
	}
	w.writer.WriteByte('\n')
}
