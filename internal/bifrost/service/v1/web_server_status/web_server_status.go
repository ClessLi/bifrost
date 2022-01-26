package web_server_status

import storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"

type webServerStatusService struct {
	store storev1.StoreFactory
}

func NewWebServerStatusService(store storev1.StoreFactory) *webServerStatusService {
	return &webServerStatusService{store: store}
}
