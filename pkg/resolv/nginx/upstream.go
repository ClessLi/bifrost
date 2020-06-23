package nginx

type Upstream struct {
	BasicContext `json:"upstream"`
}

func (u *Upstream) QueryAll(kw KeyWords) (parsers []Parser) {
	if u.filter(kw) {
		parsers = append(parsers, u)
	}
	return u.subQueryAll(parsers, kw)
}

func (u *Upstream) Query(kw KeyWords) (parser Parser) {
	if u.filter(kw) {
		parser = u
	}
	return u.subQuery(kw)
}

func (u *Upstream) BitSize(_ Order, _ int) byte {
	return 0
}

func (u *Upstream) BitLen(_ Order) int {
	return 0
}

func (u *Upstream) Size(_ Order) int {
	return 0
}

func NewUpstream(value string) *Upstream {
	return &Upstream{BasicContext{
		Name:     TypeUpstream,
		Value:    value,
		Children: nil,
	}}
}
