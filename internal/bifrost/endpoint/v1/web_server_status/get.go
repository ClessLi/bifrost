package web_server_status

import (
	"context"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/go-kit/kit/endpoint"
)

func (w *webServerStatusEndpoints) EndpointGet() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if _, ok := request.(*pbv1.Null); ok {
			return w.svc.WebServerStatus().Get(ctx)
		}

		return nil, code.Errorf("invalid get request, need *pbv1.Null, not %T", request)
	}
}
