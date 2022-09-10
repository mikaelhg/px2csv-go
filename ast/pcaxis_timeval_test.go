package ast_test

import (
	"testing"

	"github.com/alecthomas/repr"
	"github.com/mikaelhg/gpcaxis/ast"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestPxRowTimeValSimple(t *testing.T) {
	text := `TIMEVAL[sv]("Besiktningsår")=TLIST(A1),"2017","2018","2019","2020","2021";`

	sv := "sv"

	er := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:  "TIMEVAL",
			Language: &sv,
			Specifiers: &[]string{
				"Besiktningsår",
			},
		},
		Value: ast.PxValue{
			Integer: nil,
			Times: &[]ast.PxTimeVal{
				{
					Units: "A1",
					Times: &[]string{
						"2017",
						"2018",
						"2019",
						"2020",
						"2021",
					},
				},
			},
			String: nil,
			List:   nil,
		},
	}

	r, err := ast.PxParser.ParseString("", text)
	repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}

	assert.Check(t, cmp.DeepEqual(er, *r))
}
