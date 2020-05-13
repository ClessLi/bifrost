package resolv

import (
	"bytes"
	"regexp"
)

func GetHTTP(ctx Context) *Http {
	return (ctx.QueryAll(KeywordHTTP)[0]).(*Http)
}

//func GetHTTPServers(ctx resolv.Context, tagger func([]*resolv.Server) []int) []*resolv.Server {
//	servers := GetHTTP(ctx).Servers()
//	if tagger != nil {
//		tags := tagger(servers)
//		servers = ServersInsertionSort(tags, servers)
//	}
//	return servers
//}

func GetHTTPServers(ctx Context, orders ...Order) []Parser {
	servers := GetHTTP(ctx).Servers()
	if orders != nil {
		//	ServersInsertionSort(&servers, orders...)
		SortByOrders(&servers, orders...)
	}
	//ServersInsertionSort(&servers, orders...)
	//SortByOrders(&servers, orders...)
	return servers
}

//func ServersInsertionSort(slice *[]Parser, orders ...func(Parser) int) {
//	n := len(*slice)
//	for _, order := range orders {
//
//		if order == nil {
//			break
//		}
//		cache := map[Parser]int{}
//		//cache := map[*resolv.Server]int{}
//		if n <= 1 {
//			return
//		}
//
//		for i := 1; i < n; i++ {
//			tmp := (*slice)[i]
//			d, tmpOK := cache[tmp]
//			if !tmpOK {
//				d = order(tmp)
//				cache[tmp] = d
//			}
//			j := i - 1
//			for ; j >= 0; j-- {
//				c, ok := cache[(*slice)[j]]
//				if !ok {
//					c = order((*slice)[j])
//					cache[(*slice)[j]] = c
//				}
//
//				if c > d {
//					(*slice)[j+1] = (*slice)[j]
//				} else {
//					break
//				}
//
//			}
//			(*slice)[j+1] = tmp
//		}
//	}
//	return
//}

//func OrderByPort(parser Parser) int {
//	if ctx, ok := parser.(Context); ok {
//		weight, err := strconv.Atoi(stripSpace(GetPorts(ctx)[0].(*Key).Value))
//		if err != nil {
//			weight = 0
//		}
//		return weight
//	} else {
//		return 0
//	}
//}

//func ServerName(server *resolv.Server) int {
//	serverName := GetServerNames(server)
//	if serverName == nil {
//		return 0
//	}
//	sn := stripSpace(serverName[0].(*resolv.Key).Value)
//	//var weightSTR string
//	//for _, b := range sn {
//	//	weightSTR = fmt.Sprintf("%s%d", weightSTR, b)
//	//}
//	//bs := []byte(sn)
//	//n := len(bs)
//	//weight := 0
//	//for i := n; i > 0; i-- {
//	//	m := int(bs[n-i])
//	//	weight += m * int(math.Pow(float64(1000), float64(i)))
//	//}
//	//weight, err := strconv.ParseInt(weightSTR, 10, 64)
//	//weight64, err := base64.RawURLEncoding.DecodeString(sn)
//	weight64, _ := base64.RawURLEncoding.DecodeString(sn)
//	//if err != nil {
//	//	return 0
//	//}
//	weightBig := new(big.Int)
//	weightBig.SetBytes(weight64)
//	weight := int(weightBig.Int64())
//	return weight
//}

func GetStream(ctx Context) *Stream {
	return ctx.QueryAll(KeywordStream)[0].(*Stream)
}

func GetServerNames(ctx Context) []Parser {
	return ctx.QueryAll(KeywordSvrName)
}

func GetServerName(ctx Context) Parser {
	return ctx.Query(KeywordSvrName)
}

func GetPorts(ctx Context) []Parser {
	return ctx.QueryAll(KeywordPort)
}

func GetLocations(ctx Context) []Parser {
	return ctx.QueryAll(KeywordLocations)
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

func stripSpace(s string) string {
	spaceReg := regexp.MustCompile(`\s{2,}`)
	return string(spaceReg.ReplaceAll(bytes.TrimSpace([]byte(s)), []byte(" ")))
}
