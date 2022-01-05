package nginx

import (
	"context"
	v1 "github.com/ClessLi/bifrost/api/bifrost/v1"
	storev1 "github.com/ClessLi/bifrost/internal/bifrost/store/v1"
	"github.com/ClessLi/bifrost/internal/pkg/code"
	"github.com/ClessLi/bifrost/pkg/resolv/V2/nginx/configuration"
	"github.com/marmotedu/errors"
)

type webServerConfigStore struct {
	configs map[string]configuration.Configuration
}

func (w *webServerConfigStore) Get(ctx context.Context, name *v1.ServerName) (*v1.WebServerConfig, error) {
	if conf, has := w.configs[name.Name]; has {
		jdata := conf.Json()
		if len(jdata) == 0 {
			return nil, errors.WithCode(code.ErrInvalidConfig, "nginx server config '%s' is null", name.Name)
		}
		return &v1.WebServerConfig{
			ServerName: name,
			JsonData:   jdata,
		}, nil
	}
	return nil, errors.WithCode(code.ErrConfigurationNotFound, "nginx server config '%s' not found", name.Name)
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
