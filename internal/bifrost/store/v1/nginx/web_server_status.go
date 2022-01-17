package nginx

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/monitor"
)

type webServerStatusStore struct {
	m                    monitor.Monitor
	webServersStatusFunc func() []*v1.WebServerInfo
	bifrostVersionFunc   func() string
}

func (w *webServerStatusStore) Get(ctx context.Context) (*v1.Metrics, error) {
	//TODO implement me
	panic("implement me")
}

func newWebServerStatusStore(store *webServerStore) *webServerStatusStore {
	//TODO implement me
	panic("implement me")
}
