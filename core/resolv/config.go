package resolv

type Parser interface {
	String() []string
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
	default:
		return svrs
	}
}

func (c *Config) Server() *Server {
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
