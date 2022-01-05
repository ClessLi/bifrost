package web_server_config

import (
	"context"
	"github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"
)

func (w *webServerConfigEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.ServerName); ok {
			return w.svc.WebServerConfig().Get(ctx, req)
		}
		return nil, errors.Errorf("invalid get request, need *v1.ServerName, not %T", request)
	}
}
