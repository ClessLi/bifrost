package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
	"regexp"
	"strings"
)

type KeyWords interface {
	Match(parser Parser) bool
}

type keyWord struct {
	parserType parser_type.ParserType
	value      string
	isReg      bool
}

func (k keyWord) Match(parser Parser) bool {
	matched := parser.GetType() == k.parserType
	if matched {
		matched = false
		if k.isReg {
			var err error
			matched, err = regexp.MatchString(k.value, parser.GetValue())
			if err != nil {
				return false
			}
		} else {
			matched = strings.EqualFold(k.value, parser.GetValue())
		}
	}
	return matched
}
