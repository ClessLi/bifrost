package parser

import (
	"bytes"

	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/dumper"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_indention"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/parser_type"
)

type LuaBlock struct {
	BasicContext `json:"upstream"`
	Name         string `json:"name"`
}

func NewLuaContext(name string, value string, indention parser_indention.Indention) Context {
	children := make([]Parser, 0)
	basicContext := BasicContext{
		Name:      parser_type.TypeLuaBlock,
		Children:  children,
		indention: indention,
		Value:     value,
	}

	return &LuaBlock{
		BasicContext: basicContext,
		Name:         name,
	}
}

func (l LuaBlock) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(l.indention.GlobalIndents() + l.headString()))
	for _, child := range l.Children {
		if cmt, ok := child.(*Comment); ok && cmt.Inline {
			buff.Truncate(buff.Len() - 1)
		}
		child.SetGlobalDeep(l.indention.GlobalDeep() + 1)
		buff.Write(child.Bytes())
	}
	buff.WriteString(l.indention.GlobalIndents() + l.tailString())

	return buff.Bytes()
}

func (l LuaBlock) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(l.Position.ConfigIndents()+string(l.GetType()))
	// debug config Position end*/
	dumper.Write(l.Position.Id(), []byte(l.indention.ConfigIndents()+l.headString()))
	for _, child := range l.Children {
		err := child.Dump(dumper)
		if err != nil {
			return err
		}
	}
	dumper.Write(l.Position.Id(), []byte(l.indention.ConfigIndents()+l.tailString()))

	return nil
}

func (l LuaBlock) headString() string {
	contextTitle := ""
	contextTitle += l.Name

	if l.Value != "" {
		contextTitle += " " + l.Value
	}

	contextTitle += " {\n"

	return contextTitle
}
