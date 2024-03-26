package handler //nolint:dupl

import (
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"sync"

	"github.com/go-kit/kit/transport/grpc"

	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"
)

type WebServerStatusHandlers interface {
	HandlerGet() grpc.Handler
}

var _ WebServerStatusHandlers = &webServerStatusHandlers{}

type webServerStatusHandlers struct {
	onceGet             sync.Once
	singletonHandlerGet grpc.Handler
	eps                 epv1.WebServerStatusEndpoints
	decoder             decoder.Decoder
	encoder             encoder.Encoder
}

func (w *webServerStatusHandlers) HandlerGet() grpc.Handler {
	w.onceGet.Do(func() {
		if w.singletonHandlerGet == nil {
			w.singletonHandlerGet = NewHandler(w.eps.EndpointGet(), w.decoder, w.encoder)
		}
	})
	if w.singletonHandlerGet == nil {
		logV1.Fatal("web server status handler `Get` is nil")

		return nil
	}

	return w.singletonHandlerGet
}

func NewWebServerStatusHandlers(eps epv1.EndpointsFactory) WebServerStatusHandlers {
	return &webServerStatusHandlers{
		onceGet: sync.Once{},
		eps:     eps.WebServerStatus(),
		decoder: decoder.NewWebServerStatusDecoder(),
		encoder: encoder.NewWebServerStatusEncoder(),
	}
}
