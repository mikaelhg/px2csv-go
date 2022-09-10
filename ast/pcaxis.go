package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	PxLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z-_\d]*`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Integer`, Pattern: `\d+`},
		{Name: `Punct`, Pattern: `[][=;(),"]`},
		{Name: `whitespace`, Pattern: `\s+`},
		{Name: `EOL`, Pattern: `[\n\r]+`},
	})
	PxParser = participle.MustBuild[PxFileHeader](
		participle.Lexer(PxLexer),
		participle.Unquote("String"),
	)
)

type PxKeyword struct {
	Keyword    string    `parser:" @Ident "`
	Language   *string   `parser:"( '[' @Ident ']' )?"`
	Specifiers *[]string `parser:"( '(' @String ( ',' @String )* ')' )?"`
}

type PxValue struct {
	Integer *int      `parser:"   @Integer"`
	String  *string   `parser:"  | @Ident "`
	List    *[]string `parser:"  | ( @String ( ',' @String )* ) "`
}

type PxRow struct {
	Keyword PxKeyword `parser:" @@ '=' "`
	Value   PxValue   `parser:" @@ ';' "`
}

type PxFileHeader struct {
	Row PxRow `parser:"( @@ EOL )* 'DATA=' "`
}
