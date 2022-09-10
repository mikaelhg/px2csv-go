package ast_test

import (
	"github.com/mikaelhg/gpcaxis/ast"

	"testing"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

var (
	l = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z-_\d]*`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Integer`, Pattern: `\d+`},
		{Name: "Punct", Pattern: `[][=;(),"]`},
		{Name: "whitespace", Pattern: `\s+`},
		{Name: "EOL", Pattern: `[\n\r]+`},
	})
	rp = participle.MustBuild[ast.PxRow](
		participle.Lexer(l),
		participle.Unquote("String"),
	)
	vp = participle.MustBuild[ast.PxValue](
		participle.Lexer(l),
		participle.Unquote("String"),
	)
)

func TestPxRow(t *testing.T) {
	text := `SUBJECT-AREA[sv]="Besiktningar av personbilar";`
	t.Error("Foo")

	r, err := rp.ParseString("", text)
	repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))

	if err != nil {
		panic(err)
	}
}

func TestPxValueStrings(t *testing.T) {
	text := `"AAA","BBB","CCC"`
	t.Error("Foo")

	r, err := vp.ParseString("", text)
	repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))

	if err != nil {
		panic(err)
	}
}

func TestPxValueInt(t *testing.T) {
	text := `123`
	t.Error("Foo")

	r, err := vp.ParseString("", text)
	repr.Println(r, repr.Indent("  "), repr.OmitEmpty(false))

	if err != nil {
		panic(err)
	}
}
