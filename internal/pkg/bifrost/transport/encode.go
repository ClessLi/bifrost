package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func EncodeOperateResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.OperateResponse)
	if resp.Error != nil {
		return &bifrostpb.OperateResponse{
			Ret: nil,
			Err: resp.Error.Error(),
		}, nil
	}
	return &bifrostpb.OperateResponse{
		Ret: resp.Result,
		Err: "",
	}, nil
}

func EncodeConfigResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.ConfigResponse)
	if resp.Error != nil {
		return &bifrostpb.ConfigResponse{
			Ret: nil,
			Err: resp.Error.Error(),
		}, nil
	}
	return &bifrostpb.ConfigResponse{
		Ret: &bifrostpb.Config{JData: resp.Result.JData},
		Err: "",
	}, nil
}

func EncodeHealthCheckResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.HealthResponse)
	status := grpc_health_v1.HealthCheckResponse_NOT_SERVING
	if resp.Status {
		status = grpc_health_v1.HealthCheckResponse_SERVING
	}
	return &grpc_health_v1.HealthCheckResponse{
		Status: status,
	}, nil
}
