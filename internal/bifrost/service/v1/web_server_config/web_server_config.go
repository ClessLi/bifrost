package web_server_config

import (
	svcv1 "github.com/ClessLi/bifrost/internal/bifrost/service/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type webServerConfigService struct {
	store storev1.StoreFactory
}

var _ svcv1.WebServerConfigService = &webServerConfigService{}

func NewWebServerConfigService(store storev1.StoreFactory) svcv1.WebServerConfigService {
	return &webServerConfigService{store: store}
}
