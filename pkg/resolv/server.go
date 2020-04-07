package resolv

import (
	"encoding/json"
)

type Server struct {
	BasicContext
}

func (s *Server) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Server []Parser `json:"server"`
	}{s.Children})
}

func (s *Server) UnmarshalJSON(b []byte) error {
	server := struct {
		Server []Parser `json:"server"`
	}{}
	err := json.Unmarshal(b, &server)
	if err != nil {
		return nil
	}

	s.Name = "server"
	s.Children = server.Server
	return nil
}

func NewServer() *Server {
	return &Server{BasicContext{
		Name:     "server",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
