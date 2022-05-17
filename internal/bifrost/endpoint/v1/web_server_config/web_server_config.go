package web_server_config

import (
	svcv1 "github.com/yongPhone/bifrost/internal/bifrost/service/v1"
)

type webServerConfigEndpoints struct {
	svc svcv1.ServiceFactory
}

func NewWebServerConfigEndpoints(svc svcv1.ServiceFactory) *webServerConfigEndpoints {
	return &webServerConfigEndpoints{svc: svc}
}
