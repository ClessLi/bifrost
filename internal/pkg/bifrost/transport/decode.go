package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func DecodeHealthCheckRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if _, ok := r.(*grpc_health_v1.HealthCheckRequest); ok {
		return endpoint.HealthRequest{}, nil
	}
	return nil, ErrInvalidHealthCheckRequest
}

func DecodeViewRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.ViewRequest); ok {
		return newViewRequestInfo(req), nil
	}
	return nil, ErrUnknownRequest
}

func newViewRequestInfo(req *bifrostpb.ViewRequest) *endpoint.ViewRequestInfo {
	return &endpoint.ViewRequestInfo{
		ViewType:   req.ViewType,
		ServerName: req.ServerName,
		Token:      req.Token,
	}
}

func DecodeUpdateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.UpdateRequest); ok {
		return newUpdateRequestInfo(req), nil
	}
	return nil, ErrUnknownRequest
}

func newUpdateRequestInfo(req *bifrostpb.UpdateRequest) *endpoint.UpdateRequestInfo {
	return &endpoint.UpdateRequestInfo{
		UpdateType: req.UpdateType,
		ServerName: req.ServerName,
		Token:      req.Token,
		Data:       req.Data,
	}
}

func DecodeWatchRequest(ctx context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.WatchRequest); ok {
		return newWatchRequestInfo(req), nil
	}
	return nil, ErrUnknownRequest
}

func newWatchRequestInfo(req *bifrostpb.WatchRequest) *endpoint.WatchRequestInfo {
	return &endpoint.WatchRequestInfo{
		WatchType:   req.WatchType,
		ServerName:  req.ServerName,
		Token:       req.Token,
		WatchObject: req.WatchObject,
	}
}
