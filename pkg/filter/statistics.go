package filter

import "github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"

func HTTPServers(ctx resolv.Context) int {
	return len(GetHTTP(ctx).Servers())
}

func HTTPPorts(ctx resolv.Context) []string {
	var ports []string
	for _, parser := range GetPorts(GetHTTP(ctx)) {
		ports = appendNewString(ports, parser.(*resolv.Key).Value)
	}
	return ports
}

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

func StreamPorts(ctx resolv.Context) []string {
	var ports []string
	for _, parser := range GetPorts(GetStream(ctx)) {
		ports = appendNewString(ports, parser.(*resolv.Key).Value)
	}
	return ports
}
