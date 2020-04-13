package resolv

type Http struct {
	BasicContext `json:"http"`
}

func NewHttp() *Http {
	return &Http{BasicContext{
		Name:     "http",
		Value:    "",
		depth:    0,
		Children: nil,
	}}
}
