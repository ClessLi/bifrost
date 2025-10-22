package v1

import (
	"github.com/go-kit/kit/endpoint"
)

type WebServerConfigEndpoints interface {
	EndpointGetServerNames() endpoint.Endpoint
	EndpointGet() endpoint.Endpoint
	EndpointConnectivityCheckOfProxiedServers() endpoint.Endpoint
	EndpointUpdate() endpoint.Endpoint
}
