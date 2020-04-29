package filter

import (
	"bytes"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"regexp"
)

func GetHTTP(context resolv.Context) *resolv.Http {
	return (context.Filter(KeywordHTTP)[0]).(*resolv.Http)
}

func GetStream(context resolv.Context) *resolv.Stream {
	return context.Filter(KeywordStream)[0].(*resolv.Stream)
}

//func GetServers(context resolv.Context) []*resolv.Server {
//	return context.Servers()
//}

func GetServerName(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordSvrName)
}

func GetPorts(context resolv.Context) []resolv.Parser {
	return context.Filter(KeywordPorts)
}

func GetLocations(context resolv.Context) []resolv.Parser {
	return context.Filter(keywordLocations)
}

func appendIfNotExist(list []string, elem string) []string {
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
