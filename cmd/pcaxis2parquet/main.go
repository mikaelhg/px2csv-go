package main

import (
	"os"

	"github.com/mikaelhg/gpcaxis/ast"

	"github.com/alecthomas/repr"
)

func main() {
	ini, err := ast.Parser.Parse("", os.Stdin)
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
