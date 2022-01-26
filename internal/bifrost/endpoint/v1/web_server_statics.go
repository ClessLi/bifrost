package v1

import "github.com/go-kit/kit/endpoint"

type WebServerStatisticsEndpoints interface {
	EndpointGet() endpoint.Endpoint
}
