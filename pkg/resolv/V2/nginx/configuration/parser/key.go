package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Key struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	position  parser_position.ParserPosition
	indention parser_indention.Indention
}

func (k Key) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(k.Position.ConfigIndents()+string(k.GetType()))
	// debug config Position end*/
	dumper.Write(k.position.Id(), []byte(k.indention.ConfigIndents()+k.string()))
	return nil
}

func (k Key) Bytes() []byte {
	return []byte(k.indention.GlobalIndents() + k.string())
}

func (k Key) GetType() parser_type.ParserType {
	return parser_type.TypeKey
}

func (k Key) GetValue() string {
	return k.Name + " " + k.Value
}

func (k *Key) setPosition(p string) error {
	k.position = parser_position.NewPosition(p)
	return nil
}

func (k Key) GetIndention() parser_indention.Indention {
	return k.indention
}

func (k Key) GetPosition() string {
	return k.position.Id()
}

func (k *Key) SetGlobalDeep(deep int) {
	k.indention.SetGlobalDeep(deep)
}

func (k *Key) Match(words KeyWords) bool {
	return words.Match(k)
}

func (k Key) string() string {
	return k.Name + " " + k.Value + ";\n"
}

func NewKey(k, v string, indention parser_indention.Indention) Parser {
	return &Key{
		Name:      k,
		Value:     v,
		indention: indention,
	}
}
