package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func DecodeHealthCheckRequest(_ context.Context, r interface{}) (interface{}, error) {
	if _, ok := r.(*grpc_health_v1.HealthCheckRequest); ok {
		return endpoint.HealthRequestInfo{}, nil
	}
	return nil, ErrInvalidHealthCheckRequest
}

func DecodeViewRequest(_ context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.ViewRequest); ok {
		return &endpoint.ViewRequestInfo{
			ViewType:   req.ViewType,
			ServerName: req.ServerName,
			Token:      req.Token,
		}, nil
	}
	return nil, ErrUnknownRequest
}

func DecodeUpdateRequest(_ context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.UpdateRequest); ok {
		return &endpoint.UpdateRequestInfo{
			UpdateType: req.UpdateType,
			ServerName: req.ServerName,
			Token:      req.Token,
			Data:       req.Data,
		}, nil
	}
	return nil, ErrUnknownRequest
}

func DecodeWatchRequest(_ context.Context, r interface{}) (interface{}, error) {
	if req, ok := r.(*bifrostpb.WatchRequest); ok {
		return &endpoint.WatchRequestInfo{
			WatchType:   req.WatchType,
			ServerName:  req.ServerName,
			Token:       req.Token,
			WatchObject: req.WatchObject,
		}, nil
	}
	return nil, ErrUnknownRequest
}
