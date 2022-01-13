package web_server_statistics

import (
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type webServerStatisticsEndpoints struct {
	svc svcv1.ServiceFactory
}

func NewWebServerStatisticsEndpoints(svc svcv1.ServiceFactory) *webServerStatisticsEndpoints {
	return &webServerStatisticsEndpoints{svc: svc}
}
