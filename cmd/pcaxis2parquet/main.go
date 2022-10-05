package main

import (
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	f, err := os.Open("../data/010_kats_tau_101.px")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pxParser := Parser{}
	pxParser.ParseHeader(f)

}

type Parser struct {
	HeaderParserState internal.HeaderParseState
	DataParserState   internal.DataParserState
	RowAccumulator    internal.RowAccumulator
	Headers           map[internal.PxHeaderKeyword]internal.PxHeaderValue
}

func (p Parser) ParseDataDense() {
	p.HeaderParserState.Count += 1
}

func (p Parser) ParseHeader(input *os.File) {

}

func (p Parser) ParseHeaderCharacter(c int8) bool {
	inQuotes := p.HeaderParserState.Quotes%2 == 0
	inParenthesis := p.HeaderParserState.ParenthesisOpen > p.HeaderParserState.ParenthesisClose
	inKey := p.HeaderParserState.Semicolons == p.HeaderParserState.Equals
	inLanguage := inKey && p.HeaderParserState.SquarebracketOpen > p.HeaderParserState.SquarebracketClose
	inSubkey := inKey && inParenthesis

	if c == '"' {
		p.HeaderParserState.Quotes += 1
	} else if (c == '\n' || c == '\r') && inQuotes {
		println("There can't be newlines inside quoted strings.")
	} else if (c == '\n' || c == '\r') && !inQuotes {
		return false
	} else if c == '[' && inKey && !inQuotes {
		p.HeaderParserState.SquarebracketOpen += 1
	} else if c == ']' && inKey && !inQuotes {
		p.HeaderParserState.SquarebracketClose += 1
	} else if c == '(' && inKey && !inQuotes {
		p.HeaderParserState.ParenthesisOpen += 1
	} else if c == '(' && !inKey && !inQuotes {
		// TLIST opening quote
		p.HeaderParserState.ParenthesisOpen += 1
		p.RowAccumulator.Value += c
	} else if c == ')' && inKey && !inQuotes {
		p.HeaderParserState.ParenthesisClose += 1
		p.RowAccumulator.Subkeys = append(p.RowAccumulator.Subkeys, p.RowAccumulator.Subkey)
		p.RowAccumulator.Subkey = ""
	} else if c == ')' && !inKey && !inQuotes {
		// TLIST closing quote
		p.HeaderParserState.ParenthesisClose += 1
		p.RowAccumulator.Value += c
	} else if c == ',' && inSubkey && !inQuotes {
		p.RowAccumulator.Subkeys = append(p.RowAccumulator.Subkeys, p.RowAccumulator.Subkey)
		p.RowAccumulator.Subkey = ""
	} else if c == ',' && !inKey && !inQuotes && !inParenthesis {
		p.RowAccumulator.Values = append(p.RowAccumulator.Values, p.RowAccumulator.Value)
		p.RowAccumulator.Value = ""
	} else if c == '=' && !inKey && !inQuotes {
		println("Found a second equals sign without a matching semicolon. Unexpected keyword terminator.")
	} else if c == '=' && inKey && !inQuotes {
		if p.RowAccumulator.Keyword == "DATA" {
			return true
		}
		p.HeaderParserState.Equals += 1
	} else if c == ';' && inKey && !inQuotes {
		println("Found a second equals sign without a matching semicolon. Unexpected keyword terminator.")
	} else if c == ';' && !inKey && !inQuotes {
		if len(p.RowAccumulator.Value) > 0 {
			p.RowAccumulator.Values = append(p.RowAccumulator.Values, p.RowAccumulator.Value)
		}
		p.HeaderParserState.Semicolons += 1
		p.Headers[p.RowAccumulator.ToKeyword()] = p.RowAccumulator.ToValue()
		p.RowAccumulator = internal.RowAccumulator{}
		return false
	} else if inSubkey {
		p.RowAccumulator.Subkey += c
	} else if inLanguage {
		p.RowAccumulator.Language += c
	} else if inKey {
		p.RowAccumulator.Keyword += c
	} else {
		p.RowAccumulator.Value += c
	}

	return false
}
