package endpoint

import (
	"github.com/go-kit/kit/endpoint"

	epv1 "github.com/yongPhone/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/yongPhone/bifrost/pkg/client/bifrost/v1/transport"
)

type webServerStatusEndpoints struct {
	transport txpclient.WebServerStatusTransport
}

func (w *webServerStatusEndpoints) EndpointGet() endpoint.Endpoint {
	return w.transport.Get().Endpoint()
}

func newWebServerStatusEndpoints(factory *factory) epv1.WebServerStatusEndpoints {
	return &webServerStatusEndpoints{transport: factory.transport.WebServerStatus()}
}
