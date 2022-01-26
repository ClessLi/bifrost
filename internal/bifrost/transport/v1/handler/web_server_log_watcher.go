package handler

import (
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/go-kit/kit/transport/grpc"
	"sync"
)

type WebServerLogWatcherHandlers interface {
	HandlerWatch() grpc.Handler
}

var _ WebServerLogWatcherHandlers = &webServerLogWatcherHandlers{}

type webServerLogWatcherHandlers struct {
	onceWatch             sync.Once
	singletonHandlerWatch grpc.Handler
	eps                   epv1.WebServerLogWatcherEndpoints
	decoder               decoder.Decoder
	encoder               encoder.Encoder
}

func (lw *webServerLogWatcherHandlers) HandlerWatch() grpc.Handler {
	lw.onceWatch.Do(func() {
		if lw.singletonHandlerWatch == nil {
			lw.singletonHandlerWatch = NewHandler(lw.eps.EndpointWatch(), lw.decoder, lw.encoder)
		}
	})

	if lw.singletonHandlerWatch == nil {
		log.Fatal("web server log watcher handler `Watch` is nil")

		return nil
	}

	return lw.singletonHandlerWatch

}

func NewWebServerLogWatcherHandlers(eps epv1.EndpointsFactory) WebServerLogWatcherHandlers {
	return &webServerLogWatcherHandlers{
		onceWatch:             sync.Once{},
		singletonHandlerWatch: nil,
		eps:                   eps.WebServerLogWatcher(),
		decoder:               decoder.NewWebServerLogWatcherDecoder(),
		encoder:               encoder.NewWebServerLogWatcherEncoder(),
	}
}
