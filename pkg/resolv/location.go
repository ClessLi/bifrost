package resolv

type Location struct {
	BasicContext `json:"location"`
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     "location",
		Value:    value,
		Children: nil,
	}}
}
