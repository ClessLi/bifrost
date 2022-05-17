package parser

import (
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Context interface {
	Parser
	Insert(parser Parser, index int) error
	Remove(index int) error
	Modify(parser Parser, index int) error
	Query(words KeyWords) (Context, int)
	QueryAll(words KeyWords) map[Context][]int
	GetChild(index int) (Parser, error)
	Len() int
}

func NewContext(value string, parserType parser_type.ParserType, indention parser_indention.Indention) Context {
	children := make([]Parser, 0)
	basicContext := BasicContext{
		Name:      parserType,
		Children:  children,
		indention: indention,
	}
	switch parserType { //nolint:exhaustive
	case parser_type.TypeConfig:
		basicContext.Value = value
		basicContext.Position = parser_position.NewPosition(value)
		// basicContext.Position = parser_position.NewParserPosition(Position.ConfigAbsPath(), Position.ConfigDeep(), 0)
		return &Config{basicContext}
	case parser_type.TypeEvents:
		return &Events{basicContext}
	case parser_type.TypeGeo:
		basicContext.Value = value

		return &Geo{basicContext}
	case parser_type.TypeHttp:
		return &Http{basicContext}
	case parser_type.TypeIf:
		basicContext.Value = value

		return &If{basicContext}
	case parser_type.TypeInclude:
		basicContext.Value = value

		return &Include{basicContext}
	case parser_type.TypeLimitExcept:
		basicContext.Value = value

		return &LimitExcept{basicContext}
	case parser_type.TypeLocation:
		basicContext.Value = value

		return &Location{basicContext}
	case parser_type.TypeMap:
		basicContext.Value = value

		return &Map{basicContext}
	case parser_type.TypeServer:
		return &Server{basicContext}
	case parser_type.TypeStream:
		return &Stream{basicContext}
	case parser_type.TypeTypes:
		return &Types{basicContext}
	case parser_type.TypeUpstream:
		basicContext.Value = value

		return &Upstream{basicContext}
	default:
		return nil
	}
}
