package resolv

type Stream struct {
	BasicContext `json:"stream"`
}

func (s *Stream) Filter(kw KeyWords) (parsers []Parser) {
	if s.filter(kw) {
		parsers = append(parsers, s)
	}
	return s.subFilter(parsers, kw)
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     TypeStream,
		Value:    "",
		Children: nil,
	}}
}
