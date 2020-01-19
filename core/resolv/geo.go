package resolv

type Geo struct {
	BasicContext
}

func NewGeo() *Geo {
	return &Geo{BasicContext{
		Name:     "geo",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
