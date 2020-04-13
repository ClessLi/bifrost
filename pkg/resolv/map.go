package resolv

type Map struct {
	BasicContext `json:"map"`
}

func NewMap(value string) *Map {
	return &Map{BasicContext{
		Name:     "map",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
