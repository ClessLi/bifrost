package web_server_statistics

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/go-kit/kit/endpoint"
)

func (w *webServerStatisticsEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.ServerName); ok {
			return w.svc.WebServerStatistics().Get(ctx, req)
		}

		return nil, code.Errorf("invalid get request, need *v1.ServerName, not %T", request)
	}
}
