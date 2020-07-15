package nginx

type Map struct {
	BasicContext `json:"map"`
}

func (m *Map) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if m.filter(kw) {
		parsers = append(parsers, m)
	}
	if kw.IsRec {
		return m.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (m *Map) QueryByKeywords(kw Keywords) (parser Parser) {
	if m.filter(kw) {
		return m
	}
	if kw.IsRec {
		return m.subQuery(kw)
	} else {
		return
	}
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
