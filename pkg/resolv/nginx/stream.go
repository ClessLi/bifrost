package nginx

type Stream struct {
	BasicContext `json:"stream"`
}

func (s *Stream) QueryAll(kw KeyWords) (parsers []Parser) {
	if s.filter(kw) {
		parsers = append(parsers, s)
	}
	return s.subQueryAll(parsers, kw)
}

func (s *Stream) Query(kw KeyWords) (parser Parser) {
	if s.filter(kw) {
		parser = s
	}
	return s.subQuery(kw)
}

func (s *Stream) BitSize(_ Order, _ int) byte {
	return 0
}

func (s *Stream) BitLen(_ Order) int {
	return 0
}

func (s *Stream) Size(_ Order) int {
	return 0
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     TypeStream,
		Value:    "",
		Children: nil,
	}}
}