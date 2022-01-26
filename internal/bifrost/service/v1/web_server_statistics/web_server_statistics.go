package web_server_statistics

import (
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
)

type webServerStatisticsService struct {
	store storev1.StoreFactory
}

func NewWebServerStatisticsService(store storev1.StoreFactory) *webServerStatisticsService {
	return &webServerStatisticsService{store: store}
}
