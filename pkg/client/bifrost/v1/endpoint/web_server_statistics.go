package endpoint

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"

	"github.com/go-kit/kit/endpoint"
)

type webServerStatisticsEndpoints struct {
	transport txpclient.WebServerStatisticsTransport
}

func (w *webServerStatisticsEndpoints) EndpointGet() endpoint.Endpoint {
	return w.transport.Get().Endpoint()
}

func newWebServerStatisticsEndpoints(factory *factory) epv1.WebServerStatisticsEndpoints {
	return &webServerStatisticsEndpoints{transport: factory.transport.WebServerStatistics()}
}
