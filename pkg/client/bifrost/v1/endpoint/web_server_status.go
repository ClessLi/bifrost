package endpoint

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"

	"github.com/go-kit/kit/endpoint"
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
