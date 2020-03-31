package resolv

type Geo struct {
	BasicContext `json:"geo"`
}

func NewGeo(value string) *Geo {
	return &Geo{BasicContext{
		Name:     "geo",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
