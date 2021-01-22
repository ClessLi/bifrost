package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Parser interface {
	Bytes() []byte
	Dump(dumper dumper.Dumper) error
	GetType() parser_type.ParserType
	GetValue() string
	GetIndention() parser_indention.Indention
	SetGlobalDeep(int)
	GetPosition() string
	setPosition(string) error
	Match(words KeyWords) bool
}
