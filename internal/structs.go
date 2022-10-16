package internal

import "golang.org/x/exp/slices"

type PxHeaderRow struct {
	Keyword  string
	Language string
	Subkeys  []string
	Values   []string
}

func (r *PxHeaderRow) Equals(keyword string, language string, subkeys []string) bool {
	return r.Keyword == keyword && r.Language == language && slices.Equal(r.Subkeys, subkeys)
}

type RowAccumulator struct {
	Keyword  string
	Language string
	Subkey   string
	Subkeys  []string
	Value    string
	Values   []string
}

func (r *RowAccumulator) ToRow() PxHeaderRow {
	return PxHeaderRow{
		Keyword:  r.Keyword,
		Language: r.Language,
		Subkeys:  r.Subkeys,
		Values:   r.Values,
	}
}

type HeaderParseState struct {
	Count              int
	Quotes             int
	Semicolons         int
	Equals             int
	SquarebracketOpen  int
	SquarebracketClose int
	ParenthesisOpen    int
	ParenthesisClose   int
}
