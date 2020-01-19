package resolv

type Location struct {
	BasicContext
}

func NewLocation(value string) *Location {
	return &Location{BasicContext{
		Name:     "location",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
