package web_server_bin_cmd

import (
	"context"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"

	"github.com/go-kit/kit/endpoint"
)

func (w *webServerBinCMDEndpoints) EndpointExec() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		if req, ok := request.(*v1.ExecuteRequest); ok {
			return w.svc.WebServerBinCMD().Exec(ctx, req)
		}

		return nil, code.Errorf("invalid get request, need *pbv1.ExecuteRequest, not %T", request)
	}
}
