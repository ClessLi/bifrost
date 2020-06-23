package nginx

type Map struct {
	BasicContext `json:"map"`
}

func (m *Map) QueryAll(kw KeyWords) (parsers []Parser) {
	if m.filter(kw) {
		parsers = append(parsers, m)
	}
	return m.subQueryAll(parsers, kw)
}

func (m *Map) Query(kw KeyWords) (parser Parser) {
	if m.filter(kw) {
		parser = m
	}
	return m.subQuery(kw)
}

func (m *Map) BitSize(_ Order, _ int) byte {
	return 0
}

func (m *Map) BitLen(_ Order) int {
	return 0
}

func (m *Map) Size(_ Order) int {
	return 0
}

func NewMap(value string) *Map {
	return &Map{BasicContext{
		Name:     TypeMap,
		Value:    value,
		Children: nil,
	}}
}
