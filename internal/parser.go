package internal

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

const DataValueWidth = 128 // max width of 32 bit float string in bytes

type Parser struct {
	hps        HeaderParseState
	row        RowAccumulator
	headers    []PxHeaderRow
	CubeWriter StatCubeWriter
}

func (p *Parser) Header(keyword string, language string, subkeys []string) []string {
	for _, v := range p.headers {
		if v.Equals(keyword, language, subkeys) {
			return v.Values
		}
	}
	return nil
}

func joinStringSlice(ss []string) string {
	return strings.Join(ss, " ")
}

func (p *Parser) valuesHeader(subkey string) []string {
	return p.Header("VALUES", "", []string{subkey})
}

func (p *Parser) ParseDataDense(reader *bufio.Reader) {
	stub := p.Header("STUB", "", []string{})
	stubValues := MapXtoY(stub, p.valuesHeader)
	stubFlattener := NewCartesianProduct(stubValues)
	stubWidth := len(stub)

	heading := p.Header("HEADING", "", []string{})
	headingValues := MapXtoY(heading, p.valuesHeader)
	headingFlattener := NewCartesianProduct(headingValues)
	headingFlattened := headingFlattener.All()
	headingWidth := len(headingFlattened)
	headingCsv := MapXtoY(headingFlattened, joinStringSlice)

	quotes := 0
	bufLength := 0
	currentValue := 0
	buf := make([]byte, headingWidth*DataValueWidth)
	values := make([][]byte, headingWidth)
	valueLengths := make([]int, headingWidth)
	theseStubs := make([]*string, stubWidth)

	p.CubeWriter.WriteHeading(stub, headingCsv)

	// This is the most performance-critical part of the whole program,
	// and we'll want to avoid any heap allocations inside the parser loop.
	// That's why we indulge with all this funky pointer arithmetic and
	// pre-allocation.
parser:
	for {
		c, err := reader.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break parser
			} else {
				panic(err)
			}
		}
		if c == '"' {
			quotes += 1

		} else if c == ' ' || c == '\n' || c == '\r' || c == ';' {
			if bufLength > 0 {
				base := DataValueWidth * currentValue
				values[currentValue] = buf[base : base+bufLength]
				valueLengths[currentValue] = bufLength
				bufLength = 0
				currentValue += 1
			}
			if currentValue == headingWidth {
				currentValue = 0
				stubFlattener.NextP(&theseStubs)
				p.CubeWriter.WriteRow(&theseStubs, &values,
					&valueLengths, stubWidth, headingWidth)
			}
		} else {
			buf[bufLength+(DataValueWidth*currentValue)] = c
			bufLength += 1
		}
	}

	p.CubeWriter.WriteFooting()
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

// TIMEVAL and HIERARCHY not yet supported beyond passing them through.
func (p *Parser) ParseHeaderCharacter(c byte) (stop bool, err error) {
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
		return true, errors.New("found a semicolon without a matching equals sign, value terminator without keyword terminator")

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
