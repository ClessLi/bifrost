package resolv

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) QueryAll(kw KeyWords) (parsers []Parser) {
	if g.filter(kw) {
		parsers = append(parsers, g)
	}
	return g.subQueryAll(parsers, kw)
}

func (g *Geo) Query(kw KeyWords) (parser Parser) {
	if g.filter(kw) {
		parser = g
	}
	return g.subQuery(kw)
}

func (g *Geo) BitSize(order Order, bit int) byte {
	return 0
}

func (g *Geo) BitLen(order Order) int {
	return 0
}

func (g *Geo) Size(order Order) int {
	return 0
}

func NewGeo(value string) *Geo {
	return &Geo{BasicContext{
		Name:     TypeGeo,
		Value:    value,
		Children: nil,
	}}
}
