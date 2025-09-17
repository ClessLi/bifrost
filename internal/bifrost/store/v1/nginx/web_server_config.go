package nginx

import (
	"context"
	"encoding/json"

	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration"
	utilsV3 "github.com/ClessLi/bifrost/pkg/resolv/V3/nginx/configuration/utils"

	"github.com/marmotedu/errors"
)

type webServerConfigStore struct {
	configs map[string]configuration.NginxConfig
}

func (w *webServerConfigStore) GetServerNames(ctx context.Context) (*v1.ServerNames, error) {
	serverNames := make(v1.ServerNames, 0)
	for name := range w.configs {
		serverNames = append(serverNames, v1.ServerName{Name: name})
	}

	return &serverNames, nil
}

func (w *webServerConfigStore) Get(ctx context.Context, servername *v1.ServerName) (*v1.WebServerConfig, error) {
	if conf, has := w.configs[servername.Name]; has {
		jdata := conf.Json()
		if len(jdata) == 0 {
			return nil, errors.WithCode(code.ErrInvalidConfig, "the web server config '%s' is null", servername.Name)
		}

		fp := utilsV3.NewConfigFingerprinter(conf.Dump()).Fingerprints()
		fpdata, err := json.Marshal(fp)
		if err != nil {
			return nil, errors.WithCode(code.ErrInvalidConfig, "failed to marshal fingerprints of the web server config, cased by: %v", err)
		}

		return &v1.WebServerConfig{
			ServerName:           servername,
			JsonData:             jdata,
			OriginalFingerprints: fpdata,
		}, nil
	}

	return nil, errors.WithCode(code.ErrConfigurationNotFound, "the web server config '%s' not found", servername.Name)
}

func (w *webServerConfigStore) Update(ctx context.Context, config *v1.WebServerConfig) error {
	// check original fingerprints
	var ofp utilsV3.ConfigFingerprints
	err := json.Unmarshal(config.OriginalFingerprints, &ofp)
	if err != nil {
		return errors.WithCode(code.ErrInvalidConfig, "failed to unmarshal fingerprints of the web server config, cased by: %v", err)
	}
	if utilsV3.NewConfigFingerprinter(w.configs[config.ServerName.Name].Dump()).Diff(ofp) {
		return errors.WithCode(code.ErrInvalidConfig, "the original fingerprints(%v) submitted for the config update request does not match the web server config fingerprints", ofp)
	}

	if conf, has := w.configs[config.ServerName.Name]; has {
		return conf.UpdateFromJsonBytes(config.JsonData)
	}

	return errors.WithCode(code.ErrConfigurationNotFound, "the web server config '%s' not found", config.ServerName.Name)
}

var _ storev1.WebServerConfigStore = &webServerConfigStore{}

func newNginxConfigStore(store *webServerStore) storev1.WebServerConfigStore {
	return &webServerConfigStore{configs: store.configsManger.GetConfigs()}
}
