package resolv

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func (l *LimitExcept) QueryAll(kw KeyWords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	return l.subQueryAll(parsers, kw)
}

func (l *LimitExcept) Query(kw KeyWords) (parser Parser) {
	if l.filter(kw) {
		parser = l
	}
	return l.subQuery(kw)
}

func (l *LimitExcept) BitSize(order Order, bit int) byte {
	return 0
}

func (l *LimitExcept) BitLen(order Order) int {
	return 0
}

func (l *LimitExcept) Size(order Order) int {
	return 0
}

func NewLimitExcept(value string) *LimitExcept {
	return &LimitExcept{BasicContext{
		Name:     TypeLimitExcept,
		Value:    value,
		Children: nil,
	}}
}
