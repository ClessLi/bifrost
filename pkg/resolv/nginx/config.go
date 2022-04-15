package nginx

import (
	"io/ioutil"
)

type Order int

type Parser interface {
	String() []string
	string(*Caches, int) []string
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

func (c *Config) Save() (Caches, error) {
	caches := NewCaches()
	dumps, derr := c.dump(c.Value, &caches, 0)
	if derr != nil {
		return nil, derr
	}
	for path, lines := range dumps {
		data := make([]byte, 0)

		for _, line := range lines {
			data = append(data, []byte(line)...)
		}

		werr := ioutil.WriteFile(path, data, 0755)

		cacheErr := caches.SetCache(caches[path].config, data)
		if cacheErr != nil {
			return nil, cacheErr
		}

		if werr != nil {
			return nil, werr
		}
	}

	return caches, nil
}

func (c Config) getCaches() (Caches, error) {
	caches := NewCaches()
	dumps, derr := c.dump(c.Value, &caches, 0)
	if derr != nil {
		return nil, derr
	}
	for path, lines := range dumps {
		data := make([]byte, 0)

		for _, line := range lines {
			data = append(data, []byte(line)...)
		}

		cacheErr := caches.SetCache(caches[path].config, data)
		if cacheErr != nil {
			return nil, cacheErr
		}
	}
	return caches, nil
}

func (c Config) String() []string {
	caches := NewCaches()
	return c.string(&caches, 0)
}

func (c *Config) string(caches *Caches, deep int) []string {
	if _, ok := (*caches)[c.Value]; ok {
		return nil
	}
	ret := make([]string, 0)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			ret = append(ret, child.string(caches, deep)[0])
		case Context:
			ret = append(ret, child.string(caches, deep)...)
		}
	}

	if ret != nil {
		ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	}

	(*caches)[c.Value], _ = newCache(c, hashForString)

	return ret
}

//func (c Config) Dump() (map[string][]byte, error) {
//	caches := NewCaches()
//	strDmps, err := c.dump("", &caches, 0)
//	if err != nil {
//		return nil, err
//	}
//
//	dumps := make(map[string][]byte)
//	for path, strings := range strDmps {
//		for _, s := range strings {
//			dumps[path] = append(dumps[path], []byte(s)...)
//		}
//	}
//
//	return dumps, nil
//}

func (c *Config) dump(_ string, caches *Caches, deep int) (map[string][]string, error) {
	dumps := make(map[string][]string)
	if caches.IsCached(c.Value) {
		return dumps, IsInCaches
	}

	var err error

	err = caches.SetCache(c, hashForDumpTemp)
	if err != nil {
		return nil, err
	}

	dmp := make([]string, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			dmp = append(dmp, child.string(caches, deep)...)
		case Context:

			newDumps, err := child.(Context).dump(c.Value, caches, deep)

			if err != nil && err != IsInCaches {
				return nil, err
			} else if err == IsInCaches {
				break
			}

			for dmpPath, data := range newDumps {
				if _, ok := dumps[dmpPath]; ok || dmpPath == c.Value {
					continue
				}
				dumps[dmpPath] = data
			}

			if d, ok := newDumps[c.Value]; ok {
				dmp = append(dmp, d...)
			}
		default:
			lines := child.string(caches, deep)
			if lines != nil {
				dmp = append(dmp, lines...)
			}
		}
	}
	dumps[c.Value] = dmp
	return dumps, nil
}

func (c *Config) List() (caches Caches, err error) {
	caches = NewCaches()
	err = caches.SetCache(c, hashForGetList)
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
					err = caches.SetCache(cache.config, hashForGetList)
					if err != nil && err != IsInCaches {
						return nil, err
					}
				}
			}
		}
	}
	return
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
