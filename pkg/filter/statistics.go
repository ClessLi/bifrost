package filter

import "github.com/ClessLi/go-nginx-conf-parser/pkg/resolv"

func HTTPServers(ctx resolv.Context) int {
	svrs := 0
	https := GetHTTP(ctx)
	for _, http := range https {
		svrs += len(GetServers(http.(resolv.Context)))
	}
	return svrs
}

func HTTPPorts(ctx resolv.Context) []string {
	var ports []string
	https := GetHTTP(ctx)
	for _, http := range https {
		keys := GetPorts(http.(resolv.Context))
		for _, key := range keys {
			ports = appendIsNotExist(ports, key.(*resolv.Key).Value)
		}
	}
	return ports
}

func HTTPServerNames(ctx resolv.Context) []string {
	var serverNames []string
	https := GetHTTP(ctx)
	for _, http := range https {
		keys := GetServerName(http.(resolv.Context))
		for _, key := range keys {
			serverNames = appendIsNotExist(serverNames, key.(*resolv.Key).Value)
		}
	}
	return serverNames
}

func HTTPLocations(ctx resolv.Context) int {
	locations := 0
	https := GetHTTP(ctx)
	for _, http := range https {
		locations += len(GetLocations(http.(resolv.Context)))
	}
	return locations
}

func StreamServers(ctx resolv.Context) int {
	svrs := 0
	streams := GetStream(ctx)
	for _, stream := range streams {
		svrs += len(GetServers(stream.(resolv.Context)))
	}
	return svrs
}

func StreamPorts(ctx resolv.Context) []string {
	var ports []string
	streams := GetStream(ctx)
	for _, stream := range streams {
		keys := GetPorts(stream.(resolv.Context))
		for _, key := range keys {
			ports = appendIsNotExist(ports, key.(*resolv.Key).Value)
		}
	}
	return ports
}
