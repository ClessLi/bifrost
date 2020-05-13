package statistics

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"strconv"
)

func HTTPServers(ctx resolv.Context) int {
	return len(resolv.GetHTTP(ctx).Servers())
}

func HTTPPorts(ctx resolv.Context) []int {
	var ports []int
	for _, parser := range resolv.GetPorts(resolv.GetHTTP(ctx)) {
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = resolv.SortInsertUniqInt(ports, port)
	}
	return ports
}

func HTTPServerNames(ctx resolv.Context) []string {
	var serverNames []string
	for _, parser := range resolv.GetServerNames(resolv.GetHTTP(ctx)) {
		serverNames = resolv.AppendNewString(serverNames, parser.(*resolv.Key).Value)
	}
	return serverNames
}

func HTTPLocations(ctx resolv.Context) int {
	return len(resolv.GetLocations(resolv.GetHTTP(ctx)))
}

func StreamServers(ctx resolv.Context) int {
	return len(resolv.GetStream(ctx).Servers())
}

func StreamPorts(ctx resolv.Context) []int {
	//var ports []string
	var ports []int
	for _, parser := range resolv.GetPorts(resolv.GetStream(ctx)) {
		//ports = appendNewString(ports, parser.(*resolv.Key).Value)
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = resolv.SortInsertUniqInt(ports, port)
	}
	return ports
}
