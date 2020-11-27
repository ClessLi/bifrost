package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
)

func DecodeViewConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "ViewConfig")
}

func DecodeGetConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "GetConfig")
}

func DecodeUpdateConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "UpdateConfig")
}

func DecodeViewStatisticsRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "ViewStatistics")
}

func DecodeStatusRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return decodeRequest(ctx, r, "Status")
}

func DecodeWatchLogRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
	return r, nil
}

func decodeRequest(ctx context.Context, r interface{}, requestType string) (interface{}, error) {
	switch requestType {
	case "UpdateConfig":
		if req, ok := r.(*bifrostpb.ConfigRequest); ok {
			return endpoint.ConfigRequest{
				RequestType: requestType,
				Token:       req.Token,
				SvrName:     req.SvrName,
				Ret:         endpoint.Config{JData: req.Req.JData},
			}, nil
		}
		return nil, ErrInvalidConfigRequest
	default:
		if req, ok := r.(*bifrostpb.OperateRequest); ok {
			return endpoint.OperateRequest{
				RequestType: requestType,
				Token:       req.Token,
				SvrName:     req.SvrName,
				Param:       req.Param,
			}, nil
		}
		return nil, ErrInvalidOperateRequest
	}
}
