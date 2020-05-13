package resolv

type Types struct {
	BasicContext `json:"types"`
}

func (t *Types) QueryAll(kw KeyWords) (parsers []Parser) {
	if t.filter(kw) {
		parsers = append(parsers, t)
	}
	return t.subQueryAll(parsers, kw)
}

func (t *Types) Query(kw KeyWords) (parser Parser) {
	if t.filter(kw) {
		parser = t
	}
	return t.subQuery(kw)
}

func (t *Types) BitSize(order Order, bit int) byte {
	return 0
}

func (t *Types) BitLen(order Order) int {
	return 0
}

func (t *Types) Size(order Order) int {
	return 0
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     TypeTypes,
		Value:    "",
		Children: nil,
	}}
}
