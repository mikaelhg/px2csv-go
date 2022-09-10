package ast_test

import (
	"github.com/mikaelhg/gpcaxis/ast"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"testing"

	"github.com/alecthomas/participle/v2"
)

var (
	rowParser = participle.MustBuild[ast.PxRow](
		participle.Lexer(ast.PxLexer),
		participle.Unquote("String"),
		participle.Elide("whitespace", "EOL"),
	)
)

func TestPxRowWithLang(t *testing.T) {
	text := `SUBJECT-AREA[sv]="Besiktningar av personbilar";`

	sv := "sv"

	er := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:    "SUBJECT-AREA",
			Language:   &sv,
			Specifiers: nil,
		},
		Value: ast.PxValue{
			Integer: nil,
			String:  nil,
			List: &[]string{
				"Besiktningar av personbilar",
			},
		},
	}

	r, err := rowParser.ParseString("", text)
	// repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}

	assert.Check(t, cmp.DeepEqual(er, *r))
}

func TestPxRow(t *testing.T) {
	text := `SUBJECT-AREA="Besiktningar av personbilar";`

	er := ast.PxRow{
		Keyword: ast.PxKeyword{
			Keyword:    "SUBJECT-AREA",
			Language:   nil,
			Specifiers: nil,
		},
		Value: ast.PxValue{
			Integer: nil,
			String:  nil,
			List: &[]string{
				"Besiktningar av personbilar",
			},
		},
	}

	r, err := rowParser.ParseString("", text, participle.AllowTrailing(true))
	// repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))
	if err != nil {
		panic(err)
	}

	assert.Check(t, cmp.DeepEqual(er, *r))
}
