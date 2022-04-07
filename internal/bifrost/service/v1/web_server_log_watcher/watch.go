package web_server_log_watcher

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerLogWatcherService) Watch(
	ctx context.Context,
	request *v1.WebServerLogWatchRequest,
) (*v1.WebServerLog, error) {
	return w.store.WebServerLogWatcher().Watch(ctx, request)
}
