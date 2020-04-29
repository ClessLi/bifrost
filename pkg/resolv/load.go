package resolv

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func Load(path string) (*Config, error) {

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return load(absPath)
}

func load(path string) (*Config, error) {

	dir, dirErr := filepath.Abs(filepath.Dir(path))
	if dirErr != nil {
		return nil, dirErr
	}

	f, err := NewConf(nil, path)
	if err == ErrConfigIsExist {
		return f, nil
	} else if err != nil && err != ErrConfigIsExist {
		return nil, err
	}

	data, err := readConf(path)
	if err != nil {
		return nil, err
	}

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

			k := checkInclude(k, dir)

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
			continue
		}

		break
	}

	return f, nil
}

func checkInclude(k *Key, dir string) Parser {
	if k.Name == TypeInclude {
		return NewInclude(dir, k.Value)
	}
	return k
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
