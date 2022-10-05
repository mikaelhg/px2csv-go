package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/mikaelhg/gpcaxis/internal"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer f.Close()

	pxParser := Parser{
		Headers: make(map[internal.PxHeaderKeyword]internal.PxHeaderValue),
	}
	pxParser.ParseHeader(f)

}

type Parser struct {
	HeaderParserState internal.HeaderParseState
	DataParserState   internal.DataParserState
	RowAccumulator    internal.RowAccumulator
	Headers           map[internal.PxHeaderKeyword]internal.PxHeaderValue
}

func (p Parser) ParseDataDense() {
}

func (p *Parser) ParseHeader(input *os.File) {
	reader := bufio.NewReader(input)
	buffer := make([]byte, 4096)
out:
	for {
		n, err := reader.Read(buffer)
		if err != nil {
			panic(err)
		}
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			stop, err := p.ParseHeaderCharacter(buffer[i])
			if err != nil {
				fmt.Printf("%#v\n", p.HeaderParserState)
				panic(err)
			}
			if stop {
				// rewind the reader
				break out
			}
		}
	}
	fmt.Printf("%#v\n", p.HeaderParserState)
}

func (p *Parser) ParseHeaderCharacter(c byte) (bool, error) {
	inQuotes := p.HeaderParserState.Quotes%2 == 1
	inParenthesis := p.HeaderParserState.ParenthesisOpen > p.HeaderParserState.ParenthesisClose
	inKey := p.HeaderParserState.Semicolons == p.HeaderParserState.Equals
	inLanguage := inKey && p.HeaderParserState.SquarebracketOpen > p.HeaderParserState.SquarebracketClose
	inSubkey := inKey && inParenthesis

	p.HeaderParserState.Count += 1

	if c == '"' {
		p.HeaderParserState.Quotes += 1

	} else if (c == '\n' || c == '\r') && inQuotes {
		return true, errors.New("There can't be newlines inside quoted strings.")

	} else if (c == '\n' || c == '\r') && !inQuotes {
		return false, nil

	} else if c == '[' && inKey && !inQuotes {
		p.HeaderParserState.SquarebracketOpen += 1

	} else if c == ']' && inKey && !inQuotes {
		p.HeaderParserState.SquarebracketClose += 1

	} else if c == '(' && inKey && !inQuotes {
		p.HeaderParserState.ParenthesisOpen += 1

	} else if c == '(' && !inKey && !inQuotes {
		// TLIST opening quote
		p.HeaderParserState.ParenthesisOpen += 1
		p.RowAccumulator.Value += string(c)

	} else if c == ')' && inKey && !inQuotes {
		p.HeaderParserState.ParenthesisClose += 1
		p.RowAccumulator.Subkeys = append(p.RowAccumulator.Subkeys, p.RowAccumulator.Subkey)
		p.RowAccumulator.Subkey = ""

	} else if c == ')' && !inKey && !inQuotes {
		// TLIST closing quote
		p.HeaderParserState.ParenthesisClose += 1
		p.RowAccumulator.Value += string(c)

	} else if c == ',' && inSubkey && !inQuotes {
		p.RowAccumulator.Subkeys = append(p.RowAccumulator.Subkeys, p.RowAccumulator.Subkey)
		p.RowAccumulator.Subkey = ""

	} else if c == ',' && !inKey && !inQuotes && !inParenthesis {
		p.RowAccumulator.Values = append(p.RowAccumulator.Values, p.RowAccumulator.Value)
		p.RowAccumulator.Value = ""

	} else if c == '=' && !inKey && !inQuotes {
		return true, errors.New("Found a second equals sign without a matching semicolon. Unexpected keyword terminator.")

	} else if c == '=' && inKey && !inQuotes {
		if p.RowAccumulator.Keyword == "DATA" {
			return true, nil
		}
		p.HeaderParserState.Equals += 1

	} else if c == ';' && inKey && !inQuotes {
		return true, errors.New("Found a second equals sign without a matching semicolon. Unexpected keyword terminator.")

	} else if c == ';' && !inKey && !inQuotes {
		if len(p.RowAccumulator.Value) > 0 {
			p.RowAccumulator.Values = append(p.RowAccumulator.Values, p.RowAccumulator.Value)
		}
		p.HeaderParserState.Semicolons += 1
		p.Headers[p.RowAccumulator.ToKeyword()] = p.RowAccumulator.ToValue()
		p.RowAccumulator = internal.RowAccumulator{}
		return false, nil

	} else if inSubkey {
		p.RowAccumulator.Subkey += string(c)

	} else if inLanguage {
		p.RowAccumulator.Language += string(c)

	} else if inKey {
		p.RowAccumulator.Keyword += string(c)

	} else {
		p.RowAccumulator.Value += string(c)
	}

	return false, nil
}
