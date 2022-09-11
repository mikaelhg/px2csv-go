package ast_test

import (
	"os"

	"github.com/mikaelhg/gpcaxis/ast"
	"gotest.tools/v3/assert"

	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/repr"
)

func parseFile(t *testing.T, filename string) {
	r, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	header, err := ast.PxParser.Parse("", r, participle.AllowTrailing(true))
	repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}
	assert.Check(t, header != nil)
}

func TestPxFileHeader(t *testing.T) {
	parseFile(t, "../data/statfin_ehk_pxt_005_en.px")
	parseFile(t, "../data/010_kats_tau_101.px")
}

func TestTerminate(t *testing.T) {
	text := `A=1;
VALUENOTE[en]("Information","Median mileage")="Median mileage";
DATA=
1564581 174000 162000 21 1243095 321486 
`
	header, err := ast.PxParser.ParseString("", text, participle.AllowTrailing(true))
	// repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}
	assert.Check(t, header != nil)
}

func TestTerminate2(t *testing.T) {
	text := `A=1;
VALUES("Data")="Annual change %","Quartal change %","Value, M";
PRECISION("Data","Value, M")=1;
DATA=
"." "." 325.3 
60.078759 "." 520.8 
`
	header, err := ast.PxParser.ParseString("", text, participle.AllowTrailing(true))
	// repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}
	assert.Check(t, header != nil)
}
