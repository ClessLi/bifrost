package resolv

type Server struct {
	BasicContext `json:"server"`
}

func (s *Server) Filter(kw KeyWords) (parsers []Parser) {
	if s.filter(kw) {
		parsers = append(parsers, s)
	}
	return s.subFilter(parsers, kw)
}

func NewServer() *Server {
	return &Server{BasicContext{
		Name:     TypeServer,
		Value:    "",
		Children: nil,
	}}
}
