package resolv

type If struct {
	BasicContext `json:"if"`
}

func (i *If) Filter(kw KeyWords) (parsers []Parser) {
	if i.filter(kw) {
		parsers = append(parsers, i)
	}
	return i.subFilter(parsers, kw)
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     TypeIf,
		Value:    value,
		Children: nil,
	}}
}
