package internal

import (
	"bufio"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/exp/slices"
)

type Parser struct {
	hps     HeaderParseState
	row     RowAccumulator
	headers []PxHeaderRow
}

func NewParser() Parser {
	return Parser{}
}

func (p *Parser) Header(keyword string, language string, subkeys []string) []string {
	for _, v := range p.headers {
		if v.Keyword.Keyword == keyword &&
			v.Keyword.Language == language &&
			slices.Equal(v.Keyword.Subkeys, subkeys) {
			return v.Value.Values
		}
	}
	return nil
}

func (p *Parser) ParseDataDense(reader *bufio.Reader) {
	fn := func(x string) []string {
		return p.Header("VALUES", "", []string{x})
	}

	stub := p.Header("STUB", "", []string{})
	stubValues := MapXtoY(stub, fn)

	heading := p.Header("HEADING", "", []string{})
	headingValues := MapXtoY(heading, fn)
	headingCsv := MapXtoY(headingValues, func(x []string) string { return strings.Join(x, " ") })

	fmt.Printf("stub: %#v\n", stub)
	fmt.Printf("heading: %#v\n", heading)
	fmt.Printf("stubValues: %#v\n", stubValues)
	fmt.Printf("headingValues: %#v\n", headingValues)
	fmt.Printf("headingCsv: %#v\n", headingCsv)

}

func (p *Parser) ParseDataDenseCharacter(c byte) (bool, error) {
	return false, nil
}

func (p *Parser) ParseHeader(reader *bufio.Reader) {
	for {
		c, err := reader.ReadByte()
		if err != nil {
			panic(err)
		}
		stop, err := p.ParseHeaderCharacter(c)
		if err != nil {
			fmt.Printf("%#v\n", p.hps)
			// fmt.Printf("%#v\n", p.headers)
			panic(err)
		}
		if stop {
			fmt.Printf("%#v\n", p.hps)
			// fmt.Printf("%#v\n", p.headers)
			return
		}
	}
}

func (p *Parser) ParseHeaderCharacter(c byte) (bool, error) {
	inQuotes := p.hps.Quotes%2 == 1
	inParenthesis := p.hps.ParenthesisOpen > p.hps.ParenthesisClose
	inKey := p.hps.Semicolons == p.hps.Equals
	inLanguage := inKey && p.hps.SquarebracketOpen > p.hps.SquarebracketClose
	inSubkey := inKey && inParenthesis

	p.hps.Count += 1

	if c == '"' {
		p.hps.Quotes += 1

	} else if (c == '\n' || c == '\r') && inQuotes {
		return true, errors.New("there can't be newlines inside quoted strings")

	} else if (c == '\n' || c == '\r') && !inQuotes {
		return false, nil

	} else if c == '[' && inKey && !inQuotes {
		p.hps.SquarebracketOpen += 1

	} else if c == ']' && inKey && !inQuotes {
		p.hps.SquarebracketClose += 1

	} else if c == '(' && inKey && !inQuotes {
		p.hps.ParenthesisOpen += 1

	} else if c == '(' && !inKey && !inQuotes {
		// TLIST opening quote
		p.hps.ParenthesisOpen += 1
		p.row.Value += string(c)

	} else if c == ')' && inKey && !inQuotes {
		p.hps.ParenthesisClose += 1
		p.row.Subkeys = append(p.row.Subkeys, p.row.Subkey)
		p.row.Subkey = ""

	} else if c == ')' && !inKey && !inQuotes {
		// TLIST closing quote
		p.hps.ParenthesisClose += 1
		p.row.Value += string(c)

	} else if c == ',' && inSubkey && !inQuotes {
		p.row.Subkeys = append(p.row.Subkeys, p.row.Subkey)
		p.row.Subkey = ""

	} else if c == ',' && !inKey && !inQuotes && !inParenthesis {
		p.row.Values = append(p.row.Values, p.row.Value)
		p.row.Value = ""

	} else if c == '=' && !inKey && !inQuotes {
		return true, errors.New("found a second equals sign without a matching semicolon, unexpected keyword terminator")

	} else if c == '=' && inKey && !inQuotes {
		if p.row.Keyword == "DATA" {
			return true, nil
		}
		p.hps.Equals += 1

	} else if c == ';' && inKey && !inQuotes {
		return true, errors.New("found a second equals sign without a matching semicolon. unexpected keyword terminator")

	} else if c == ';' && !inKey && !inQuotes {
		if len(p.row.Value) > 0 {
			p.row.Values = append(p.row.Values, p.row.Value)
		}
		p.hps.Semicolons += 1
		p.headers = append(p.headers, p.row.ToRow())
		p.row = RowAccumulator{}
		return false, nil

	} else if inSubkey {
		p.row.Subkey += string(c)

	} else if inLanguage {
		p.row.Language += string(c)

	} else if inKey {
		p.row.Keyword += string(c)

	} else {
		p.row.Value += string(c)
	}

	return false, nil
}
