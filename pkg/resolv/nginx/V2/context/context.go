package context

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
)

type Context interface {
	parser.Parser
	Insert(parser parser.Parser, index int) error
	Remove(index int) error
	Modify(parser parser.Parser, index int) error
	QueryAll(words parser.KeyWords) []parser.Parser
	Len() int
}

func NewContext(value string, parserType parser_type.ParserType, position parser_position.ParserPosition) Context {
	children := make([]parser.Parser, 0)
	basicContext := BasicContext{
		Name:     parserType,
		Children: children,
		Position: position,
	}
	switch parserType {
	case parser_type.TypeConfig:
		basicContext.Value = value
		basicContext.Position = parser_position.NewParserPosition(position.ConfigAbsPath(), position.ConfigDeep(), 0)
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
		return &location{basicContext}
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
