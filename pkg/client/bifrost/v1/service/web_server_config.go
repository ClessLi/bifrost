package service

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
)

type WebServerConfigService interface {
	GetServerNames() (*v1.ServerNames, error)
	Get(servername *v1.ServerName) (*v1.WebServerConfig, error)
	Update(config *v1.WebServerConfig) error
}

type webServerConfigService struct {
	eps epv1.WebServerConfigEndpoints
}

func (w *webServerConfigService) GetServerNames() (*v1.ServerNames, error) {
	resp, err := w.eps.EndpointGetServerNames()(GetContext(), nil)
	if err != nil {
		return nil, err
	}

	return resp.(*v1.ServerNames), nil
}

func (w *webServerConfigService) Get(servername *v1.ServerName) (*v1.WebServerConfig, error) {
	resp, err := w.eps.EndpointGet()(GetContext(), servername)
	if err != nil {
		return nil, err
	}

	return resp.(*v1.WebServerConfig), nil
}

func (w *webServerConfigService) Update(config *v1.WebServerConfig) error {
	resp, err := w.eps.EndpointUpdate()(GetContext(), config)
	if err != nil {
		return err
	}
	log.Infof("Update result: %s", resp.(*v1.Response).Message)
	return nil
}

func newWebServerConfigService(factory *factory) WebServerConfigService {
	return &webServerConfigService{eps: factory.eps.WebServerConfig()}
}
