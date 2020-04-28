package resolv

type If struct {
	BasicContext `json:"if"`
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     "if",
		Value:    value,
		Children: nil,
	}}
}
