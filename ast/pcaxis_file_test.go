package ast_test

import (
	"fmt"
	"os"
	"runtime"

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
		t.Fatal(err)
	}
	assert.Check(t, header != nil)
}

func AATestPxFileHeader(t *testing.T) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// parseFile(t, "../data/statfin_ehk_pxt_005_en.px")
	parseFile(t, "../data/010_kats_tau_101.px")

	runtime.ReadMemStats(&m2)
	fmt.Println("total:", m2.TotalAlloc-m1.TotalAlloc)
	fmt.Println("mallocs:", m2.Mallocs-m1.Mallocs)

	t.Error("test")
}

func BenchmarkPxFileHeader(b *testing.B) {
	for i := 0; i < b.N; i++ {

		r, err := os.Open("../data/010_kats_tau_101.px")
		if err != nil {
			panic(err)
		}
		defer r.Close()
		header, err := ast.PxParser.Parse("", r, participle.AllowTrailing(true))
		// repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
		if err != nil {
			fmt.Println("errored")
			b.Fatal(err)
		}
		assert.Check(b, header != nil)
		// b.Error("test")
	}
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
		t.Fatal(err)
	}
	assert.Check(t, header != nil)
}

func TestTerminate2(t *testing.T) {
	text := `A=1;
VALUES("Data")="Annual change %","Quartal change %","Value, M";
PRECISION("Data","Value, M")=1;
DATA=
"." "." 325.3 
60.078759 "." 520.8 
`
	header, err := ast.PxParser.ParseString("", text, participle.AllowTrailing(true))
	// repr.Println(header, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		t.Fatal(err)
	}
	assert.Check(t, header != nil)
}
