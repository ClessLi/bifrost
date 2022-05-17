package parser

import (
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Parser interface {
	Bytes() []byte
	Dump(dumper dumper.Dumper) error
	GetType() parser_type.ParserType
	// key.GetValue:
	// return key's name and key's value
	GetValue() string
	GetIndention() parser_indention.Indention
	SetGlobalDeep(int)
	GetPosition() string
	setPosition(string) error
	Match(words KeyWords) bool
}
