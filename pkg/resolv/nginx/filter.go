package nginx

import (
	"bytes"
	"regexp"
	"strconv"
)

func GetHTTP(ctx Context) *Http {
	if ps := ctx.QueryAllByKeywords(KeywordHTTP); len(ps) > 0 {
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
	if ps := ctx.QueryAllByKeywords(KeywordStream); len(ps) > 0 {
		if stream, ok := ps[0].(*Stream); ok {
			return stream
		}
	}
	return nil
}

//func GetServerNames(ctx Context) []Parser {
//	return ctx.QueryAllByKeywords(KeywordSvrName)
//}

func GetServerName(ctx Context) Parser {
	return ctx.QueryByKeywords(KeywordSvrName)
}

func GetPorts(ctx Context) []Parser {
	return ctx.QueryAllByKeywords(KeywordPort)
}

func GetPort(ctx Context) int {
	portValue := ctx.QueryByKeywords(KeywordPort).(*Key).Value
	if RegPortValue.MatchString(portValue) {
		portStr := RegPortValue.FindStringSubmatch(portValue)[1]
		port, err := strconv.Atoi(portStr)
		if err != nil {
			port = -1
		}
		return port
	}
	return -1
}

//func GetLocations(ctx Context) []Parser {
//	http := GetHTTP(ctx)
//	if http != nil {
//		return ctx.QueryAllByKeywords(KeywordLocations)
//	} else {
//		return []Parser{}
//	}
//}

//func AppendNewString(slice []string, elements ...string) []string {
//	for _, element := range elements {
//		element = StripSpace(element)
//		isDiff := true
//		for _, s := range slice {
//			if s == element {
//				isDiff = false
//				break
//			}
//		}
//		if isDiff {
//			slice = append(slice, element)
//		}
//	}
//	return slice
//}

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

func StripSpace(s string) string {
	spaceReg := regexp.MustCompile(`\s+`)
	return string(spaceReg.ReplaceAll(bytes.TrimSpace([]byte(s)), []byte(" ")))
}
