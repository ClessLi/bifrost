package nginx

import "io/ioutil"

type Order int

type Parser interface {
	String(*Caches) []string
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

func (c *Config) Save() error {
	caches := NewCaches()
	dumps, derr := c.dump(c.Value, &caches)
	if derr != nil {
		return derr
	}
	data := make([]byte, 0)
	for path, lines := range dumps {

		for _, line := range lines {
			data = append(data, []byte(line)...)
		}

		werr := ioutil.WriteFile(path, data, 0755)
		if werr != nil {
			return werr
		}
	}

	return nil

}

func (c Config) getCaches() (Caches, error) {
	caches := NewCaches()
	dumps, derr := c.dump(c.Value, &caches)
	if derr != nil {
		return nil, derr
	}
	for path, lines := range dumps {
		cache, ok := caches[path]
		if !ok {
			continue
		}
		data := make([]byte, 0)

		for _, line := range lines {
			data = append(data, []byte(line)...)
		}
		hash, herr := getHash(cache.config.Value, data)
		if herr != nil {
			return nil, herr
		}
		cache.hash = hash
		caches[path] = cache
	}
	return caches, nil
}

func (c *Config) String(caches *Caches) []string {
	if _, ok := (*caches)[c.Value]; ok {
		return nil
	}
	ret := make([]string, 0)

	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			ret = append(ret, child.String(caches)[0])
		case Context:
			ret = append(ret, child.String(caches)...)
		}
	}

	if ret != nil {
		ret[len(ret)-1] = RegEndWithCR.ReplaceAllString(ret[len(ret)-1], "}\n")
	}

	(*caches)[c.Value], _ = newCache(c, hashForString)

	return ret
}

func (c *Config) dump(_ string, caches *Caches) (map[string][]string, error) {
	if _, ok := (*caches)[c.Value]; ok {
		return nil, IsInCaches
	}

	dumps := make(map[string][]string)
	var err error

	err = caches.setCache(c, hashForDump)
	if err != nil {
		return nil, err
	}

	dmp := make([]string, 0)
	for _, child := range c.Children {
		switch child.(type) {
		case *Key, *Comment:
			dmp = append(dmp, child.String(caches)...)
		case Context:
			dumps, err = child.(Context).dump(c.Value, caches)

			if err != nil {
				return nil, err
			}

			if d, ok := dumps[c.Value]; ok {
				dmp = append(dmp, d...)
			}
		default:
			dmp = append(dmp, child.String(caches)...)
		}
	}
	dumps[c.Value] = dmp
	return dumps, nil
}

func (c *Config) List() (caches Caches, err error) {
	caches = NewCaches()
	err = caches.setCache(c, hashForGetList)
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
					err = caches.setCache(cache.config, hashForGetList)
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
