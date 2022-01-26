package endpoint

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"
	"github.com/go-kit/kit/endpoint"
)

type webServerLogWatcherEndpoints struct {
	transport txpclient.WebServerLogWatcherTransport
}

func (w *webServerLogWatcherEndpoints) EndpointWatch() endpoint.Endpoint {
	return w.transport.Watch().Endpoint()
}

func newWebServerLogWatcherEndpoints(factory *factory) epv1.WebServerLogWatcherEndpoints {
	return &webServerLogWatcherEndpoints{transport: factory.transport.WebServerLogWatcher()}
}
