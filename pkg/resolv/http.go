package resolv

import (
	"encoding/json"
)

type Http struct {
	BasicContext `json:"http"`
}

func (h *Http) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Http []Parser `json:"http"`
	}{Http: h.Children})
}

func (h *Http) UnmarshalJSON(b []byte) error {
	http := struct {
		Http []Parser `json:"http"`
	}{}
	err := json.Unmarshal(b, &http)
	if err != nil {
		return err
	}

	h.Name = "http"
	h.Children = http.Http
	return nil
}

func NewHttp() *Http {
	return &Http{BasicContext{
		Name:     "http",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
