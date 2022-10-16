package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	pxFilename := flag.String("px", "", "PC-AXIS input file")
	csvFilename := flag.String("csv", "", "CSV output file")
	flag.Parse()

	inf, err := os.Open(*pxFilename)
	if err != nil {
		panic(err)
	}
	defer inf.Close()

	outf, err := os.OpenFile(*csvFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	reader, writer := bufio.NewReader(inf), bufio.NewWriter(outf)

	pxParser := internal.Parser{}
	pxParser.ParseHeader(reader)
	pxParser.ParseDataDense(reader, writer)
	writer.Flush()
}
