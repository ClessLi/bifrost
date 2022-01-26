package web_server_log_watcher

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"
)

func (w *webServerLogWatcherEndpoints) EndpointWatch() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.WebServerLogWatchRequest); ok {
			return w.svc.WebServerLogWatcher().Watch(ctx, req)
		}
		return nil, errors.Errorf("invalid get request, need *v1.WebServerLogWatcherRequest, not %T", request)
	}
}
