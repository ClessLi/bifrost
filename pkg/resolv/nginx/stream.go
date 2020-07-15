package nginx

type Stream struct {
	BasicContext `json:"stream"`
}

func (s *Stream) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if s.filter(kw) {
		parsers = append(parsers, s)
	}
	if kw.IsRec {
		return s.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (s *Stream) QueryByKeywords(kw Keywords) (parser Parser) {
	if s.filter(kw) {
		return s
	}
	if kw.IsRec {
		return s.subQuery(kw)
	} else {
		return
	}
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
