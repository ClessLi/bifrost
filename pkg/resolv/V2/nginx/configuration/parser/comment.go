package parser

import (
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type Comment struct {
	Comments  string `json:"comments"`
	Inline    bool   `json:"inline"`
	position  parser_position.ParserPosition
	indention parser_indention.Indention
}

func (c Comment) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(c.Position.ConfigIndents()+string(c.GetType()))
	// debug config Position end*/
	if c.Inline {
		err := dumper.Truncate(c.position.Id(), dumper.Len(c.position.Id())-1)
		if err != nil {
			return err
		}
	}
	dumper.Write(c.position.Id(), []byte(c.string()))
	return nil
}

func (c Comment) Bytes() []byte {
	return []byte(c.string())
}

func (c Comment) GetType() parser_type.ParserType {
	return parser_type.TypeComment
}

func (c Comment) GetValue() string {
	return c.Comments
}

func (c *Comment) GetPosition() string {
	return c.position.Id()
}

func (c *Comment) setPosition(p string) error {
	c.position = parser_position.NewPosition(p)
	return nil
}

func (c Comment) GetIndention() parser_indention.Indention {
	return c.indention
}

func (c *Comment) SetGlobalDeep(deep int) {
	c.indention.SetGlobalDeep(deep)
}

func (c *Comment) Match(words KeyWords) bool {
	return words.Match(c)
}

func (c Comment) string() string {
	if c.Inline {
		return "  # " + c.Comments + "\n"
	}
	return c.indention.GlobalIndents() + "# " + c.Comments + "\n"
}

func NewComment(comments string, inline bool, indention parser_indention.Indention) Parser {
	return &Comment{
		Comments:  comments,
		Inline:    inline,
		indention: indention,
	}
}
