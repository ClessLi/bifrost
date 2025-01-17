package v1

import (
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_bin_cmd"
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_config"
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_log_watcher"
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_statistics"
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_status"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type EndpointsFactory interface {
	WebServerConfig() WebServerConfigEndpoints
	WebServerStatistics() WebServerStatisticsEndpoints
	WebServerStatus() WebServerStatusEndpoints
	WebServerLogWatcher() WebServerLogWatcherEndpoints
	WebServerBinCMD() WebServerBinCMDEndpoints
}

var _ EndpointsFactory = &endpoints{}

type endpoints struct {
	svc svcv1.ServiceFactory
}

func NewEndpoints(svc svcv1.ServiceFactory) EndpointsFactory {
	return EndpointsFactory(&endpoints{svc: svc})
}

func (e *endpoints) WebServerConfig() WebServerConfigEndpoints {
	return web_server_config.NewWebServerConfigEndpoints(e.svc)
}

func (e *endpoints) WebServerStatistics() WebServerStatisticsEndpoints {
	return web_server_statistics.NewWebServerStatisticsEndpoints(e.svc)
}

func (e *endpoints) WebServerStatus() WebServerStatusEndpoints {
	return web_server_status.NewWebServerStatusEndpoints(e.svc)
}

func (e *endpoints) WebServerLogWatcher() WebServerLogWatcherEndpoints {
	return web_server_log_watcher.NewWebServerLogWatcherEndpoints(e.svc)
}

func (e *endpoints) WebServerBinCMD() WebServerBinCMDEndpoints {
	return web_server_bin_cmd.NewWebServerBinCMDEndpoints(e.svc)
}
