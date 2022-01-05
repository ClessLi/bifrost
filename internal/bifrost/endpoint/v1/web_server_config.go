package v1

import (
	"github.com/go-kit/kit/endpoint"
)

type WebServerConfigEndpoints interface {
	EndpointGet() endpoint.Endpoint
	EndpointUpdate() endpoint.Endpoint
}
