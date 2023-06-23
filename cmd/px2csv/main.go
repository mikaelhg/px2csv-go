package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	inputFilename := flag.String("px", "/dev/stdin", "PC-AXIS input file")
	outputFilename := flag.String("csv", "/dev/stdout", "CSV output file")
	flag.Parse()

	inf, err := os.Open(*inputFilename)
	if err != nil {
		panic(err)
	}
	defer inf.Close()

	outf, err := os.OpenFile(*outputFilename, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer outf.Close()

	reader, writer := bufio.NewReader(inf), bufio.NewWriter(outf)
	cubeWriter := internal.StatCubeCsvWriter{Writer: writer}
	pxParser := internal.PxParser{CubeWriter: &cubeWriter}
	pxParser.ParseHeader(reader)
	pxParser.ParseDataDense(reader)
	err = writer.Flush()
	if err != nil {
		panic(err)
	}
}
