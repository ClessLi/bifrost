package service

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
)

type WebServerLogWatcherService interface {
	Watch(request *v1.WebServerLogWatchRequest) (<-chan []byte, context.CancelFunc, error)
}

type webServerLogWatcherService struct {
	eps epv1.WebServerLogWatcherEndpoints
}

func (w *webServerLogWatcherService) Watch(request *v1.WebServerLogWatchRequest) (<-chan []byte, context.CancelFunc, error) {
	reqCtx, cancel := context.WithCancel(GetContext())
	resp, err := w.eps.EndpointWatch()(reqCtx, request)
	if err != nil {
		cancel()
		return nil, cancel, err
	}
	return resp.(*v1.WebServerLog).Lines, cancel, nil
}

func newWebServerLogWatcherService(factory *factory) WebServerLogWatcherService {
	return &webServerLogWatcherService{eps: factory.eps.WebServerLogWatcher()}
}
