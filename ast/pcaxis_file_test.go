package ast_test

import (
	"os"

	"github.com/mikaelhg/gpcaxis/ast"

	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
)

func TestPxFileHeader(t *testing.T) {
	r, err := os.Open("../data/010_kats_tau_101.px")
	if err != nil {
		panic(err)
	}
	defer r.Close()
	header, err := ast.PxParser.Parse("", r, participle.AllowTrailing(true))
	repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}
}

func TestTerminate(t *testing.T) {
	text := `VALUENOTE[en]("Information","Median mileage")="Median mileage";
DATA=
1564581 174000 162000 21 1243095 321486 
`
	header, err := ast.PxParser.ParseString("", text, participle.AllowTrailing(true))
	repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}

	t.Error("yes")
}
