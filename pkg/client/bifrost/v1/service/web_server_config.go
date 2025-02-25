package service

import (
	"encoding/json"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"
	logV1 "github.com/ClessLi/component-base/pkg/log/v1"
	"github.com/marmotedu/errors"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	epv1 "github.com/ClessLi/bifrost/internal/bifrost/endpoint/v1"
)

type WebServerConfigService interface {
	GetServerNames() (servernames []string, err error)
	Get(servername string) (config configuration.NginxConfig, originalFingerprinter utilsV3.ConfigFingerprinter, err error)
	Update(servername string, config configuration.NginxConfig, originalFingerprints utilsV3.ConfigFingerprints) error
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

func (w *webServerConfigService) Get(servername string) (configuration.NginxConfig, utilsV3.ConfigFingerprinter, error) {
	resp, err := w.eps.EndpointGet()(GetContext(), &v1.ServerName{Name: servername})
	if err != nil {
		return nil, nil, err
	}
	response := resp.(*v1.WebServerConfig)
	if response.ServerName.Name != servername {
		return nil, nil, errors.Errorf(
			"get incorrect web server config: get `%s`, want `%s`",
			response.ServerName.Name,
			servername,
		)
	}
	config, err := configuration.NewNginxConfigFromJsonBytes(response.JsonData)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to unmarshal the web server(%s) config", response.ServerName.Name)
	}
	var ofp utilsV3.ConfigFingerprints
	err = json.Unmarshal(response.OriginalFingerprints, &ofp)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to unmarshal fingerprints of the web server(%s) config", response.ServerName.Name)
	}
	return config, utilsV3.SimpleConfigFingerprinter(ofp), nil
}

func (w *webServerConfigService) Update(servername string, config configuration.NginxConfig, ofp utilsV3.ConfigFingerprints) error {
	ofpdata, err := json.Marshal(ofp)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal fingerprints of the web server config")
	}
	resp, err := w.eps.EndpointUpdate()(GetContext(), &v1.WebServerConfig{
		ServerName:           &v1.ServerName{Name: servername},
		JsonData:             config.Json(),
		OriginalFingerprints: ofpdata,
	})
	if err != nil {
		return err
	}
	logV1.Infof("Update result: %s", resp.(*v1.Response).Message)

	return nil
}

func newWebServerConfigService(factory *factory) WebServerConfigService {
	return &webServerConfigService{eps: factory.eps.WebServerConfig()}
}
