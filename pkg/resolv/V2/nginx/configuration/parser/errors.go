package parser

import "errors"

var (
	ErrIndexOutOfRange       = errors.New("index out of range")
	ErrInsertParserTypeError = errors.New("insert parser type error")
	ErrNullPosition          = errors.New("null position")
)
