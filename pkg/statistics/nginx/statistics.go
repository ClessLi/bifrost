package nginx

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"strconv"
)

//func HTTPServersNum(ctx nginx.Context) int {
//	http := nginx.GetHTTP(ctx)
//	if http == nil {
//		return 0
//	}
//	return len(http.Servers())
//}

func HTTPPorts(ctx nginx.Context) []int {
	var ports []int
	http := nginx.GetHTTP(ctx)
	if http == nil {
		return nil
	}
	for _, parser := range nginx.GetPorts(http) {
		portValue := parser.(*nginx.Key).Value
		if nginx.RegPortValue.MatchString(portValue) {
			portStr := nginx.RegPortValue.FindStringSubmatch(portValue)[1]
			port, err := strconv.Atoi(portStr)
			if err != nil {
				continue
			}
			ports = nginx.SortInsertUniqInt(ports, port)
		}
	}
	return ports
}

//func HTTPServerNames(ctx nginx.Context) (serverNames []string) {
//	for _, parser := range nginx.GetHTTPServers(ctx, nginx.ServerName) {
//		if serverNameKey := nginx.GetServerName(parser.(*nginx.Server)); serverNameKey != nil {
//			serverNames = nginx.AppendNewString(serverNames, serverNameKey.(*nginx.Key).Value)
//		}
//	}
//	return
//}

func HTTPServers(ctx nginx.Context) (int, map[string][]int) {
	// 初始化serversInfo
	serversInfo := make(map[string][]int, 0)

	// 获取按server name排序后的servers
	servers := nginx.GetHTTPServers(ctx, nginx.ServerName, nginx.ServerPort)
	// 获取servers总数
	l := len(servers)

	// 生成serversInfo
	for _, parser := range servers {
		var tmpServerName string
		if serverNameKey := nginx.GetServerName(parser.(*nginx.Server)); serverNameKey != nil {
			tmpServerName = nginx.StripSpace(serverNameKey.(*nginx.Key).Value)
		} else {
			tmpServerName = ""
		}

		port := nginx.GetPort(parser.(*nginx.Server))
		if port < 0 {
			continue
		}

		if serversInfo[tmpServerName] == nil {
			serversInfo[tmpServerName] = make([]int, 0)
		}
		serversInfo[tmpServerName] = append(serversInfo[tmpServerName], port)
	}

	return l, serversInfo
}

//func HTTPLocationsNum(ctx nginx.Context) int {
//	return len(nginx.GetLocations(ctx))
//}
//
//func StreamServersNum(ctx nginx.Context) int {
//	stream := nginx.GetStream(ctx)
//	if stream == nil {
//		return 0
//	}
//	return len(stream.Servers())
//}

func StreamServers(ctx nginx.Context) (int, []int) {
	//var ports []string
	var ports []int
	stream := nginx.GetStream(ctx)
	if stream == nil {
		return 0, nil
	}
	for _, parser := range nginx.GetPorts(stream) {
		//ports = appendNewString(ports, parser.(*resolv.Key).Value)
		port, err := strconv.Atoi(parser.(*nginx.Key).Value)
		if err != nil {
			continue
		}
		ports = nginx.SortInsertUniqInt(ports, port)
	}
	return len(ports), ports
}
