package resolv

type Types struct {
	BasicContext `json:"types"`
}

func NewTypes() *Types {
	return &Types{BasicContext{
		Name:     "types",
		Value:    "",
		Children: nil,
	}}
}
