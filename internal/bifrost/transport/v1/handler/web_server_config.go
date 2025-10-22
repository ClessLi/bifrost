package handler

import (
	"sync"

	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/go-kit/kit/transport/grpc"
)

type WebServerConfigHandlers interface {
	HandlerGetServerNames() grpc.Handler
	HandlerGet() grpc.Handler
	HandlerConnectivityCheckOfProxiedServers() grpc.Handler
	HandlerUpdate() grpc.Handler
}

var _ WebServerConfigHandlers = &webServerConfigHandlers{}

type webServerConfigHandlers struct {
	onceGetServerNames                                sync.Once
	onceGet                                           sync.Once
	onceCheckProxiedServersConnectivity               sync.Once
	onceUpdate                                        sync.Once
	singletonHandlerGetServerNames                    grpc.Handler
	singletonHandlerGet                               grpc.Handler
	singletonHandlerConnectivityCheckOfProxiedServers grpc.Handler
	singletonHandlerUpdate                            grpc.Handler
	eps                                               epv1.WebServerConfigEndpoints
	decoder                                           decoder.Decoder
	encoder                                           encoder.Encoder
}

func (wsc *webServerConfigHandlers) HandlerGetServerNames() grpc.Handler {
	wsc.onceGetServerNames.Do(func() {
		if wsc.singletonHandlerGetServerNames == nil {
			wsc.singletonHandlerGetServerNames = NewHandler(wsc.eps.EndpointGetServerNames(), wsc.decoder, wsc.encoder)
		}
	})
	if wsc.singletonHandlerGetServerNames == nil {
		logV1.Fatal("web server config handler `GetServerNames` is nil")

		return nil
	}

	return wsc.singletonHandlerGetServerNames
}

func (wsc *webServerConfigHandlers) HandlerGet() grpc.Handler {
	wsc.onceGet.Do(func() {
		if wsc.singletonHandlerGet == nil {
			wsc.singletonHandlerGet = NewHandler(wsc.eps.EndpointGet(), wsc.decoder, wsc.encoder)
		}
	})
	if wsc.singletonHandlerGet == nil {
		logV1.Fatal("web server config handler `Get` is nil")

		return nil
	}

	return wsc.singletonHandlerGet
}

func (wsc *webServerConfigHandlers) HandlerConnectivityCheckOfProxiedServers() grpc.Handler {
	wsc.onceCheckProxiedServersConnectivity.Do(func() {
		if wsc.singletonHandlerConnectivityCheckOfProxiedServers == nil {
			wsc.singletonHandlerConnectivityCheckOfProxiedServers = NewHandler(wsc.eps.EndpointConnectivityCheckOfProxiedServers(), wsc.decoder, wsc.encoder)
		}
	})
	if wsc.singletonHandlerConnectivityCheckOfProxiedServers == nil {
		logV1.Fatal("web server config handler `ConnectivityCheckOfProxiedServers` is nil")

		return nil
	}

	return wsc.singletonHandlerConnectivityCheckOfProxiedServers
}

func (wsc *webServerConfigHandlers) HandlerUpdate() grpc.Handler {
	wsc.onceUpdate.Do(func() {
		if wsc.singletonHandlerUpdate == nil {
			wsc.singletonHandlerUpdate = NewHandler(wsc.eps.EndpointUpdate(), wsc.decoder, wsc.encoder)
		}
	})
	if wsc.singletonHandlerUpdate == nil {
		logV1.Fatal("web server config handler `Update` is nil")

		return nil
	}

	return wsc.singletonHandlerUpdate
}

func NewWebServerConfigHandler(eps epv1.EndpointsFactory) WebServerConfigHandlers {
	return &webServerConfigHandlers{
		onceGetServerNames: sync.Once{},
		onceGet:            sync.Once{},
		onceUpdate:         sync.Once{},
		eps:                eps.WebServerConfig(),
		decoder:            decoder.NewWebServerConfigDecoder(),
		encoder:            encoder.NewWebServerConfigEncoder(),
	}
}
