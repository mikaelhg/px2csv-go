package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	filename := flag.String("file", "", "PX file")
	flag.Parse()
	f, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	pxParser := internal.NewParser()
	pxParser.ParseHeader(reader)
	pxParser.ParseDataDense(reader)
}
