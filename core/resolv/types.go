package resolv

type Types struct {
	BasicContext
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     "types",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
