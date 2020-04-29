package resolv

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func (l *LimitExcept) Filter(kw KeyWords) (parsers []Parser) {
	if l.filter(kw) {
		parsers = append(parsers, l)
	}
	return l.subFilter(parsers, kw)
}

func NewLimitExcept(value string) *LimitExcept {
	return &LimitExcept{BasicContext{
		Name:     TypeLimitExcept,
		Value:    value,
		Children: nil,
	}}
}
