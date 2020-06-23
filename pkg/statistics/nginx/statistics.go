package nginx

import (
	"github.com/ClessLi/bifrost/pkg/resolv/nginx"
	"strconv"
)

func HTTPServersNum(ctx nginx.Context) int {
	http := nginx.GetHTTP(ctx)
	if http == nil {
		return 0
	}
	return len(http.Servers())
}

func HTTPPorts(ctx nginx.Context) []int {
	var ports []int
	http := nginx.GetHTTP(ctx)
	if http == nil {
		return nil
	}
	for _, parser := range nginx.GetPorts(http) {
		port, err := strconv.Atoi(parser.(*nginx.Key).Value)
		if err != nil {
			continue
		}
		ports = nginx.SortInsertUniqInt(ports, port)
	}
	return ports
}

func HTTPServerNames(ctx nginx.Context) (serverNames []string) {
	for _, parser := range nginx.GetHTTPServers(ctx, nginx.ServerName) {
		if serverNameKey := nginx.GetServerName(parser.(*nginx.Server)); serverNameKey != nil {
			serverNames = nginx.AppendNewString(serverNames, serverNameKey.(*nginx.Key).Value)
		}
	}
	return
}

func HTTPLocationsNum(ctx nginx.Context) int {
	return len(nginx.GetLocations(ctx))
}

func StreamServersNum(ctx nginx.Context) int {
	stream := nginx.GetStream(ctx)
	if stream == nil {
		return 0
	}
	return len(stream.Servers())
}

func StreamPorts(ctx nginx.Context) []int {
	//var ports []string
	var ports []int
	stream := nginx.GetStream(ctx)
	if stream == nil {
		return nil
	}
	for _, parser := range nginx.GetPorts(stream) {
		//ports = appendNewString(ports, parser.(*resolv.Key).Value)
		port, err := strconv.Atoi(parser.(*nginx.Key).Value)
		if err != nil {
			continue
		}
		ports = nginx.SortInsertUniqInt(ports, port)
	}
	return ports
}
