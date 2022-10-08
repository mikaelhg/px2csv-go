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

func (p *Parser) ParseDataDense(reader *bufio.Reader, writer *bufio.Writer) {
	valuesHeader := func(x string) []string {
		return p.Header("VALUES", "", []string{x})
	}
	joinStringSlice := func(x []string) string {
		return strings.Join(x, " ")
	}

	stub := p.Header("STUB", "", []string{})
	stubValues := MapXtoY(stub, valuesHeader)
	stubFlattener := NewCartesianProduct(stubValues)
	stubWidth := len(stub)

	heading := p.Header("HEADING", "", []string{})
	headingValues := MapXtoY(heading, valuesHeader)
	headingFlattener := NewCartesianProduct(headingValues)
	headingFlattened := headingFlattener.All()
	headingWidth := len(headingFlattened)
	headingCsv := MapXtoY(headingFlattened, joinStringSlice)

	writer.WriteString("\"")
	writer.WriteString(strings.Join(stub, "\";\""))
	writer.WriteString("\";\"")
	writer.WriteString(strings.Join(headingCsv, "\";\""))
	writer.WriteString("\"\n")

	quotes := 0
	bufLength := 0
	currentValue := 0
	var buf bytes.Buffer
	values := make([]string, headingWidth)
	buf2 := make([]byte, headingWidth*16)
	values2 := make([][]byte, headingWidth)
	valueLengths := make([]int, headingWidth)
	for {
		c, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			} else {
				panic(err)
			}
		}
		if c == '"' {
			quotes += 1

		} else if c == ' ' || c == '\n' || c == '\r' || c == ';' {
			if bufLength > 0 {
				values2[currentValue] = buf2[0:bufLength]
				valueLengths[currentValue] = bufLength
				values[currentValue] = buf.String()
				buf.Reset()
				bufLength = 0
				currentValue += 1
			}
			if currentValue == headingWidth {
				stubs, _ := stubFlattener.Next()
				writer.WriteByte('"')
				for i, s := range stubs {
					writer.WriteString(s)
					if i < stubWidth-1 {
						writer.WriteByte('"')
						writer.WriteByte(';')
						writer.WriteByte('"')
					}
				}
				writer.WriteByte('"')
				writer.WriteByte(';')
				for i, s := range values2 {
					// writer.WriteString(s)
					writer.Write(s[0:valueLengths[i]])
					if i < headingWidth-1 {
						writer.WriteByte(';')
					}
				}
				writer.WriteByte('\n')
				currentValue = 0
			}
		} else {
			buf2[bufLength+(16*currentValue)] = c
			buf.WriteByte(c)
			bufLength += 1
		}
	}

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
