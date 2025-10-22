package web_server_config

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerConfigService) GetServerNames(ctx context.Context) (*v1.ServerNames, error) {
	return w.store.WebServerConfig().GetServerNames(ctx)
}

func (w *webServerConfigService) Get(ctx context.Context, servername *v1.ServerName) (*v1.WebServerConfig, error) {
	return w.store.WebServerConfig().Get(ctx, servername)
}

func (w *webServerConfigService) ConnectivityCheckOfProxiedServers(ctx context.Context, pos *v1.WebServerConfigContextPos) (*v1.ContextData, error) {
	return w.store.WebServerConfig().ConnectivityCheckOfProxiedServers(ctx, pos)
}
