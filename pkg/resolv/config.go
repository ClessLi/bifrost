package resolv

import (
	"encoding/json"
	"io/ioutil"
)

type Parser interface {
	String() []string
	//UnmarshalJSON(b []byte) error
}
type Config struct {
	BasicContext
}

func (c *Config) Servers() []*Server {
	svrs := make([]*Server, 0)
	for _, child := range c.Children {
		svrs = appendServer(svrs, child)
	}
	return svrs
}

func appendServer(svrs []*Server, p Parser) []*Server {
	switch p.(type) {
	case *Server:
		return append(svrs, p.(*Server))
	case *Http:
		for _, child := range p.(*Http).Children {
			svrs = appendServer(svrs, child)
		}
		return svrs
	case *Include:
		for _, child := range p.(*Include).Children {
			for _, includechild := range child.(*Config).Children {
				svrs = appendServer(svrs, includechild)
			}
			//fmt.Println(len(child.(*Config).Servers()))
		}
		return svrs
	default:
		return svrs
	}
}

func (c *Config) Server() *Server {
	return c.Servers()[0]
}

func (c *Config) String() []string {
	ret := make([]string, 0)

	title := c.getTitle()
	ret = append(ret, "# "+title)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			ret = append(ret, child.String()[0])
		case Context:
			ret = append(ret, child.String()...)
		}
	}

	if ret != nil {
		ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	}

	ret = append(ret, "#End# "+c.Name+": "+c.Value+"}\n\n")

	return ret
}

func (c *Config) save() error {
	conf, derr := c.dump()
	if derr != nil {
		return derr
	}
	data := make([]byte, 0)
	for _, str := range conf {
		data = append(data, []byte(str)...)
	}

	werr := ioutil.WriteFile(c.Value, data, 0755)
	if werr != nil {
		return werr
	}

	return nil

}

func (c *Config) dump() ([]string, error) {
	ret := make([]string, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			ret = append(ret, child.String()...)
		case Context:
			strs, err := child.(Context).dump()

			if err != nil {
				return ret, err
			}

			ret = append(ret, strs...)
		default:
			ret = append(ret, child.String()...)
		}
	}
	return ret, nil
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path   string   `json:"path"`
		Config []Parser `json:"config"`
	}{Path: c.Value, Config: c.Children})
}

func (c *Config) UnmarshalJSON(b []byte) error {
	// TODO:json 反序列化方法
	var config struct {
		Path   string   `json:"path"`
		Config []Parser `json:"config"`
	}
	//err := unmarshal(b, c)
	err := json.Unmarshal(b, &config)
	if err != nil {
		return err
	}

	//err := json.Unmarshal(b, &config)
	//if err != nil {
	//	return err
	//}

	c.Name = "Config"
	c.Value = config.Path
	c.Children = config.Config

	return nil
}

//func childrenUnmarshal(maps [][]byte) (children []Parser, error) {
//	for _, child := range maps {
//
//	}
//	err := json.Unmarshal(b, p)
//	if err != nil {
//		return err
//	}

//}

func NewConf(conf []Parser, value string) *Config {
	return &Config{BasicContext{
		Name:     "Config",
		Value:    value,
		depth:    0,
		Children: conf,
	}}
}
