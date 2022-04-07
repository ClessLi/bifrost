package grpc_health_v1

import (
	"context"

	"github.com/marmotedu/errors"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type healthCheckRequest struct {
	Service string `json:"service"`
}

type healthCheckResponse struct {
	Status HealthStatus `json:"status"`
}

func encodeClientRequest(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case healthCheckRequest:
		return &grpc_health_v1.HealthCheckRequest{Service: r.Service}, nil
	}

	return nil, errors.New("invalid request for gRPC health check")
}

func decodeClientResponse(ctx context.Context, r interface{}) (interface{}, error) {
	switch r := r.(type) {
	case *grpc_health_v1.HealthCheckResponse:
		return healthCheckResponse{Status: HealthStatus(r.GetStatus())}, nil
	}

	return nil, errors.New("invalid health check response from gRPC server")
}
