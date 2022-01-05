package web_server_config

import (
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	v12 "github.com/ClessLi/bifrost/internal/bifrost/transport/v1"
	"github.com/go-kit/kit/transport/grpc"
	"sync"
)

var _ pbv1.WebServerConfigServer = &webServerConfigServer{}

type webServerConfigServer struct {
	handler v12.HandlersFactory
	Options *v12.Options
}

func NewWebServerConfigServer(handler v12.HandlersFactory, opts *v12.Options) pbv1.WebServerConfigServer {
	// TODO: default options item
	return &webServerConfigServer{
		handler: handler,
		Options: opts,
	}
}

var _ v12.WebServerConfigHandler = &webServerConfigHandler{}

var (
	handlerGetOnce         = sync.Once{}
	handlerUpdateOnce      = sync.Once{}
	singletonHandlerGet    grpc.Handler
	singletonHandlerUpdate grpc.Handler
)

type webServerConfigHandler struct {
	eps     epv1.EndpointsFactory
	decoder v12.Decoder
	encoder v12.Encoder
}

func (w *webServerConfigHandler) HandlerGet() grpc.Handler {
	handlerGetOnce.Do(func() {
		if singletonHandlerGet == nil {
			singletonHandlerGet = v12.NewHandler(w.eps.WebServerConfig().EndpointGet(), w.decoder, w.encoder)
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
			singletonHandlerUpdate = v12.NewHandler(w.eps.WebServerConfig().EndpointUpdate(), w.decoder, w.encoder)
		}
	})
	if singletonHandlerUpdate == nil {
		// logs.Fatel
		panic("web server config handler `Update` is nil")
	}
	return singletonHandlerUpdate
}

func NewWebServerConfigHandler(eps epv1.EndpointsFactory) v12.WebServerConfigHandler {
	return &webServerConfigHandler{
		eps:     eps,
		decoder: NewDecoder(),
		encoder: NewEncoder(),
	}
}
