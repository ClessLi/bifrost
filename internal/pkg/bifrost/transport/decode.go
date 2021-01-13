package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

//func DecodeViewConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "ViewConfig")
//}
//
//func DecodeGetConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "GetConfig")
//}
//
//func DecodeUpdateConfigRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "UpdateConfig")
//}
//
//func DecodeViewStatisticsRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "ViewStatistics")
//}
//
//func DecodeStatusRequest(ctx context.Context, r interface{}) (interface{}, error) {
//	return decodeRequest(ctx, r, "Status")
//}
//
//func DecodeWatchLogRequest(ctx context.Context, r interface{}) (request interface{}, err error) {
//	//fmt.Println("decoding watch log")
//	return decodeRequest(ctx, r, "WatchLog")
//}

func DecodeHealthCheckRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if _, ok := r.(*grpc_health_v1.HealthCheckRequest); ok {
		return endpoint.HealthRequest{}, nil
	}
	return nil, ErrInvalidHealthCheckRequest
}

func DecodeViewRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.ViewRequest); ok {
		return newViewRequest(req), nil
	}
	return nil, ErrUnknownRequest
}

func newViewRequest(req *bifrostpb.ViewRequest) *endpoint.ViewRequest {
	return &endpoint.ViewRequest{
		ViewType:   req.ViewType,
		ServerName: req.ServerName,
		Token:      req.Token,
	}
}

func DecodeUpdateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.UpdateRequest); ok {
		return newUpdateRequest(req), nil
	}
	return nil, ErrUnknownRequest
}

func newUpdateRequest(req *bifrostpb.UpdateRequest) *endpoint.UpdateRequest {
	return &endpoint.UpdateRequest{
		UpdateType: req.UpdateType,
		ServerName: req.ServerName,
		Token:      req.Token,
		Data:       req.Data,
	}
}

func DecodeWatchRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.WatchRequest); ok {
		return newWatchRequest(req), nil
	}
	return nil, ErrUnknownRequest
}

func newWatchRequest(req *bifrostpb.WatchRequest) *endpoint.WatchRequest {
	return &endpoint.WatchRequest{
		WatchType:   req.WatchType,
		ServerName:  req.ServerName,
		Token:       req.Token,
		WatchObject: req.WatchObject,
	}
}
