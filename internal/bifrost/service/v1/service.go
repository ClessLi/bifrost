package v1

import (
	"github.com/ClessLi/bifrost/internal/bifrost/service/v1/web_server_config"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type ServiceFactory interface {
	WebServerConfig() WebServerConfigService
}

var _ ServiceFactory = &serviceFactory{}

type serviceFactory struct {
	store storev1.StoreFactory
}

func (s *serviceFactory) WebServerConfig() WebServerConfigService {
	return web_server_config.NewWebServerConfigService(s.store)
}

func NewServiceFactory(store storev1.StoreFactory) ServiceFactory {
	return &serviceFactory{store: store}
}
