package resolv

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) Filter(kw KeyWords) (parsers []Parser) {
	if g.filter(kw) {
		parsers = append(parsers, g)
	}
	return g.subFilter(parsers, kw)
}

func NewGeo(value string) *Geo {
	return &Geo{BasicContext{
		Name:     TypeGeo,
		Value:    value,
		Children: nil,
	}}
}
