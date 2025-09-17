package web_server_config

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"

	"github.com/go-kit/kit/endpoint"
	"github.com/marmotedu/errors"
)

func (w *webServerConfigEndpoints) EndpointGetServerNames() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if _, ok := request.(*pbv1.Null); ok {
			return w.svc.WebServerConfig().GetServerNames(ctx)
		}

		return nil, errors.Errorf("invalid get request, need *pbv1.Null, not %T", request)
	}
}

func (w *webServerConfigEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.ServerName); ok {
			return w.svc.WebServerConfig().Get(ctx, req)
		}

		return nil, errors.Errorf("invalid get request, need *v1.ServerName, not %T", request)
	}
}
