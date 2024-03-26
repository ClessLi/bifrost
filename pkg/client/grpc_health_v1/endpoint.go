package grpc_health_v1

import (
	"context"
	"errors"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"io"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type endpoints struct {
	healthCheckEndpoint endpoint.Endpoint
	healthWatchEndpoint endpoint.Endpoint
}

func (e *endpoints) Check(ctx context.Context, service string) (HealthStatus, error) {
	req := healthCheckRequest{Service: service}
	resp, err := e.healthCheckEndpoint(ctx, req)
	if err != nil {
		return UNKNOWN, err
	}
	if response, ok := resp.(healthCheckResponse); ok {
		return response.Status, nil
	}

	return UNKNOWN, errors.New("failed to check health, invalid response")
}

func (e *endpoints) Watch(ctx context.Context, service string) (<-chan HealthStatus, error) {
	req := healthCheckRequest{Service: service}
	resp, err := e.healthWatchEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	if c, ok := resp.(grpc_health_v1.Health_WatchClient); ok {
		statusChan := make(chan HealthStatus)
		go func(statusC chan<- HealthStatus) {
			for {
				select {
				case <-ctx.Done():
					logV1.Info("health watch connect closed by client")
					close(statusC)
					err := c.CloseSend()
					if err != nil {
						logV1.Warnf(err.Error())
					}

					return
				default:
					s, err := c.Recv()
					if err != nil {
						if errors.Is(err, io.EOF) {
							logV1.Info("health watch connect closed by server")

							return
						}
						logV1.Error(err.Error())

						return
					}
					statusC <- HealthStatus(s.GetStatus())
				}
			}
		}(statusChan)

		return statusChan, nil
	}

	return nil, errors.New("failed to watch health, invalid response")
}

func makeEndpoints(conn *grpc.ClientConn) *endpoints {
	return &endpoints{
		healthCheckEndpoint: grpctransport.NewClient(
			conn,
			"grpc.health.v1.Health", // "grpc.health.v1"
			"Check",
			encodeClientRequest,
			decodeClientResponse,
			new(grpc_health_v1.HealthCheckResponse),
		).Endpoint(),
		healthWatchEndpoint: func(ctx context.Context, request interface{}) (response interface{}, err error) {
			request, err = encodeClientRequest(ctx, request)
			if err != nil {
				return nil, err
			}
			client := grpc_health_v1.NewHealthClient(conn)

			return client.Watch(ctx, request.(*grpc_health_v1.HealthCheckRequest))
		},
	}
}
