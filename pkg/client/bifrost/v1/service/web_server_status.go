package service

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	pbv1 "github.com/ClessLi/bifrost/api/protobuf-spec/bifrostpb/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
)

type WebServerStatusService interface {
	Get() (*v1.Metrics, error)
}

type webServerStatusService struct {
	eps epv1.WebServerStatusEndpoints
}

func (w *webServerStatusService) Get() (*v1.Metrics, error) {
	resp, err := w.eps.EndpointGet()(GetContext(), new(pbv1.Null))
	if err != nil {
		return nil, err
	}

	return resp.(*v1.Metrics), nil
}

func newWebServerStatusService(factory *factory) WebServerStatusService {
	return &webServerStatusService{eps: factory.eps.WebServerStatus()}
}
