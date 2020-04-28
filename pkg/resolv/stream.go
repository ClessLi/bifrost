package resolv

type Stream struct {
	BasicContext `json:"stream"`
}

func NewStream() *Stream {
	return &Stream{BasicContext{
		Name:     "stream",
		Value:    "",
		Children: nil,
	}}
}
