package endpoint

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/ClessLi/bifrost/pkg/client/bifrost/v1/transport"
)

type Factory interface {
	WebServerConfig() epv1.WebServerConfigEndpoints
	WebServerStatistics() epv1.WebServerStatisticsEndpoints
}

type factory struct {
	transport txpclient.Factory
}

func (f *factory) WebServerConfig() epv1.WebServerConfigEndpoints {
	return newWebServerConfigEndpoints(f)
}

func (f *factory) WebServerStatistics() epv1.WebServerStatisticsEndpoints {
	return newWebServerStatisticsEndpoints(f)
}

func New(transport txpclient.Factory) Factory {
	return &factory{transport: transport}
}
