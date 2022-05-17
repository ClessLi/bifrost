package nginx

import (
	"context"

	"github.com/marmotedu/errors"

	v1 "github.com/yongPhone/bifrost/api/bifrost/v1"
	storev1 "github.com/yongPhone/bifrost/internal/bifrost/store/v1"
	"github.com/yongPhone/bifrost/internal/pkg/code"
	"github.com/yongPhone/bifrost/pkg/resolv/V2/nginx/configuration"
)

type webServerConfigStore struct {
	configs map[string]configuration.Configuration
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
			return nil, errors.WithCode(code.ErrInvalidConfig, "nginx server config '%s' is null", servername.Name)
		}

		return &v1.WebServerConfig{
			ServerName: servername,
			JsonData:   jdata,
		}, nil
	}

	return nil, errors.WithCode(code.ErrConfigurationNotFound, "nginx server config '%s' not found", servername.Name)
}

func (w *webServerConfigStore) Update(ctx context.Context, config *v1.WebServerConfig) error {
	if conf, has := w.configs[config.ServerName.Name]; has {
		return conf.UpdateFromJsonBytes(config.JsonData)
	}

	return errors.WithCode(code.ErrConfigurationNotFound, "nginx server config '%s' not found", config.ServerName.Name)
}

var _ storev1.WebServerConfigStore = &webServerConfigStore{}

func newNginxConfigStore(store *webServerStore) storev1.WebServerConfigStore {
	return &webServerConfigStore{configs: store.cms.GetConfigs()}
}
