package parser

import (
	"bytes"

	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/parser_position"
)

type Include struct {
	BasicContext `json:"include"`
}

func (i Include) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(""))
	for _, child := range i.Children {
		if cmt, ok := child.(*Comment); ok && cmt.Inline {
			buff.Truncate(buff.Len() - 1)
		}
		child.SetGlobalDeep(i.indention.GlobalDeep())
		buff.Write(child.Bytes())
	}

	return buff.Bytes()
}

func (i Include) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(i.Position.ConfigIndents()+string(i.GetType()))
	// debug config Position end*/

	// dump itself
	dumper.Write(i.Position.Id(), []byte(i.indention.ConfigIndents()+"include "+i.Value+";\n"))

	// dump included configs
	for _, child := range i.Children {
		err := child.Dump(dumper)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Include) Insert(parser Parser, index int) error {
	if _, ok := parser.(*Config); !ok {
		return ErrInsertParserTypeError
	}

	return i.BasicContext.Insert(parser, index)

	//if index > len(i.Children) {
	//	return ErrIndexOutOfRange
	//}
	//n := i.Len()
	//i.Children = append(i.Children, nil)
	//for j := n; j > index; j-- {
	//	i.Children[j] = i.Children[j-1]
	//}
	//i.Children[index] = parser
	//return nil
}

func (i *Include) Modify(parser Parser, index int) error {
	if _, ok := parser.(*Config); !ok {
		return ErrInsertParserTypeError
	}

	return i.BasicContext.Modify(parser, index)
}

func (i *Include) setPosition(p string) error {
	i.Position = parser_position.NewPosition(p)

	return nil
}
