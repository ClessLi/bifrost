package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func EncodeBifrostServiceResponse(_ context.Context, r interface{}) (interface{}, error) {
	switch r.(type) {
	case endpoint.Response:
		errStr := ""
		if r.(endpoint.Response).Error != nil {
			errStr = r.(endpoint.Response).Error.Error()
		}
		return &bifrostpb.Response{
			Ret: r.(endpoint.Response).Result,
			Err: errStr,
		}, nil
	case *endpoint.Watcher:
		return r, nil
	default:
		return nil, ErrUnknownResponse
	}
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
