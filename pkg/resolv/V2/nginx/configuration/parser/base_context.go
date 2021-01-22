package parser

import (
	"bytes"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type BasicContext struct {
	Name      parser_type.ParserType `json:"-"`
	Value     string                 `json:"value,omitempty"`
	Children  []Parser               `json:"param,omitempty"`
	Position  parser_position.ParserPosition
	indention parser_indention.Indention
}

func (b BasicContext) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(b.indention.GlobalIndents() + b.headString()))
	for _, child := range b.Children {
		child.SetGlobalDeep(b.indention.GlobalDeep() + 1)
		buff.Write(child.Bytes())
	}
	buff.WriteString(b.indention.GlobalIndents() + b.tailString())
	return buff.Bytes()
}

func (b BasicContext) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(b.Position.ConfigIndents()+string(b.GetType()))
	// debug config Position end*/
	dumper.Write(b.Position.Id(), []byte(b.indention.ConfigIndents()+b.headString()))
	for _, child := range b.Children {
		err := child.Dump(dumper)
		if err != nil {
			return err
		}
	}
	dumper.Write(b.Position.Id(), []byte(b.indention.ConfigIndents()+b.tailString()))
	return nil
}

func (b BasicContext) GetType() parser_type.ParserType {
	return b.Name
}

func (b BasicContext) GetValue() string {
	return b.Value
}

func (b BasicContext) GetPosition() string {
	return b.Position.Id()
}

func (b *BasicContext) setPosition(p string) error {
	b.Position = parser_position.NewPosition(p)
	for i := range b.Children {
		err := b.Children[i].setPosition(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b BasicContext) GetIndention() parser_indention.Indention {
	return b.indention
}

func (b BasicContext) SetGlobalDeep(deep int) {
	b.indention.SetGlobalDeep(deep)
}

func (b *BasicContext) Insert(parser Parser, index int) error {
	if index > len(b.Children) {
		return ErrIndexOutOfRange
	}

	if b.Position != nil && parser.GetType() != parser_type.TypeConfig {
		err := parser.setPosition(b.Position.Id())
		if err != nil {
			return err
		}
	}
	n := b.Len()
	b.Children = append(b.Children, nil)
	for i := n; i > index; i-- {
		b.Children[i] = b.Children[i-1]
	}
	b.Children[index] = parser
	return nil
}

func (b *BasicContext) Remove(index int) error {
	if index > len(b.Children)-1 {
		return ErrIndexOutOfRange
	}
	b.Children = append(b.Children[:index], b.Children[index+1:]...)
	return nil
}

func (b *BasicContext) Modify(parser Parser, index int) error {
	if index > len(b.Children)-1 {
		return ErrIndexOutOfRange
	}
	err := b.Remove(index)
	if err != nil {
		return err
	}
	return b.Insert(parser, index)
}

func (b BasicContext) Match(words KeyWords) bool {
	return words.Match(&b)
}

func (b BasicContext) Query(words KeyWords) (Context, int) {
	for idx, child := range b.Children {
		if child.Match(words) {
			return &b, idx
		}
		if c, ok := child.(Context); ok {
			ctx, index := c.Query(words)
			if ctx != nil {
				return ctx, index
			}
		}
	}
	return nil, 0
}

func (b BasicContext) QueryAll(words KeyWords) map[Context][]int {
	result := make(map[Context][]int)
	for idx, child := range b.Children {
		if child.Match(words) {
			result[&b] = append(result[&b], idx)
		}
		if c, ok := child.(Context); ok {
			subResult := c.QueryAll(words)
			for ctx, indexes := range subResult {
				result[ctx] = append(result[ctx], indexes...)
			}
		}
	}
	return result
}

func (b BasicContext) GetChild(index int) (Parser, error) {
	if index < 0 || index >= b.Len() {
		return nil, ErrIndexOutOfRange
	}
	return b.Children[index], nil
}

func (b BasicContext) Len() int {
	return len(b.Children)
}

func (b BasicContext) headString() string {
	contextTitle := ""
	contextTitle += b.Name.String()

	if b.Value != "" {
		contextTitle += " " + b.Value
	}

	contextTitle += " {\n"
	return contextTitle
}

func (b BasicContext) tailString() string {
	return "}\n"
}
