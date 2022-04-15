package web_server_config

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
)

func (w *webServerConfigEndpoints) EndpointUpdate() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.WebServerConfig); ok {
			err = w.svc.WebServerConfig().Update(ctx, req)
			if err != nil {
				return nil, err
			}

			return &v1.Response{Message: "update success"}, nil
		}

		return nil, errors.Errorf("invalid update request, need *v1.ServerConfig, not %T", request)
	}
}
