package v1

import "github.com/go-kit/kit/endpoint"

type WebServerBinCMDEndpoints interface {
	EndpointExec() endpoint.Endpoint
}
