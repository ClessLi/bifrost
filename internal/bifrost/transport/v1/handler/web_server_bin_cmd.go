package handler //nolint:dupl

import (
	"sync"

	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/decoder"
	"github.com/ClessLi/bifrost/internal/bifrost/transport/v1/encoder"

	logV1 "github.com/ClessLi/component-base/pkg/log/v1"

	"github.com/go-kit/kit/transport/grpc"
)

type WebServerBinCMDHandlers interface {
	HandlerExec() grpc.Handler
}

var _ WebServerBinCMDHandlers = &webServerBinCMDHandlers{}

type webServerBinCMDHandlers struct {
	onceExec             sync.Once
	singletonHandlerExec grpc.Handler
	eps                  epv1.WebServerBinCMDEndpoints
	decoder              decoder.Decoder
	encoder              encoder.Encoder
}

func (w *webServerBinCMDHandlers) HandlerExec() grpc.Handler {
	w.onceExec.Do(func() {
		if w.singletonHandlerExec == nil {
			w.singletonHandlerExec = NewHandler(w.eps.EndpointExec(), w.decoder, w.encoder)
		}
	})
	if w.singletonHandlerExec == nil {
		logV1.Fatal("web server binary command handler `Exec` is nil")

		return nil
	}

	return w.singletonHandlerExec
}

func NewWebServerBinCMDHandlers(eps epv1.EndpointsFactory) WebServerBinCMDHandlers {
	return &webServerBinCMDHandlers{
		onceExec: sync.Once{},
		eps:      eps.WebServerBinCMD(),
		decoder:  decoder.NewWebServerBinCMDDecoder(),
		encoder:  encoder.NewWebServerBinCMDEncoder(),
	}
}
