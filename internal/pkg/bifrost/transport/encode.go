package transport

import (
	"github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb"
	"github.com/ClessLi/bifrost/internal/pkg/bifrost/endpoint"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

type watchResponder struct {
	bytesResponseChan chan *bifrostpb.BytesResponse
	signalChan        chan int
}

func (wr watchResponder) Respond() <-chan *bifrostpb.BytesResponse {
	return wr.bytesResponseChan
}

func (wr watchResponder) Close() error {
	select {
	case wr.signalChan <- 9:
		return nil
	case <-time.After(time.Second * 30):
		return ErrWatcherCloseTimeout
	}
}

func newWatchResponder(epWatchResponder endpoint.WatchResponder) *watchResponder {
	responder := &watchResponder{
		bytesResponseChan: make(chan *bifrostpb.BytesResponse),
		signalChan:        make(chan int),
	}
	go func() {
		for {
			select {
			case byteResponder := <-epWatchResponder.Respond():
				responder.bytesResponseChan <- &bifrostpb.BytesResponse{
					Ret: byteResponder.Respond(),
					Err: byteResponder.Error(),
				}
			case signal := <-responder.signalChan:
				if signal == 9 {
					epWatchResponder.Close()
					return
				}
			}
		}
	}()
	return responder
}

func EncodeViewResponse(_ context.Context, r interface{}) (interface{}, error) {
	if resp, ok := r.(endpoint.BytesResponder); ok {
		return &bifrostpb.BytesResponse{
			Ret: resp.Respond(),
			Err: resp.Error(),
		}, nil
	}
	return nil, ErrUnknownResponse
}

func EncodeUpdateResponse(_ context.Context, r interface{}) (interface{}, error) {
	if resp, ok := r.(endpoint.ErrorResponder); ok {
		return &bifrostpb.ErrorResponse{Err: resp.Error()}, nil
	}
	return nil, ErrUnknownResponse
}

func EncodeWatchResponse(_ context.Context, r interface{}) (interface{}, error) {
	if resp, ok := r.(endpoint.WatchResponder); ok {
		return newWatchResponder(resp), nil
	}
	return nil, ErrUnknownResponse
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
