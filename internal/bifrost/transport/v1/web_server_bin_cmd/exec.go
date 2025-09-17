package web_server_bin_cmd

import (
	"context"

	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
)

func (w webServerBinCMDServer) Exec(ctx context.Context, request *pbv1.ExecuteRequest) (*pbv1.ExecuteResponse, error) {
	_, resp, err := w.handler.HandlerExec().ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	response := resp.(*pbv1.ExecuteResponse)

	return response, nil
}
