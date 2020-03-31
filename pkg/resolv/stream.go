package resolv

type Stream struct {
	BasicContext
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     "stream",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
