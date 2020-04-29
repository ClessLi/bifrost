package filter

import (
	"bytes"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"regexp"
)

func GetHTTP(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordHTTP)
}

func GetStream(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordStream)
}

func GetServers(context resolv.Context) []*resolv.Server {
	return context.Servers()
}

func GetServerName(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordSvrName)
}

func GetPorts(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordPorts)
}

func GetLocations(context resolv.Context) []resolv.Parser {
	return context.Filter(keywordLocations)
}

func appendIsNotExist(list []string, elem string) []string {
	elem = stripSpace(elem)
	var tmp []string
	for _, s := range list {
		if s == elem {
			return list
		}
		tmp = append(tmp, s)
	}
	tmp = append(tmp, elem)
	return tmp
}

func stripSpace(s string) string {
	spaceReg := regexp.MustCompile(`\s{2,}`)
	return string(spaceReg.ReplaceAll(bytes.TrimSpace([]byte(s)), []byte(" ")))
}
