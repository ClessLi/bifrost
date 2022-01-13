package service

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
)

type WebServerStatisticsService interface {
	Get(servername *v1.ServerName) (*v1.Statistics, error)
}

type webServerStatisticsService struct {
	eps epv1.WebServerStatisticsEndpoints
}

func (w *webServerStatisticsService) Get(servername *v1.ServerName) (*v1.Statistics, error) {
	resp, err := w.eps.EndpointGet()(GetContext(), servername)
	if err != nil {
		return nil, err
	}

	return resp.(*v1.Statistics), nil
}

func newWebServerStatisticsService(factory *factory) WebServerStatisticsService {
	return &webServerStatisticsService{eps: factory.eps.WebServerStatistics()}
}
