package nginx

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if g.filter(kw) {
		parsers = append(parsers, g)
	}
	if kw.IsRec {
		return g.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (g *Geo) QueryByKeywords(kw Keywords) (parser Parser) {
	if g.filter(kw) {
		return g
	}
	if kw.IsRec {
		return g.subQuery(kw)
	} else {
		return
	}
}

func (g *Geo) BitSize(_ Order, _ int) byte {
	return 0
}

func (g *Geo) BitLen(_ Order) int {
	return 0
}

func (g *Geo) Size(_ Order) int {
	return 0
}

func NewGeo(value string) *Geo {
	return &Geo{BasicContext{
		Name:     TypeGeo,
		Value:    value,
		Children: nil,
	}}
}
