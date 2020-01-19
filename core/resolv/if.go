package resolv

type If struct {
	BasicContext
}

func NewIf(value string) *If {
	return &If{BasicContext{
		Name:     "if",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
