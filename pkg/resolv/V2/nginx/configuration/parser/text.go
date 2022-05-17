package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Text struct {
	Value     string `json:"value"`
	position  parser_position.ParserPosition
	indention parser_indention.Indention
}

func (t Text) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(t.Position.ConfigIndents()+string(t.GetType()))
	// debug config Position end*/
	dumper.Write(t.position.Id(), []byte(t.indention.ConfigIndents()+t.string()))

	return nil
}

func (t Text) Bytes() []byte {
	return []byte(t.indention.GlobalIndents() + t.string())
}

func (t Text) GetType() parser_type.ParserType {
	return parser_type.TypeKey
}

func (t Text) GetValue() string {
	return t.Value
}

func (t *Text) setPosition(p string) error {
	t.position = parser_position.NewPosition(p)

	return nil
}

func (t Text) GetIndention() parser_indention.Indention {
	return t.indention
}

func (t Text) GetPosition() string {
	return t.position.Id()
}

func (t *Text) SetGlobalDeep(deep int) {
	t.indention.SetGlobalDeep(deep)
}

func (t *Text) Match(words KeyWords) bool {
	return words.Match(t)
}

func (t Text) string() string {
	return t.Value + "\n"
}

func NewText(v string, indention parser_indention.Indention) Parser {
	return &Text{
		Value:     v,
		indention: indention,
	}
}
