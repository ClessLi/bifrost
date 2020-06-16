package resolv

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func Load(path string) (*Config, error) {
	ngDir, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		return nil, err
	}
	relativePath := filepath.Base(path)
	var configCaches []string
	return load(ngDir, relativePath, &configCaches)
}

func load(ngDir, relativePath string, configCaches *[]string) (conf *Config, err error) {
	absPath := filepath.Join(ngDir, relativePath)
	if inCaches(absPath, configCaches) {
		return nil, fmt.Errorf("config '%s' is already loaded", absPath)
	}

	f, err := NewConf(nil, absPath)
	if err == ErrConfigIsExist {
		return f, nil
	} else if err != nil && err != ErrConfigIsExist {
		return nil, err
	}

	data, err := readConf(absPath)
	if err != nil {
		return nil, err
	}
	//var newCaches []string
	//copy(newCaches, *configCaches)
	newCaches := *configCaches
	//configCaches = append(configCaches, absPath)
	newCaches = append(newCaches, absPath)

	index := 0
	var lopen []Parser

	parseContext := func(reg *regexp.Regexp) bool {
		var context Parser

		if m := reg.FindStringIndex(data[index:]); m != nil {
			switch reg {
			case RegEventsHead:
				context = NewEvents()
			case RegHttpHead:
				context = NewHttp()
			case RegStreamHead:
				context = NewStream()
			case RegServerHead:
				context = NewServer()
			case RegLocationHead:
				context = NewLocation(reg.FindStringSubmatch(data[index:])[1])
			case RegIfHead:
				context = NewIf(reg.FindStringSubmatch(data[index:])[1])
			case RegUpstreamHead:
				context = NewUpstream(reg.FindStringSubmatch(data[index:])[1])
			case RegGeoHead:
				context = NewGeo(reg.FindStringSubmatch(data[index:])[1])
			case RegMapHead:
				context = NewMap(reg.FindStringSubmatch(data[index:])[1])
			case RegLimitExceptHead:
				context = NewLimitExcept(reg.FindStringSubmatch(data[index:])[1])
			case RegTypesHead:
				context = NewTypes()
			}
			if context != nil {
				lopen = append([]Parser{context}, lopen...)
				index += m[len(m)-1]
				return true
			}
		}
		return false
	}

	parseComment := func(reg *regexp.Regexp) bool {
		if ms := reg.FindStringSubmatch(data[index:]); len(ms) == 3 {
			m := reg.FindStringIndex(data[index:])
			c := NewComment(ms[2], !strings.Contains(ms[1], "\n"))
			if ct, ok := checkContext(lopen); ok {
				ct.Add(c)
				lopen[0] = ct
			} else {
				f.Add(c)
			}
			index += m[len(m)-1] - 1
			return true
		} else {
			return false
		}
	}

	parseContextEnd := func(reg *regexp.Regexp) bool {
		if m := reg.FindStringIndex(data[index:]); m != nil {
			if lc, isLowerContext := checkContext(lopen); isLowerContext {
				lopen = lopen[1:]
				if uc, isUpperContext := checkContext(lopen); isUpperContext {
					uc.Add(lc)
					lopen[0] = uc
				} else {
					f.Add(lc)
				}
			}
			index += m[len(m)-1]
			return true
		} else {
			return false
		}
	}

	parseKey := func(reg *regexp.Regexp) bool {
		var k *Key
		if m := reg.FindStringIndex(data[index:]); m != nil {
			ms := reg.FindStringSubmatch(data[index:])
			switch reg {
			case RegKey:
				k = NewKey(ms[1], "")
			case RegKeyValue:
				k = NewKey(ms[1], ms[2])
			}

			k, ckErr := checkInclude(k, ngDir, &newCaches)
			if ckErr != nil {
				err = ckErr
				return false
			}

			if ct, isContext := checkContext(lopen); isContext {
				ct.Add(k)
				lopen[0] = ct
			} else {
				f.Add(k)
			}
			index += m[len(m)-1]
			return true
		}
		return false
	}

	for {

		switch {
		case parseContext(RegEventsHead),
			parseContext(RegHttpHead),
			parseContext(RegStreamHead),
			parseContext(RegServerHead),
			parseContext(RegLocationHead),
			parseContext(RegIfHead),
			parseContext(RegUpstreamHead),
			parseContext(RegGeoHead),
			parseContext(RegMapHead),
			parseContext(RegLimitExceptHead),
			parseContext(RegTypesHead),
			parseComment(RegCommentHead),
			parseContextEnd(RegContextEnd),
			parseKey(RegKeyValue),
			parseKey(RegKey):
			if err != nil {
				return nil, err
			}
			continue
		}

		break
	}

	return f, err
}

func inCaches(path string, caches *[]string) bool {
	for _, cache := range *caches {
		if path == cache {
			return true
		}
	}
	return false
}

func checkInclude(k *Key, dir string, configCaches *[]string) (Parser, error) {
	if k.Name == TypeInclude {
		return NewInclude(dir, k.Value, configCaches)
	}
	return k, nil
}

func checkContext(lopen []Parser) (Context, bool) {
	if len(lopen) > 0 {
		ct, isContext := lopen[0].(Context)
		return ct, isContext
	} else {
		return nil, false
	}
}

func readConf(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 解析、加载配置文件对象后，方便清除相关缓存
//func ReleaseConfigsCache() {
//	configs = []*Config{}
//}
