package resolv

type Types struct {
	BasicContext `json:"types"`
}

func (t *Types) Filter(kw KeyWords) (parsers []Parser) {
	if t.filter(kw) {
		parsers = append(parsers, t)
	}
	return t.subFilter(parsers, kw)
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     TypeTypes,
		Value:    "",
		Children: nil,
	}}
}
