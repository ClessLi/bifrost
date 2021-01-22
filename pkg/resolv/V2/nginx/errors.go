package nginx

import "errors"

var (
	KeywordStringError = errors.New("unknown keyword string")
	ErrNotFound        = errors.New("parser not found")
	ErrConfigObject    = errors.New("error config object")
)
