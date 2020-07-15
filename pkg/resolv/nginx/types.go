package nginx

type Types struct {
	BasicContext `json:"types"`
}

func (t *Types) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if t.filter(kw) {
		parsers = append(parsers, t)
	}
	if kw.IsRec {
		return t.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (t *Types) QueryByKeywords(kw Keywords) (parser Parser) {
	if t.filter(kw) {
		return t
	}
	if kw.IsRec {
		return t.subQuery(kw)
	} else {
		return
	}
}

func (t *Types) BitSize(_ Order, _ int) byte {
	return 0
}

func (t *Types) BitLen(_ Order) int {
	return 0
}

func (t *Types) Size(_ Order) int {
	return 0
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     TypeTypes,
		Value:    "",
		Children: nil,
	}}
}
