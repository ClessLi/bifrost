package v1

import "github.com/go-kit/kit/endpoint"

type WebServerStatusEndpoints interface {
	EndpointGet() endpoint.Endpoint
}
