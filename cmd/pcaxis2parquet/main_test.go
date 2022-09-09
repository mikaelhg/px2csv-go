package main

import (
	"github.com/mikaelhg/gpcaxis/ast"

	"testing"

	"github.com/alecthomas/repr"
)

func TestHello(t *testing.T) {
	ini, err := ast.Parser.ParseString("", "foo=\"bar\"")
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
