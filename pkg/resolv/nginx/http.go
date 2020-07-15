package nginx

type Http struct {
	BasicContext `json:"http"`
}

func (h *Http) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if h.filter(kw) {
		parsers = append(parsers, h)
	}
	if kw.IsRec {
		return h.subQueryAll(parsers, kw)
	} else {
		return parsers
	}
}

func (h *Http) QueryByKeywords(kw Keywords) (parser Parser) {
	if h.filter(kw) {
		return h
	}
	if kw.IsRec {
		return h.subQuery(kw)
	} else {
		return
	}
}

func (h *Http) BitSize(_ Order, _ int) byte {
	return 0
}

func (h *Http) BitLen(_ Order) int {
	return 0
}

func (h *Http) Size(_ Order) int {
	return 0
}

func NewHttp() *Http {
	return &Http{BasicContext{
		Name:     TypeHttp,
		Value:    "",
		Children: nil,
	}}
}
