package service

import (
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
	log "github.com/ClessLi/bifrost/pkg/log/v1"
	"github.com/marmotedu/errors"
)

type WebServerConfigService interface {
	GetServerNames() (servernames []string, err error)
	Get(servername string) ([]byte, error)
	Update(servername string, config []byte) error
}

type webServerConfigService struct {
	eps epv1.WebServerConfigEndpoints
}

func (w *webServerConfigService) GetServerNames() (servernames []string, err error) {
	resp, err := w.eps.EndpointGetServerNames()(GetContext(), nil)
	if err != nil {
		return nil, err
	}

	for _, servername := range *resp.(*v1.ServerNames) {
		servernames = append(servernames, servername.Name)
	}
	return
}

func (w *webServerConfigService) Get(servername string) ([]byte, error) {
	resp, err := w.eps.EndpointGet()(GetContext(), &v1.ServerName{Name: servername})
	if err != nil {
		return nil, err
	}
	response := resp.(*v1.WebServerConfig)
	if response.ServerName.Name != servername {
		return nil, errors.Errorf("get incorrect web server config: get `%s`, want `%s`", response.ServerName.Name, servername)
	}

	return response.JsonData, nil
}

func (w *webServerConfigService) Update(servername string, config []byte) error {
	resp, err := w.eps.EndpointUpdate()(GetContext(), &v1.WebServerConfig{
		ServerName: &v1.ServerName{Name: servername},
		JsonData:   config,
	})
	if err != nil {
		return err
	}
	log.Infof("Update result: %s", resp.(*v1.Response).Message)
	return nil
}

func newWebServerConfigService(factory *factory) WebServerConfigService {
	return &webServerConfigService{eps: factory.eps.WebServerConfig()}
}
