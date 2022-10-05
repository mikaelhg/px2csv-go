package internal

type PxHeaderKeyword struct {
	Keyword  string
	Language string
	Subkeys  []string
}

type PxHeaderValue struct {
	Values []string
}

type PxHeaderRow struct {
	Keyword PxHeaderKeyword
	Value   PxHeaderValue
}

type RowAccumulator struct {
	Keyword  string
	Language string
	Subkey   string
	Subkeys  []string
	Value    string
	Values   []string
}

func (r *RowAccumulator) ToKeyword() PxHeaderKeyword {
	return PxHeaderKeyword{
		Keyword:  r.Keyword,
		Language: r.Language,
		Subkeys:  r.Subkeys,
	}
}

func (r *RowAccumulator) ToValue() PxHeaderValue {
	return PxHeaderValue{Values: r.Values}
}

func (r *RowAccumulator) ToRow() PxHeaderRow {
	return PxHeaderRow{Keyword: r.ToKeyword(), Value: r.ToValue()}
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

type DataParserState struct {
	Count int
}
