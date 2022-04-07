package web_server_config

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerConfigService) Update(ctx context.Context, config *v1.WebServerConfig) error {
	return w.store.WebServerConfig().Update(ctx, config)
}
