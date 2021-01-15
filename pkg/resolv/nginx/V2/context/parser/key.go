package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
)

type Key struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	position parser_position.ParserPosition
}

func (k Key) Dump(cacher dump_cacher.DumpCacher) error {
	cacher.Write(k.position.ConfigAbsPath(), []byte(k.position.ConfigIndents()+k.string()))
	return nil
}

func (k Key) Bytes() []byte {
	return []byte(k.position.GlobalIndents() + k.string())
}

func (k Key) GetType() parser_type.ParserType {
	return parser_type.TypeKey
}

func (k Key) GetValue() string {
	return k.Name + " " + k.Value
}

func (k *Key) SetGlobalDeep(deep int) {
	k.position.SetGlobalDeep(deep)
}

func (k *Key) Query(words KeyWords) Parser {
	if words.Match(k) {
		return k
	}
	return nil
}

func (k Key) string() string {
	return k.Name + " " + k.Value + ";\n"
}

func NewKey(k, v string, position parser_position.ParserPosition) Parser {
	return &Key{
		Name:     k,
		Value:    v,
		position: position,
	}
}
