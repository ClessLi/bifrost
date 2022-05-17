package endpoint

import (
	epv1 "github.com/yongPhone/bifrost/internal/bifrost/endpoint/v1"
	txpclient "github.com/yongPhone/bifrost/pkg/client/bifrost/v1/transport"
)

type Factory interface {
	WebServerConfig() epv1.WebServerConfigEndpoints
	WebServerStatistics() epv1.WebServerStatisticsEndpoints
	WebServerStatus() epv1.WebServerStatusEndpoints
	WebServerLogWatcher() epv1.WebServerLogWatcherEndpoints
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

func (f *factory) WebServerStatus() epv1.WebServerStatusEndpoints {
	return newWebServerStatusEndpoints(f)
}

func (f *factory) WebServerLogWatcher() epv1.WebServerLogWatcherEndpoints {
	return newWebServerLogWatcherEndpoints(f)
}

func New(transport txpclient.Factory) Factory {
	return &factory{transport: transport}
}
