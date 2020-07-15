package nginx

type Location struct {
	BasicContext `json:"location"`
}

func (l *Location) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	if kw.IsRec {
		return l.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (l *Location) QueryByKeywords(kw Keywords) (parser Parser) {
	if l.filter(kw) {
		return l
	}
	if kw.IsRec {
		return l.subQuery(kw)
	} else {
		return
	}
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
