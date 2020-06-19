package resolv

import (
	"bytes"
	"regexp"
)

func GetHTTP(ctx Context) *Http {
	if ps := ctx.QueryAll(KeywordHTTP); len(ps) > 0 {
		if http, ok := ps[0].(*Http); ok {
			return http
		}
	}
	return nil
}

func GetHTTPServers(ctx Context, orders ...Order) []Parser {
	http := GetHTTP(ctx)
	if http == nil {
		return []Parser{}
	}
	servers := http.Servers()
	if orders != nil {
		SortByOrders(&servers, orders...)
	}
	return servers
}

func GetStream(ctx Context) *Stream {
	if ps := ctx.QueryAll(KeywordStream); len(ps) > 0 {
		if stream, ok := ps[0].(*Stream); ok {
			return stream
		}
	}
	return nil
}

//func GetServerNames(ctx Context) []Parser {
//	return ctx.QueryAll(KeywordSvrName)
//}

func GetServerName(ctx Context) Parser {
	return ctx.Query(KeywordSvrName)
}

func GetPorts(ctx Context) []Parser {
	return ctx.QueryAll(KeywordPort)
}

func GetLocations(ctx Context) []Parser {
	http := GetHTTP(ctx)
	if http != nil {
		return ctx.QueryAll(KeywordLocations)
	} else {
		return []Parser{}
	}
}

func AppendNewString(slice []string, elem string) []string {
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
		slice[n-1] = num + 1

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

func stripSpace(s string) string {
	spaceReg := regexp.MustCompile(`\s{2,}`)
	return string(spaceReg.ReplaceAll(bytes.TrimSpace([]byte(s)), []byte(" ")))
}
