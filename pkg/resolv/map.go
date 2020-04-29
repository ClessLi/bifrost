package resolv

type Map struct {
	BasicContext `json:"map"`
}

func (m *Map) Filter(kw KeyWords) (parsers []Parser) {
	if m.filter(kw) {
		parsers = append(parsers, m)
	}
	return m.subFilter(parsers, kw)
}

func NewMap(value string) *Map {
	return &Map{BasicContext{
		Name:     TypeMap,
		Value:    value,
		Children: nil,
	}}
}
