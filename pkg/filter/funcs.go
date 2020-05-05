package filter

import (
	"bytes"
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"regexp"
	"strconv"
)

func GetHTTP(ctx resolv.Context) *resolv.Http {
	return (ctx.Filter(KeywordHTTP)[0]).(*resolv.Http)
}

func GetHTTPServers(ctx resolv.Context, tagger func([]*resolv.Server) []int) []*resolv.Server {
	servers := GetHTTP(ctx).Servers()
	if tagger != nil {
		tags := tagger(servers)
		servers = ServersInsertionSort(tags, servers)
	}
	return servers
}

func ServersInsertionSort(tags []int, servers []*resolv.Server) []*resolv.Server {
	n := len(tags)
	m := len(servers)
	if n <= 1 {
		return servers
	} else if n != m {
		return servers
	}
	for i := 1; i < n; i++ {
		tag := tags[i]
		server := servers[i]
		j := i - 1
		for ; j >= 0; j-- {
			if tags[j] > tag {
				tags[j+1] = tags[j]
				servers[j+1] = servers[j]
			} else {
				break
			}
		}
		tags[j+1] = tag
		servers[j+1] = server
	}
	return servers
}

func GetStream(ctx resolv.Context) *resolv.Stream {
	return ctx.Filter(KeywordStream)[0].(*resolv.Stream)
}

//func GetServers(context resolv.Context) []*resolv.Server {
//	return context.Servers()
//}

func GetServerName(ctx resolv.Context) []resolv.Parser {
	return ctx.Filter(KeywordSvrName)
}

func GetPorts(ctx resolv.Context) []resolv.Parser {
	return ctx.Filter(KeywordPorts)
}

func GetLocations(ctx resolv.Context) []resolv.Parser {
	return ctx.Filter(keywordLocations)
}

func appendNewString(slice []string, elem string) []string {
	elem = stripSpace(elem)
	var tmp []string
	for _, s := range slice {
		if s == elem {
			return slice
		}
		tmp = append(tmp, s)
	}
	tmp = append(tmp, elem)
	return tmp
}

func SortInsertInt(slice []int, ints ...int) []int {
	n := len(slice)
	for _, num := range ints {
		slice = append(slice, num+1)

		i := 0
		for slice[i] <= num {
			i++
		}

		for j := n; i < j; j-- {
			slice[j] = slice[j-1]
		}

		slice[i] = num
		n++

	}

	return slice
}

func SortInsertUniqInt(slice []int, ints ...int) []int {
	n := len(slice)
	for _, num := range ints {
		if n <= 0 {
			slice = append(slice, num)
			n++
			continue
		}

		if slice[n-1] == num {
			continue
		} else if slice[n-1] < num {
			slice = append(slice, num)
			n++
			continue
		}

		tmp := slice[n-1]
		slice[n-1] = num

		i := 0
		for slice[i] < num {
			i++
		}

		if slice[i] == num {
			slice[n-1] = tmp
			continue
		}

		for j := n - 1; i < j; j-- {
			slice[j] = slice[j-1]
		}

		slice[i] = num
		slice = append(slice, tmp)
		n++

	}

	return slice
}

func ServersTaggerByPort(servers []*resolv.Server) []int {
	tags := make([]int, 0, 1)
	for _, server := range servers {
		tag, err := strconv.Atoi(stripSpace(GetPorts(server)[0].(*resolv.Key).Value))
		if err != nil {
			return nil
		}

		//tags = SortInsertInt(tags, tag)
		tags = append(tags, tag)
	}
	return tags
}

func stripSpace(s string) string {
	spaceReg := regexp.MustCompile(`\s{2,}`)
	return string(spaceReg.ReplaceAll(bytes.TrimSpace([]byte(s)), []byte(" ")))
}
