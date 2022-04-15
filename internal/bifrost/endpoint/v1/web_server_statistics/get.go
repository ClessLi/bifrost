package web_server_statistics

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerStatisticsEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.ServerName); ok {
			return w.svc.WebServerStatistics().Get(ctx, req)
		}

		return nil, errors.Errorf("invalid get request, need *v1.ServerName, not %T", request)
	}
}
