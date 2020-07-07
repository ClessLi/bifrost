package nginx

import (
	"io/ioutil"
)

type Order int

type Parser interface {
	String() []string
	Query(KeyWords) Parser
	QueryAll(KeyWords) []Parser
	BitSize(Order, int) byte
	BitLen(Order) int
	Size(Order) int
}

type Config struct {
	BasicContext `json:"config"`
}

func (c *Config) String() []string {
	ret := make([]string, 0)

	// 暂取消config对象输出自身对象信息
	//title := c.getTitle()
	//ret = append(ret, "# "+title)

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

	// 暂取消config对象输出自身对象信息
	//ret = append(ret, "#End# "+c.Name+": "+c.Value+"}\n\n")

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

func (c *Config) List() (ret []string, err error) {
	//absPath, err := filepath.Abs(c.Value)
	//if err != nil {
	//	return
	//}
	//ret = []string{absPath}
	ret = []string{c.Value}
	//fmt.Println(c.Value)
	//fmt.Println("list: ", ret)
	for _, child := range c.Children {
		switch child.(type) {
		case Context:
			l, err := child.(Context).List()
			if err != nil {
				return nil, err
			} else if l != nil {
				ret = append(ret, l...)
			}
		}
	}
	return ret, nil
}

func (c *Config) QueryAll(kw KeyWords) (parsers []Parser) {
	if c.filter(kw) {
		parsers = append(parsers, c)
	}
	if kw.IsRec {
		return c.subQueryAll(parsers, kw)
	} else {
		return
	}

}

func (c *Config) Query(kw KeyWords) (parser Parser) {
	if c.filter(kw) {
		return c
	}

	if kw.IsRec {
		return c.subQuery(kw)
	} else {
		return
	}
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
	// 确定*Config的唯一性，防止多次加载
	//for _, c := range configs {
	//	if c.Value == value {
	//		return c, ErrConfigIsExist
	//	}
	//}
	f := &Config{BasicContext{
		Name:     TypeConfig,
		Value:    value,
		Children: conf,
	}}

	// 确保*Config的唯一性，将新加载的*Config加入configs
	//configs = append(configs, f)
	//return f, nil
	return f
}
