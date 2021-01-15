package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
)

type Comment struct {
	Comments string `json:"comments"`
	Inline   bool   `json:"inline"`
	position parser_position.ParserPosition
}

func (c Comment) Dump(cacher dump_cacher.DumpCacher) error {
	cacher.Write(c.position.ConfigAbsPath(), []byte(c.position.ConfigIndents()+c.string()))
	return nil
}

func (c Comment) Bytes() []byte {
	return []byte(c.position.GlobalIndents() + c.string())
}

func (c Comment) GetType() parser_type.ParserType {
	return parser_type.TypeComment
}

func (c Comment) GetValue() string {
	return c.Comments
}

func (c *Comment) SetGlobalDeep(deep int) {
	c.position.SetGlobalDeep(deep)
}

func (c *Comment) Query(words KeyWords) Parser {
	if words.Match(c) {
		return c
	}
	return nil
}

func (c Comment) string() string {
	return "# " + c.Comments + "\n"
}

func NewComment(comments string, inline bool, position parser_position.ParserPosition) Parser {
	return &Comment{
		Comments: comments,
		Inline:   inline,
		position: position,
	}
}
