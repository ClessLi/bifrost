package resolv

type Upstream struct {
	BasicContext
}

func NewUpstream(value string) *Upstream {
	return &Upstream{BasicContext{
		Name:     "upstream",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
