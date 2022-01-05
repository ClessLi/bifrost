package v1

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/web_server_config"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/transport/grpc"
)

type HandlersFactory interface {
	WebServerConfig() WebServerConfigHandler
}

type handlersFactory struct {
	eps epv1.EndpointsFactory
}

var _ HandlersFactory = &handlersFactory{}

func NewHandlersFactory(eps epv1.EndpointsFactory) HandlersFactory {
	return &handlersFactory{eps: eps}
}

func (h *handlersFactory) WebServerConfig() WebServerConfigHandler {
	return web_server_config.NewWebServerConfigHandler(h.eps)
}

func NewHandler(ep endpoint.Endpoint, decoder Decoder, encoder Encoder) grpc.Handler {
	return grpc.NewServer(ep, decoder.DecodeRequest, encoder.EncodeResponse)
}
