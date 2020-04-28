package resolv

type LimitExcept struct {
	BasicContext `json:"limit_except"`
}

func NewLimitExcept(value string) *LimitExcept {
	return &LimitExcept{BasicContext{
		Name:     "limit_except",
		Value:    value,
		Children: nil,
	}}
}
