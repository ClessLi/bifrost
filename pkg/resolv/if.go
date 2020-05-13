package resolv

type If struct {
	BasicContext `json:"if"`
}

func (i *If) QueryAll(kw KeyWords) (parsers []Parser) {
	if i.filter(kw) {
		parsers = append(parsers, i)
	}
	return i.subQueryAll(parsers, kw)
}

func (i *If) Query(kw KeyWords) (parser Parser) {
	if i.filter(kw) {
		parser = i
	}
	return i.subQuery(kw)
}

func (i *If) BitSize(order Order, bit int) byte {
	return 0
}

func (i *If) BitLen(order Order) int {
	return 0
}

func (i *If) Size(order Order) int {
	return 0
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     TypeIf,
		Value:    value,
		Children: nil,
	}}
}
