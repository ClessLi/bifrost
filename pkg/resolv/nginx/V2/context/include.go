package context

import (
	"bytes"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/context/parser"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/config_graph"
)

type Include struct {
	BasicContext
	configGraph *config_graph.Graph
}

func (i Include) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(""))
	for _, child := range i.Children {
		child.SetGlobalDeep(i.Position.GlobalDeep())
		buff.Write(child.Bytes())
	}
	return buff.Bytes()
}

func (i Include) Dump(cacher dump_cacher.DumpCacher) error {
	// dump itself
	cacher.Write(i.Position.ConfigAbsPath(), []byte(i.Position.ConfigIndents()+i.Value+";\n"))
	for _, child := range i.Children {
		err := child.Dump(cacher)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Include) Insert(parser parser.Parser, index int) error {
	if _, ok := parser.(*Config); !ok {
		return ErrInsertParserTypeError
	}
	return i.BasicContext.Insert(parser, index)
}

func (i *Include) Modify(parser parser.Parser, index int) error {
	if _, ok := parser.(*Config); !ok {
		return ErrInsertParserTypeError
	}
	return i.BasicContext.Modify(parser, index)
}
