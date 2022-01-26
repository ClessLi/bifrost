package web_server_log_watcher

import (
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
)

type webServerLogWatcherEndpoints struct {
	svc svcv1.ServiceFactory
}

func NewWebServerLogWatcherEndpoints(svc svcv1.ServiceFactory) *webServerLogWatcherEndpoints {
	return &webServerLogWatcherEndpoints{svc: svc}
}
