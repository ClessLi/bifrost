package handler

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"
	"github.com/go-kit/kit/transport/grpc"
	"sync"
)

type WebServerConfigHandler interface {
	HandlerGet() grpc.Handler
	HandlerUpdate() grpc.Handler
}

var _ WebServerConfigHandler = &webServerConfigHandler{}

var (
	handlerGetOnce         = sync.Once{}
	handlerUpdateOnce      = sync.Once{}
	singletonHandlerGet    grpc.Handler
	singletonHandlerUpdate grpc.Handler
)

type webServerConfigHandler struct {
	eps     epv1.WebServerConfigEndpoints
	decoder decoder.Decoder
	encoder encoder.Encoder
}

func (w *webServerConfigHandler) HandlerGet() grpc.Handler {
	handlerGetOnce.Do(func() {
		if singletonHandlerGet == nil {
			singletonHandlerGet = NewHandler(w.eps.EndpointGet(), w.decoder, w.encoder)
		}
	})
	if singletonHandlerGet == nil {
		// logs.Fatel
		panic("web server config handler `Get` is nil")
	}
	return singletonHandlerGet
}

func (w *webServerConfigHandler) HandlerUpdate() grpc.Handler {
	handlerUpdateOnce.Do(func() {
		if singletonHandlerUpdate == nil {
			singletonHandlerUpdate = NewHandler(w.eps.EndpointUpdate(), w.decoder, w.encoder)
		}
	})
	if singletonHandlerUpdate == nil {
		// logs.Fatel
		panic("web server config handler `Update` is nil")
	}
	return singletonHandlerUpdate
}

func NewWebServerConfigHandler(eps epv1.EndpointsFactory) WebServerConfigHandler {
	return &webServerConfigHandler{
		eps:     eps.WebServerConfig(),
		decoder: decoder.NewWebServerConfigDecoder(),
		encoder: encoder.NewWebServerConfigEncoder(),
	}
}
