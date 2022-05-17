package web_server_status

import (
	"context"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
)

func (w *webServerStatusService) Get(ctx context.Context) (*v1.Metrics, error) {
	return w.store.WebServerStatus().Get(ctx)
}
