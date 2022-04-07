package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
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
	// match parserType
	matched := parser.GetType() == k.parserType
	// match main key word
	if matched {
		matched = false //nolint:wastedassign,ineffassign
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

func NewKeyWords(pType parser_type.ParserType, isReg bool, value ...string) (KeyWords, error) {
	var kw keyWord
	if value != nil {
		switch pType { //nolint:exhaustive
		case parser_type.TypeComment,
			parser_type.TypeKey,
			parser_type.TypeConfig,
			parser_type.TypeGeo,
			parser_type.TypeIf,
			parser_type.TypeLimitExcept,
			parser_type.TypeLocation,
			parser_type.TypeMap,
			parser_type.TypeUpstream:
			kw = keyWord{
				parserType: pType,
				value:      value[0],
				isReg:      isReg,
			}
		case parser_type.TypeEvents,
			parser_type.TypeHttp,
			parser_type.TypeServer,
			parser_type.TypeStream,
			parser_type.TypeTypes:
			kw = keyWord{
				parserType: pType,
				value:      "",
				isReg:      false,
			}
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}
	} else {
		switch pType { //nolint:exhaustive
		case parser_type.TypeEvents, parser_type.TypeHttp, parser_type.TypeServer, parser_type.TypeStream, parser_type.TypeTypes:
			kw = keyWord{
				parserType: pType,
				value:      "",
				isReg:      false,
			}
		default:
			return nil, fmt.Errorf("unknown nginx context type: %s", pType)
		}
	}

	return &kw, nil
}
