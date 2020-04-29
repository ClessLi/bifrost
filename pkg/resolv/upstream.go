package resolv

type Upstream struct {
	BasicContext `json:"upstream"`
}

func (u *Upstream) Filter(kw KeyWords) (parsers []Parser) {
	if u.filter(kw) {
		parsers = append(parsers, u)
	}
	return u.subFilter(parsers, kw)
}

func NewUpstream(value string) *Upstream {
	return &Upstream{BasicContext{
		Name:     TypeUpstream,
		Value:    value,
		Children: nil,
	}}
}
