package resolv

import "encoding/json"

type Geo struct {
	BasicContext `json:"geo"`
}

func (g *Geo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Geo []Parser `json:"geo"`
	}{Geo: g.Children})
}

func (g *Geo) UnmarshalJSON(b []byte) error {
	geo := struct {
		Geo []Parser `json:"geo"`
	}{}
	err := json.Unmarshal(b, &geo)
	if err != nil {
		return err
	}

	g.Name = "geo"
	g.Children = geo.Geo
	return nil
}

func NewGeo(value string) *Geo {
	return &Geo{BasicContext{
		Name:     "geo",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
