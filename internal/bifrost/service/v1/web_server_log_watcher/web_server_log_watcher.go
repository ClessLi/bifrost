package web_server_log_watcher

import storev1 "github.com/yongPhone/bifrost/internal/bifrost/store/v1"

type webServerLogWatcherService struct {
	store storev1.StoreFactory
}

func NewWebServerLogWatcherService(store storev1.StoreFactory) *webServerLogWatcherService {
	return &webServerLogWatcherService{store: store}
}
