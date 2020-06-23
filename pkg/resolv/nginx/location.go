package nginx

type Location struct {
	BasicContext `json:"location"`
}

func (l *Location) QueryAll(kw KeyWords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	return l.subQueryAll(parsers, kw)
}

func (l *Location) Query(kw KeyWords) (parser Parser) {
	if l.filter(kw) {
		parser = l
	}
	return l.subQuery(kw)
}

func (l *Location) BitSize(_ Order, _ int) byte {
	return 0
}

func (l *Location) BitLen(_ Order) int {
	return 0
}

func (l *Location) Size(_ Order) int {
	return 0
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     TypeLocation,
		Value:    value,
		Children: nil,
	}}
}
