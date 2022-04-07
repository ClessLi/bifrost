package nginx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Load(path string) (string, Caches, error) {
	cs := NewCaches()
	ngDir, fileErr := filepath.Abs(filepath.Dir(path))
	if fileErr != nil {
		return "", nil, fileErr
	}
	relativePath := filepath.Base(path)
	configPath, loadErr := load(ngDir, relativePath, &cs)
	if loadErr != nil && loadErr != IsInCaches {
		return "", nil, loadErr
	}
	return configPath, cs, nil
}

func load(ngDir, relativePath string, caches *Caches) (configAbsPath string, err error) {
	configAbsPath = filepath.Join(ngDir, relativePath)
	if _, ok := (*caches)[configAbsPath]; ok {
		return configAbsPath, IsInCaches
	}
	f := NewConf(nil, configAbsPath)

	file, openErr := os.Open(configAbsPath)
	if openErr != nil {
		return "", openErr
	}
	defer file.Close()

	bytes, readErr := ioutil.ReadAll(file)
	if readErr != nil {
		return "", readErr
	}
	data := string(bytes)
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
			c := NewComment(ms[2], !strings.Contains(ms[1], "\n") && index != 0)
			if ct, ok := checkContext(lopen); ok {
				ct.AddByParser(c)
				lopen[0] = ct
			} else {
				f.AddByParser(c)
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
					uc.AddByParser(lc)
					lopen[0] = uc
				} else {
					f.AddByParser(lc)
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

			k, ckErr := checkInclude(k, ngDir, caches)
			if ckErr != nil {
				err = ckErr
				return false
			}

			if ct, isContext := checkContext(lopen); isContext {
				ct.AddByParser(k)
				lopen[0] = ct
			} else {
				f.AddByParser(k)
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
				return "", err
			}
			continue
		}

		break
	}

	err = caches.SetCache(f, file)

	return
}

func checkInclude(k *Key, dir string, cs *Caches) (Parser, error) {
	if k.Name == fmt.Sprintf("%s", TypeInclude) {
		// return NewInclude(dir, k.Value, allConfigs, configCaches)
		return NewInclude(dir, k.Value, cs)
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
