package resolv

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

func (l *Location) BitSize(order Order, bit int) byte {
	return 0
}

func (l *Location) BitLen(order Order) int {
	return 0
}

func (l *Location) Size(order Order) int {
	return 0
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     TypeLocation,
		Value:    value,
		Children: nil,
	}}
}
