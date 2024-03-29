package internal

import (
	"bufio"
	"errors"
	"io"

	"github.com/samber/lo"
)

const DataValueWidth = 128 // max width of 32 bit float string in bytes

type PxParser struct {
	state      PxParserState
	row        RowAccumulator
	headers    []PxHeaderRow
	CubeWriter StatCubeWriter
}

func (p *PxParser) Header(keyword string, language string, subkeys []string) []string {
	for _, v := range p.headers {
		if v.Equals(keyword, language, subkeys) {
			return v.Values
		}
	}
	return nil
}

func (p *PxParser) valuesHeader(subkey string, _ int) []string {
	return p.Header("VALUES", "", []string{subkey})
}

func (p *PxParser) denseHeading() ([][]string, int) {
	heading := p.Header("HEADING", "", []string{})
	headingValues := lo.Map(heading, p.valuesHeader)
	headingFlattener := NewCartesianProduct(headingValues)
	headingFlattened := headingFlattener.All()
	return headingFlattened, len(headingFlattened)
}

func (p *PxParser) denseStub() ([]string, CartesianProduct, int) {
	stub := p.Header("STUB", "", []string{})
	stubValues := lo.Map(stub, p.valuesHeader)
	return stub, NewCartesianProduct(stubValues), len(stub)
}

func (p *PxParser) ParseDataDense(reader *bufio.Reader) {
	stub, stubFlattener, stubWidth := p.denseStub()
	headingFlattened, headingWidth := p.denseHeading()
	p.CubeWriter.WriteHeading(stub, headingFlattened)

	base, bufLength, currentValue := 0, 0, 0
	buffer := make([]byte, headingWidth*DataValueWidth)
	valueLengths := make([]int, headingWidth)
	currentStubs := make([]*string, stubWidth)

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
		base = DataValueWidth * currentValue
		if c == '"' {
			continue parser

		} else if c == ' ' || c == '\n' || c == '\r' || c == ';' {
			if bufLength > 0 {
				valueLengths[currentValue] = bufLength
				bufLength = 0
				currentValue += 1
			}
			if currentValue == headingWidth {
				currentValue = 0
				stubFlattener.NextP(&currentStubs)
				p.CubeWriter.WriteRow(&currentStubs, &buffer,
					&valueLengths, stubWidth, headingWidth)
			}
		} else {
			buffer[base+bufLength] = c
			bufLength += 1
		}
	}

	p.CubeWriter.WriteFooting()
}

func (p *PxParser) ParseHeader(reader *bufio.Reader) {
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
func (p *PxParser) ParseHeaderCharacter(c byte) (stop bool, err error) {
	inQuotes := p.state.Quotes%2 == 1
	inParenthesis := p.state.ParenthesisOpen > p.state.ParenthesisClose
	inKey := p.state.Semicolons == p.state.Equals
	inLanguage := inKey && p.state.SquarebracketOpen > p.state.SquarebracketClose
	inSubkey := inKey && inParenthesis

	p.state.Count += 1

	if c == '"' {
		p.state.Quotes += 1

	} else if (c == '\n' || c == '\r') && inQuotes {
		return true, errors.New("there can't be newlines inside quoted strings")

	} else if (c == '\n' || c == '\r') && !inQuotes {
		return false, nil

	} else if c == '[' && inKey && !inQuotes {
		p.state.SquarebracketOpen += 1

	} else if c == ']' && inKey && !inQuotes {
		p.state.SquarebracketClose += 1

	} else if c == '(' && inKey && !inQuotes {
		p.state.ParenthesisOpen += 1

	} else if c == '(' && !inKey && !inQuotes {
		// TLIST opening quote
		p.state.ParenthesisOpen += 1
		p.row.Value += string(c)

	} else if c == ')' && inKey && !inQuotes {
		p.state.ParenthesisClose += 1
		p.row.Subkeys = append(p.row.Subkeys, p.row.Subkey)
		p.row.Subkey = ""

	} else if c == ')' && !inKey && !inQuotes {
		// TLIST closing quote
		p.state.ParenthesisClose += 1
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
		p.state.Equals += 1

	} else if c == ';' && inKey && !inQuotes {
		return true, errors.New("found a semicolon without a matching equals sign, value terminator without keyword terminator")

	} else if c == ';' && !inKey && !inQuotes {
		if len(p.row.Value) > 0 {
			p.row.Values = append(p.row.Values, p.row.Value)
		}
		p.state.Semicolons += 1
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
