package web_server_config

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerConfigService) Get(ctx context.Context, servername *v1.ServerName) (*v1.WebServerConfig, error) {
	return w.store.WebServerConfig().Get(ctx, servername)
}
