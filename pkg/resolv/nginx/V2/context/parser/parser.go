package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
)

type Parser interface {
	Bytes() []byte
	Dump(cacher dump_cacher.DumpCacher) error
	GetType() parser_type.ParserType
	GetValue() string
	SetGlobalDeep(int)
	Query(words KeyWords) Parser
}
