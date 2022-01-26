package v1

import "github.com/go-kit/kit/endpoint"

type WebServerLogWatcherEndpoints interface {
	EndpointWatch() endpoint.Endpoint
}
