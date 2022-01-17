package v1

import (
	"github.com/ClessLi/bifrost/internal/bifrost/service/v1/web_server_config"
	"github.com/ClessLi/bifrost/internal/bifrost/service/v1/web_server_statistics"
	"github.com/ClessLi/bifrost/internal/bifrost/service/v1/web_server_status"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type ServiceFactory interface {
	WebServerConfig() WebServerConfigService
	WebServerStatistics() WebServerStatisticsService
	WebServerStatus() WebServerStatusService
}

var _ ServiceFactory = &serviceFactory{}

type serviceFactory struct {
	store storev1.StoreFactory
}

func (s *serviceFactory) WebServerConfig() WebServerConfigService {
	return web_server_config.NewWebServerConfigService(s.store)
}

func (s *serviceFactory) WebServerStatistics() WebServerStatisticsService {
	return web_server_statistics.NewWebServerStatisticsService(s.store)
}

func (s *serviceFactory) WebServerStatus() WebServerStatusService {
	return web_server_status.NewWebServerStatusService(s.store)
}

func NewServiceFactory(store storev1.StoreFactory) ServiceFactory {
	return &serviceFactory{store: store}
}
