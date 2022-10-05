package main

import (
	"bufio"
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	pxParser := internal.NewParser()
	pxParser.ParseHeader(reader)
	pxParser.ParseDataDense(reader)
}
