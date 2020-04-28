package resolv

type Upstream struct {
	BasicContext `json:"upstream"`
}

func NewUpstream(value string) *Upstream {
	return &Upstream{BasicContext{
		Name:     "upstream",
		Value:    value,
		Children: nil,
	}}
}
