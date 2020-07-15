package nginx

type If struct {
	BasicContext `json:"if"`
}

func (i *If) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if i.filter(kw) {
		parsers = append(parsers, i)
	}
	if kw.IsRec {
		return i.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (i *If) QueryByKeywords(kw Keywords) (parser Parser) {
	if i.filter(kw) {
		return i
	}
	if kw.IsRec {
		return i.subQuery(kw)
	} else {
		return
	}
}

func (i *If) BitSize(_ Order, _ int) byte {
	return 0
}

func (i *If) BitLen(_ Order) int {
	return 0
}

func (i *If) Size(_ Order) int {
	return 0
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     TypeIf,
		Value:    value,
		Children: nil,
	}}
}
