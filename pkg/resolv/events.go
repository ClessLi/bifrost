package resolv

type Events struct {
	BasicContext `json:"events"`
}

func (e *Events) Filter(kw KeyWords) (parsers []Parser) {
	if e.filter(kw) {
		parsers = append(parsers, e)
	}
	return e.subFilter(parsers, kw)
}

func NewEvents() *Events {
	return &Events{BasicContext{
		Name:     TypeEvents,
		Value:    "",
		Children: nil,
	}}
}
