package resolv

type Parser interface {
	String() []string
}
type Config struct {
	BasicContext
}

func (c *Config) Servers() []Parser {
	svrs := make([]Parser, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Server:
			svrs = append(svrs, child)
		}
	}
	return svrs
}

func (c *Config) Server() interface{} {
	return c.Servers()[0]
}

func (c *Config) String() []string {
	ret := make([]string, 0)
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

	return ret
}

func NewConf(conf []Parser, value string) *Config {
	return &Config{BasicContext{
		Name:     "Config",
		Value:    value,
		depth:    0,
		Children: conf,
	}}
}
