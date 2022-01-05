package v1

import (
	"github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1/web_server_config"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type EndpointsFactory interface {
	WebServerConfig() WebServerConfigEndpoints
}

type endpoints struct {
	svc svcv1.ServiceFactory
}

func NewEndpoints(svc svcv1.ServiceFactory) EndpointsFactory {
	return EndpointsFactory(&endpoints{svc: svc})
}

func (e *endpoints) WebServerConfig() WebServerConfigEndpoints {
	return web_server_config.NewWebServerConfigEndpoints(e.svc)
}
