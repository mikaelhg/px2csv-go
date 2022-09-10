package ast

import (
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

var (
	PxLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z-_\d]*`},
		{Name: `Integer`, Pattern: `\d+`},
		{Name: `Punct`, Pattern: `[][=;(),"]`},
		{Name: `EOL`, Pattern: `[\n\r]+`},
		{Name: `whitespace`, Pattern: `\s+`},
	})
	PxParser = participle.MustBuild[PxFileHeader](
		participle.Lexer(PxLexer),
		participle.Unquote("String"),
		participle.Elide("whitespace", "EOL"),
		// participle.UseLookahead(5),
	)
)

type PxKeyword struct {
	Keyword string `parser:" ( (?! 'DATA' ) @Ident )! "`
	// Keyword    string    `parser:" @Ident  "`
	Language   *string   `parser:"( '[' @Ident ']' )?"`
	Specifiers *[]string `parser:"( '(' @String ( ',' @String )* ')' )?"`
}

type PxValue struct {
	Integer *int         `parser:"   @Integer"`
	Times   *[]PxTimeVal `parser:"  | @@ (',' @@)* "`
	String  *string      `parser:"  | @Ident "`
	List    *[]string    `parser:"  | ( @String ( ',' @String )* ) "`
}

type PxTimeVal struct {
	Units string    `parser:" 'TLIST' '(' @( 'A1' | 'H1' | 'Q1' | 'M1' | 'W1' ) "`
	Times *[]string `parser:" ( ',' @String ( ',' @String )* )? ')' ( ',' @String )* "`
}

type PxRow struct {
	Keyword PxKeyword `parser:" @@ '=' "`
	Value   PxValue   `parser:" @@ ';' "`
}

type PxFileHeader struct {
	Row []PxRow `parser:"( @@ )* 'DATA' '=' "`
}
