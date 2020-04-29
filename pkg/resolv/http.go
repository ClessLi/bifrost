package resolv

type Http struct {
	BasicContext `json:"http"`
}

func (h *Http) Filter(kw KeyWords) (parsers []Parser) {
	if h.filter(kw) {
		parsers = append(parsers, h)
	}
	return h.subFilter(parsers, kw)
}

func NewHttp() *Http {
	return &Http{BasicContext{
		Name:     TypeHttp,
		Value:    "",
		Children: nil,
	}}
}
