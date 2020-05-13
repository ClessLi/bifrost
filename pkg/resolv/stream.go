package resolv

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

func (s *Stream) BitSize(order Order, bit int) byte {
	return 0
}

func (s *Stream) BitLen(order Order) int {
	return 0
}

func (s *Stream) Size(order Order) int {
	return 0
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     TypeStream,
		Value:    "",
		Children: nil,
	}}
}
