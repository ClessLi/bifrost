package nginx

import (
	"strconv"
)

type Server struct {
	BasicContext `json:"server"`
}

func (s *Server) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if s.filter(kw) {
		parsers = append(parsers, s)
	}
	if kw.IsRec {
		return s.subQueryAll(parsers, kw)
	} else {
		return
	}
}

func (s *Server) QueryByKeywords(kw Keywords) (parser Parser) {
	if s.filter(kw) {
		return s
	}
	if kw.IsRec {
		return s.subQuery(kw)
	} else {
		return
	}
}

func (s *Server) BitSize(order Order, bit int) byte {
	switch order {
	case ServerName:
		serverName := GetServerName(s)
		if serverName == nil {
			return 0
		}
		sn := []byte(StripSpace(serverName.(*Key).Value))
		// sn := []byte(serverName.(*Key).Value)

		if len(sn) <= bit {
			return 0
		}

		return sn[bit]
	default:
		return 0
	}
}

func (s *Server) BitLen(order Order) int {
	switch order {
	case ServerName:
		serverName := GetServerName(s)
		if serverName == nil {
			return 0
		}
		sn := StripSpace(serverName.(*Key).Value)
		// sn := serverName.(*Key).Value
		return len([]byte(sn))
	default:
		return 0
	}
}

func (s *Server) Size(order Order) int {
	switch order {
	case ServerPort:
		weight, err := strconv.Atoi(StripSpace(GetPorts(s)[0].(*Key).Value))
		if err != nil {
			weight = 0
		}
		return weight
	default:
		return 0
	}
}

func NewServer() *Server {
	return &Server{BasicContext{
		Name:     TypeServer,
		Value:    "",
		Children: nil,
	}}
}
