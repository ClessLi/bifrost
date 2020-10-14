package nginx

import (
	"io/ioutil"
)

type Order int

type Parser interface {
	String() []string
	Query(parserType, bool, ...string) Parser
	QueryByKeywords(Keywords) Parser
	QueryAll(parserType, bool, ...string) []Parser
	QueryAllByKeywords(Keywords) []Parser
	BitSize(Order, int) byte
	BitLen(Order) int
	Size(Order) int
}

type Config struct {
	BasicContext `json:"config"`
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

func (c *Config) Save() error {
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

func (c *Config) List() (ret Caches, err error) {
	ret = Caches{}
	err = ret.setCache(c, hashForGetList)
	if err != nil && err != IsInCaches {
		return nil, err
	}
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			subCaches, err := child.(Context).List()
			if err != nil {
				return nil, err
			} else if subCaches != nil {
				for _, cache := range subCaches {
					err = ret.setCache(cache.config, hashForGetList)
					if err != nil && err != IsInCaches {
						return nil, err
					}
				}
			}
		}
	}
	return ret, nil
}

func (c *Config) QueryAllByKeywords(kw Keywords) (parsers []Parser) {
	if c.filter(kw) {
		parsers = append(parsers, c)
	}
	return c.subQueryAll(parsers, kw)

}

func (c *Config) QueryByKeywords(kw Keywords) (parser Parser) {
	if c.filter(kw) {
		return c
	}
	return c.subQuery(kw)
}

func (c *Config) BitSize(_ Order, _ int) byte {
	return 0
}

func (c *Config) BitLen(_ Order) int {
	return 0
}

func (c *Config) Size(_ Order) int {
	return 0
}

func NewConf(conf []Parser, value string) *Config {
	f := &Config{BasicContext{
		Name:     TypeConfig,
		Value:    value,
		Children: conf,
	}}

	return f
}
