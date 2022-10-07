package internal

import (
	"bufio"
	"bytes"
	"errors"
	"io"
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
	stubFlattener := NewCartesianProduct(stubValues)

	heading := p.Header("HEADING", "", []string{})
	headingValues := MapXtoY(heading, fn)
	headingFlattener := NewCartesianProduct(headingValues)
	headingFlattened := headingFlattener.All()
	headingWidth := len(headingFlattened)
	headingCsv := MapXtoY(headingFlattened, func(x []string) string {
		return strings.Join(x, " ")
	})

	print("\"")
	print(strings.Join(stub, "\";\""))
	print("\";\"")
	print(strings.Join(headingCsv, "\";\""))
	println("\"")

	var buf bytes.Buffer
	values := make([]string, 0)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else {
				panic(err)
			}
		}
		if c == ' ' || c == '\n' || c == '\r' {
			if buf.Len() > 0 {
				values = append(values, buf.String())
				buf.Reset()
			}
			if len(values) == headingWidth {
				stubs, _ := stubFlattener.Next()
				print("\"")
				print(strings.Join(stubs, "\";\""))
				print("\";")
				print(strings.Join(values, ";"))
				println()
				values = make([]string, 0)
			}
		} else {
			buf.WriteByte(c)
		}
	}

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
			panic(err)
		}
		if stop {
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
