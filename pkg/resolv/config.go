package resolv

import (
	"io/ioutil"
)

type Parser interface {
	String() []string
	Filter(KeyWords) []Parser
}

type Config struct {
	BasicContext `json:"config"`
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

func NewConf(conf []Parser, value string) *Config {
	return &Config{BasicContext{
		Name:     "Config",
		Value:    value,
		depth:    0,
		Children: conf,
	}}
}
