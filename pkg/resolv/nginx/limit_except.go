package nginx

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func (l *LimitExcept) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	if kw.IsRec {
		return l.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (l *LimitExcept) QueryByKeywords(kw Keywords) (parser Parser) {
	if l.filter(kw) {
		return l
	}
	if kw.IsRec {
		return l.subQuery(kw)
	} else {
		return
	}
}

func (l *LimitExcept) BitSize(_ Order, _ int) byte {
	return 0
}

func (l *LimitExcept) BitLen(_ Order) int {
	return 0
}

func (l *LimitExcept) Size(_ Order) int {
	return 0
}

func NewLimitExcept(value string) *LimitExcept {
	return &LimitExcept{BasicContext{
		Name:     TypeLimitExcept,
		Value:    value,
		Children: nil,
	}}
}
