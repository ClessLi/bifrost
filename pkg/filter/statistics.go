package filter

import (
	"github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"
	"strconv"
)

func HTTPServers(ctx resolv.Context) int {
	return len(GetHTTP(ctx).Servers())
}

func HTTPPorts(ctx resolv.Context) []int {
	var ports []int
	for _, parser := range GetPorts(GetHTTP(ctx)) {
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = SortInsertUniqInt(ports, port)
	}
	return ports
}

//func HTTPPortsSTR(ctx resolv.Context) []string {
//	var ports []string
//	for _, parser := range GetPorts(GetHTTP(ctx)) {
//		ports = appendNewString(ports, parser.(*resolv.Key).Value)
//	}
//	return ports
//}

func HTTPServerNames(ctx resolv.Context) []string {
	var serverNames []string
	for _, parser := range GetServerName(GetHTTP(ctx)) {
		serverNames = appendNewString(serverNames, parser.(*resolv.Key).Value)
	}
	return serverNames
}

func HTTPLocations(ctx resolv.Context) int {
	return len(GetLocations(GetHTTP(ctx)))
}

func StreamServers(ctx resolv.Context) int {
	return len(GetStream(ctx).Servers())
}

func StreamPorts(ctx resolv.Context) []int {
	//var ports []string
	var ports []int
	for _, parser := range GetPorts(GetStream(ctx)) {
		//ports = appendNewString(ports, parser.(*resolv.Key).Value)
		port, err := strconv.Atoi(parser.(*resolv.Key).Value)
		if err != nil {
			continue
		}
		ports = SortInsertUniqInt(ports, port)
	}
	return ports
}
