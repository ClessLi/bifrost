package web_server_config

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type webServerConfigService struct {
	store storev1.StoreFactory
}

func NewWebServerConfigService(store storev1.StoreFactory) *webServerConfigService {
	return &webServerConfigService{store: store}
}
