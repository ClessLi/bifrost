package handler

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/go-kit/kit/transport/grpc"
	"sync"
)

type WebServerStatisticsHandlers interface {
	HandlerGet() grpc.Handler
}

var _ WebServerStatisticsHandlers = &webServerStatisticsHandlers{}

type webServerStatisticsHandlers struct {
	onceGet             sync.Once
	singletonHandlerGet grpc.Handler
	eps                 epv1.WebServerStatisticsEndpoints
	decoder             decoder.Decoder
	encoder             encoder.Encoder
}

func (wss *webServerStatisticsHandlers) HandlerGet() grpc.Handler {
	wss.onceGet.Do(func() {
		if wss.singletonHandlerGet == nil {
			wss.singletonHandlerGet = NewHandler(wss.eps.EndpointGet(), wss.decoder, wss.encoder)
		}
	})

	if wss.singletonHandlerGet == nil {
		log.Fatal("web server statistics handler `Get` is nil")

		return nil
	}

	return wss.singletonHandlerGet
}

func NewWebServerStatisticsHandler(eps epv1.EndpointsFactory) WebServerStatisticsHandlers {
	return &webServerStatisticsHandlers{
		onceGet: sync.Once{},
		eps:     eps.WebServerStatistics(),
		decoder: decoder.NewWebServerStatisticsDecoder(),
		encoder: encoder.NewWebServerStatisticsEncoder(),
	}
}
