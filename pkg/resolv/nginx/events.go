package nginx

type Events struct {
	BasicContext `json:"events"`
}

func (e *Events) QueryAll(kw KeyWords) (parsers []Parser) {
	if e.filter(kw) {
		parsers = append(parsers, e)
	}
	return e.subQueryAll(parsers, kw)
}

func (e *Events) Query(kw KeyWords) (parser Parser) {
	if e.filter(kw) {
		parser = e
	}
	return e.subQuery(kw)
}

func (e *Events) BitSize(_ Order, _ int) byte {
	return 0
}

func (e *Events) BitLen(_ Order) int {
	return 0
}

func (e *Events) Size(_ Order) int {
	return 0
}

func NewEvents() *Events {
	return &Events{BasicContext{
		Name:     TypeEvents,
		Value:    "",
		Children: nil,
	}}
}
