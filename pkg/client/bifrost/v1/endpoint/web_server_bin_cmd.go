package endpoint

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"

	"github.com/go-kit/kit/endpoint"
)

type webServerBinCMDEndpoints struct {
	transport txpclient.WebServerBinCMDTransport
}

func (w *webServerBinCMDEndpoints) EndpointExec() endpoint.Endpoint {
	return w.transport.Exec().Endpoint()
}

func newWebServerBinCMDEndpoints(factory *factory) epv1.WebServerBinCMDEndpoints {
	return &webServerBinCMDEndpoints{transport: factory.transport.WebServerBinCMD()}
}
