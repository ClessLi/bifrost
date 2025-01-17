package web_server_bin_cmd

import svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"

type webServerBinCMDEndpoints struct {
	svc svcv1.ServiceFactory
}

func NewWebServerBinCMDEndpoints(svc svcv1.ServiceFactory) *webServerBinCMDEndpoints {
	return &webServerBinCMDEndpoints{svc: svc}
}
