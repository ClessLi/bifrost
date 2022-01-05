package web_server_config

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type webServerConfigEndpoints struct {
	svc svcv1.ServiceFactory
}

var _ epv1.WebServerConfigEndpoints = &webServerConfigEndpoints{}

func NewWebServerConfigEndpoints(svc svcv1.ServiceFactory) epv1.WebServerConfigEndpoints {
	return &webServerConfigEndpoints{svc: svc}
}
