package handler

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"

	epv1 "github.com/yongPhone/bifrost/internal/bifrost/endpoint/v1"
	"github.com/yongPhone/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/yongPhone/bifrost/internal/bifrost/transport/v1/encoder"
)

type HandlersFactory interface {
	WebServerConfig() WebServerConfigHandlers
	WebServerStatistics() WebServerStatisticsHandlers
	WebServerStatus() WebServerStatusHandlers
	WebServerLogWatcher() WebServerLogWatcherHandlers
}

type handlersFactory struct {
	eps epv1.EndpointsFactory
}

var _ HandlersFactory = &handlersFactory{}

func NewHandlersFactory(eps epv1.EndpointsFactory) HandlersFactory {
	return &handlersFactory{eps: eps}
}

func (h *handlersFactory) WebServerConfig() WebServerConfigHandlers {
	return NewWebServerConfigHandler(h.eps)
}

func (h *handlersFactory) WebServerStatistics() WebServerStatisticsHandlers {
	return NewWebServerStatisticsHandler(h.eps)
}

func (h *handlersFactory) WebServerStatus() WebServerStatusHandlers {
	return NewWebServerStatusHandlers(h.eps)
}

func (h *handlersFactory) WebServerLogWatcher() WebServerLogWatcherHandlers {
	return NewWebServerLogWatcherHandlers(h.eps)
}

func NewHandler(ep endpoint.Endpoint, decoder decoder.Decoder, encoder encoder.Encoder) grpc.Handler {
	return grpc.NewServer(ep, decoder.DecodeRequest, encoder.EncodeResponse)
}
