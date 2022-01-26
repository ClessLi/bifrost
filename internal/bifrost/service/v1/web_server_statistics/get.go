package web_server_statistics

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerStatisticsService) Get(ctx context.Context, servername *v1.ServerName) (*v1.Statistics, error) {
	return w.store.WebServerStatistics().Get(ctx, servername)
}
