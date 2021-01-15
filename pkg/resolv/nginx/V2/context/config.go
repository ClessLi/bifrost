package context

import (
	"bytes"
	"github.com/ClessLi/bifrost/pkg/resolv/nginx/V2/dump_cacher"
)

type Config struct {
	BasicContext
}

func (c Config) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(""))
	for _, child := range c.Children {
		child.SetGlobalDeep(c.Position.GlobalDeep())
		buff.Write(child.Bytes())
	}
	return buff.Bytes()
}

func (c Config) Dump(cacher dump_cacher.DumpCacher) error {
	for _, child := range c.Children {
		err := child.Dump(cacher)
		if err != nil {
			return err
		}
	}
	return nil
}
