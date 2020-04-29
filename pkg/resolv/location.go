package resolv

type Location struct {
	BasicContext `json:"location"`
}

func (l *Location) Filter(kw KeyWords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	return l.subFilter(parsers, kw)
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     TypeLocation,
		Value:    value,
		Children: nil,
	}}
}
