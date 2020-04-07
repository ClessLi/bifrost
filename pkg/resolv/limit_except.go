package resolv

import "encoding/json"

type LimitExcept struct {
	BasicContext
}

func (le *LimitExcept) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, le)
}

func NewLimitExcept(value string) *LimitExcept {
	return &LimitExcept{BasicContext{
		Name:     "limit_except",
		Value:    value,
		depth:    0,
		Children: nil,
	}}
}
