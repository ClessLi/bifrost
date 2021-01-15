package context

import (
	"bytes"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_position"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/parser_type"
)

type BasicContext struct {
	Name     parser_type.ParserType `json:"-"`
	Value    string                 `json:"value,omitempty"`
	Children []parser.Parser        `json:"param,omitempty"`
	Position parser_position.ParserPosition
}

func (b BasicContext) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(b.Position.GlobalIndents() + b.headString()))
	for _, child := range b.Children {
		child.SetGlobalDeep(b.Position.GlobalDeep() + 1)
		buff.Write(child.Bytes())
	}
	buff.WriteString(b.Position.GlobalIndents() + b.tailString())
	return buff.Bytes()
}

func (b BasicContext) Dump(cacher dump_cacher.DumpCacher) error {
	cacher.Write(b.Position.ConfigAbsPath(), []byte(b.Position.ConfigIndents()+b.headString()))
	for _, child := range b.Children {
		err := child.Dump(cacher)
		if err != nil {
			return err
		}
	}
	cacher.Write(b.Position.ConfigAbsPath(), []byte(b.Position.ConfigIndents()+b.tailString()))
	return nil
}

func (b BasicContext) GetType() parser_type.ParserType {
	return b.Name
}

func (b BasicContext) GetValue() string {
	return b.Value
}

func (b BasicContext) SetGlobalDeep(deep int) {
	b.Position.SetGlobalDeep(deep)
}

func (b *BasicContext) Insert(parser parser.Parser, index int) error {
	if index > len(b.Children) {
		return ErrIndexOutOfRange
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

func (b *BasicContext) Modify(parser parser.Parser, index int) error {
	if index > len(b.Children)-1 {
		return ErrIndexOutOfRange
	}
	err := b.Remove(index)
	if err != nil {
		return err
	}
	return b.Insert(parser, index)
}

func (b BasicContext) Query(words parser.KeyWords) parser.Parser {
	if words.Match(b) {
		return &b
	}
	for _, child := range b.Children {
		subParser := child.Query(words)
		if subParser != nil {
			return subParser
		}
	}
	return nil
}

func (b BasicContext) QueryAll(words parser.KeyWords) []parser.Parser {
	parsers := make([]parser.Parser, 0)
	if words.Match(b) {
		parsers = append(parsers, &b)
	}
	for _, child := range b.Children {
		if c, ok := child.(Context); ok {
			subParsers := c.QueryAll(words)
			if subParsers != nil {
				parsers = append(parsers, subParsers...)
			}
		} else {
			subParser := child.Query(words)
			if subParser != nil {
				parsers = append(parsers, subParser)
			}
		}
	}
	return parsers
}

func (b BasicContext) Len() int {
	return len(b.Children)
}

func (b BasicContext) headString() string {
	contextTitle := ""
	/*for i := 0; i < c.depth; i++ {
		contextTitle += INDENT
	}*/
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
