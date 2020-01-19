package core

import (
	"github.com/ClessLi/go-nginx-conf-parser/core/resolv"
	"io/ioutil"
	"regexp"
	"strings"
)

func Load(path string) (*resolv.Config, error) {
	data, ferr := readConf(path)
	if ferr != nil {
		return nil, ferr
	}

	return load(data), nil
}

func load(data string) *resolv.Config {
	f := resolv.NewConf(nil)
	index := 0
	var lopen []resolv.Parser

	for {
		parseContext := func(reg *regexp.Regexp) bool {
			var context resolv.Parser

			if m := reg.FindStringIndex(data[index:]); m != nil {
				switch reg {
				case resolv.RegEventsHead:
					context = resolv.NewEvents()
				case resolv.RegHttpHead:
					context = resolv.NewHttp()
				case resolv.RegStreamHead:
					context = resolv.NewStream()
				case resolv.RegServerHead:
					context = resolv.NewServer()
				case resolv.RegLocationHead:
					context = resolv.NewLocation(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegIfHead:
					context = resolv.NewIf(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegUpstreamHead:
					context = resolv.NewUpstream(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegGeoHead:
					context = resolv.NewGeo(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegMapHead:
					context = resolv.NewMap(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegLimitExceptHead:
					context = resolv.NewLimitExcept(reg.FindStringSubmatch(data[index:])[1])
				case resolv.RegTypesHead:
					context = resolv.NewTypes()
				}
				if context != nil {
					lopen = append([]resolv.Parser{context}, lopen...)
					index += m[len(m)]
					return true
				}
			}
			return false
		}

		parseComment := func(reg *regexp.Regexp) bool {
			if ms := reg.FindStringSubmatch(data[index:]); len(ms) == 3 {
				m := reg.FindStringIndex(data[index:])
				c := resolv.NewComment(ms[2], !strings.Contains(ms[1], "\n"))
				if ct, ok := lopen[0].(resolv.Context); ok {
					ct.Add(c)
					lopen[0] = ct
				} else {
					f.Add(c)
				}
				index += m[len(m)] - 1
				return true
			} else {
				return false
			}
		}

		parseContextEnd := func(reg *regexp.Regexp) bool {
			if m := reg.FindStringIndex(data[index:]); m != nil {
				if c, ok := lopen[0].(resolv.Context); ok {
					lopen = lopen[1:]
					if ct, ok2 := lopen[0].(resolv.Context); ok2 {
						ct.Add(c)
						lopen[0] = ct
					} else {
						f.Add(c)
					}
				}
				index += m[len(m)]
				return true
			} else {
				return false
			}
		}

		parseKey := func(reg *regexp.Regexp) bool {
			var k *resolv.Key
			if m := reg.FindStringIndex(data[index:]); m != nil {
				ms := reg.FindStringSubmatch(data[index:])
				switch reg {
				case resolv.RegKey:
					k = resolv.NewKey(ms[1], "")
				case resolv.RegKeyValue:
					k = resolv.NewKey(ms[1], ms[2])
				}

				if ct, ok := lopen[0].(resolv.Context); ok {
					ct.Add(k)
					lopen[0] = ct
				} else {
					f.Add(k)
				}
				index += m[len(m)]
				return true
			}
			return false
		}

		switch {
		case parseContext(resolv.RegEventsHead),
			parseContext(resolv.RegHttpHead),
			parseContext(resolv.RegStreamHead),
			parseContext(resolv.RegServerHead),
			parseContext(resolv.RegLocationHead),
			parseContext(resolv.RegIfHead),
			parseContext(resolv.RegUpstreamHead),
			parseContext(resolv.RegGeoHead),
			parseContext(resolv.RegMapHead),
			parseContext(resolv.RegLimitExceptHead),
			parseContext(resolv.RegTypesHead),
			parseComment(resolv.RegCommentHead),
			parseContextEnd(resolv.RegContextEnd),
			parseKey(resolv.RegKeyValue),
			parseKey(resolv.RegKey):
			continue
		}
		break
	}

	return f
}

func readConf(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
