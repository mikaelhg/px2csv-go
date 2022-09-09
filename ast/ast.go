/*
Package ast contains the parser for PC-AXIS files.
*/
package ast

import (
	"os"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	iniLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z_\d]*`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Float`, Pattern: `\d+(?:\.\d+)?`},
		{Name: `Punct`, Pattern: `[][=]`},
		{Name: "comment", Pattern: `[#;][^\n]*`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	Parser = participle.MustBuild[INI](
		participle.Lexer(iniLexer),
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}),
	)
)

type INI struct {
	Properties []*Property `parser:"@@*"`
	Sections   []*Section  `parser:"@@*"`
}

type Section struct {
	Identifier string      `parser:"'[' @Ident ']'"`
	Properties []*Property `parser:"@@*"`
}

type Property struct {
	Key   string `parser:"@Ident '='"`
	Value Value  `parser:"@@"`
}

type Value interface{ value() }

type String struct {
	String string `parser:"@String"`
}

func (String) value() {}

type Number struct {
	Number float64 `parser:"@Float"`
}

func (Number) value() {}

func main() {
	ini, err := Parser.Parse("", os.Stdin)
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
