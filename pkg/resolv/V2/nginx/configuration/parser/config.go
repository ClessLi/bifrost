package parser

import (
	"bytes"

	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/dumper"
)

type Config struct {
	BasicContext `json:"config"`
}

func (c Config) Bytes() []byte {
	buff := bytes.NewBuffer([]byte(""))
	for _, child := range c.Children {
		if cmt, ok := child.(*Comment); ok && cmt.Inline {
			buff.Truncate(buff.Len() - 1)
		}
		child.SetGlobalDeep(c.indention.GlobalDeep())
		buff.Write(child.Bytes())
	}

	return buff.Bytes()
}

func (c Config) Dump(dumper dumper.Dumper) error {
	/*// debug config Position
	fmt.Println(c.Position.ConfigIndents()+string(c.GetType()))
	// debug config Position end*/
	for _, child := range c.Children {
		err := child.Dump(dumper)
		if err != nil {
			return err
		}
	}
	// return nil

	return dumper.Done(c.GetValue())
}
