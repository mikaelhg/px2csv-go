package ast_test

import (
	"testing"

	"github.com/mikaelhg/gpcaxis/ast"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func parseRow(t *testing.T, expected ast.PxRow, text string) {
	r, err := rowParser.ParseString("", text)
	if err != nil {
		panic(err)
	}
	assert.Check(t, cmp.DeepEqual(expected, *r))
}

func TestPxRowTimeValSimple(t *testing.T) {
	text := `TIMEVAL[sv]("Besiktningsår")=TLIST(A1),"2017","2018","2019","2020","2021";`
	sv := "sv"
	er := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:    "TIMEVAL",
			Language:   &sv,
			Specifiers: &[]string{"Besiktningsår"},
		},
		Value: ast.PxValue{
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
		},
	}
	parseRow(t, er, text)
}

func TestPxRowTimeValMultipleRange(t *testing.T) {
	text := `TIMEVAL("aika")=TLIST(A1,"1994"-"1996"),TLIST(M1,"199609"-"199612");`
	expected := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:    "TIMEVAL",
			Specifiers: &[]string{"aika"},
		},
		Value: ast.PxValue{
			Times: &[]ast.PxTimeVal{
				{
					Units: "A1",
					Range: &[]string{"1994", "1996"},
				},
				{
					Units: "M1",
					Range: &[]string{"199609", "199612"},
				},
			},
		},
	}
	parseRow(t, expected, text)
}

func TestPxRowTimeValMultipleList(t *testing.T) {
	text := `TIMEVAL("aika")=TLIST(A1),"1994","1995","1996",TLIST(M1),"199609","199610","199611","199612";`
	expected := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:    "TIMEVAL",
			Specifiers: &[]string{"aika"},
		},
		Value: ast.PxValue{
			Integer: nil,
			Times: &[]ast.PxTimeVal{
				{
					Units: "A1",
					Times: &[]string{"1994", "1995", "1996"},
				},
				{
					Units: "M1",
					Times: &[]string{"199609", "199610", "199611", "199612"},
				},
			},
		},
	}
	parseRow(t, expected, text)
}
